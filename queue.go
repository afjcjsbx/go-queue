package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)


var l = log.New(os.Stdout, "", 0)
var s ItemQueue
type Queue int

// ItemQueue the queue of messages
type ItemQueue struct {
	counter int
	items   []Message
	lock    sync.Mutex
}

// New creates a new ItemQueue
func (s *ItemQueue) New() *ItemQueue {
	s.items = []Message{}
	s.counter = 1

	return s
}

// IsEmpty returns true if the queue is empty
func (s *ItemQueue) IsEmpty() bool {
	s.lock.Lock()
	var size = 0
	for i := 0; i < len(s.items); i++ {
		if s.items[i].Status == INQUEUE {
			size++
		}
	}
	s.lock.Unlock()
	return size == 0
}

// Size returns the number of messages in the queue
func (s *ItemQueue) Size() int {
	s.lock.Lock()
	var size = len(s.items)
	s.lock.Unlock()
	return size
}

func initQueue() *ItemQueue {
	if s.items == nil {
		s = ItemQueue{}
		s.New()
	}
	return &s
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Enqueue(t Message) {
	s.lock.Lock()

	t.ID = s.counter   // Automatically set message id
	t.Status = INQUEUE // Set message status INQUEUE

	s.items = append(s.items, t)
	s.counter = s.counter + 1 // Increment counter ids

	s.lock.Unlock()
}

// Dequeue set a message Status of first message with Status INQUEUE as SENT
func (s *ItemQueue) Dequeue() *Message {
	s.lock.Lock()

	var item Message
	for i := 0; i < len(s.items); i++ {
		if s.items[i].Status == INQUEUE {
			s.items[i].Status = SENT
			item = s.items[i]
			s.lock.Unlock()
			return &item
		}
	}

	s.lock.Unlock()
	return &item
}

// Return an array that contain all messages in
func (s *ItemQueue) GetMessages() []int {
	var messages []int
	s.lock.Lock()

	for i := 0; i < len(s.items); i++ {
			messages[i] = s.items[i].ID
	}

	s.lock.Unlock()
	return messages
}


//Handler for timeout retrasmit
func MessageThreadTimeoutRetrasmit(messageId int) {
	Log("Starting TimeoutRetrasmit thread of message id: " + strconv.Itoa(messageId))

	start := time.Now()
	for {
		for i := 0; i < s.Size(); i++ {

			Log("TimeoutRetrasmit thread of message id: " + strconv.Itoa(messageId))

			if time.Now().After(start.Add(time.Second * TIMEOUT_RETRASMIT)) {
				start = time.Now()

				if s.items[i].ID == messageId && s.items[i].Status == SENT {
					s.items[i].Status = INQUEUE
					Log("Retrasmit of message id: " + strconv.Itoa(messageId))

					return
				}

			}

			if s.items[i].ID == messageId && (s.items[i].Status == RECEIVED || s.items[i].Status == ELABORATED) {
				return
			}

		}
	}
}

//Handler for timeout based
func MessageThreadTimeoutBased(messageId int) {

	MessageThreadTimeoutRetrasmit(messageId)
	start := time.Now()

	for {
		for i := 0; i < s.Size(); i++ {

			if time.Now().After(start.Add(time.Second * TIMEOUT_VISIBILITY)) {
				start = time.Now()

				if s.items[i].ID == messageId && s.items[i].Status == RECEIVED {

					s.items[i].Status = INQUEUE
					Log("Elapsed Timeout based, Retrasmit of message id: " + strconv.Itoa(messageId))
					return
				}


			}

			if s.items[i].ID == messageId && s.items[i].Status == ELABORATED {
				s.RemoveMessageFromQueue(messageId)
				return
			}

		}
	}
}

func (s *ItemQueue) SetMessageReceived(messageId int) bool {
	s.lock.Lock()

	for i := 0; i < len(s.items); i++ {
		if s.items[i].ID == messageId && s.items[i].Status == SENT {
			s.items[i].Status = RECEIVED
			s.lock.Unlock()

			return true
		}
	}

	s.lock.Unlock()
	return false
}

func (s *ItemQueue) SetMessageElaborated(messageId int) bool {
	s.lock.Lock()

	for i := 0; i < len(s.items); i++ {
		if s.items[i].ID == messageId && s.items[i].Status == RECEIVED {
			s.items[i].Status = ELABORATED
			s.lock.Unlock()
			return true
		}
	}

	s.lock.Unlock()
	return false
}

func (s *ItemQueue) RemoveMessageFromQueue(messageId int) bool {
	s.lock.Lock()

	for i := 0; i < len(s.items); i++ {
		if s.items[i].Status == ELABORATED && s.items[i].ID == messageId {
			s.items = append(s.items[:i], s.items[i+1:]...)
			s.lock.Unlock()
			return true
		}
	}

	s.lock.Unlock()
	return false
}

func (s *ItemQueue) PrintMessages() string {

	var str = "\n"
	for i := 0; i < len(s.items); i++ {
		str = str + "Message id: " + strconv.Itoa(s.items[i].ID) + " text: " + s.items[i].Text +
			" Status: " + strconv.Itoa(s.items[i].Status) + "\n"
	}
	return str
}


func (t *Queue) GetMessageList(msg *Message, array *[]int) error {

	if !s.IsEmpty() {
		*array = s.GetMessages()
	} else {
		*array = append(*array, -1)
		Log("No messages in queue")
	}
	
	return nil
}

func (t *Queue) PullMessage(msg *Message, msg1 *Message) error {

	if !s.IsEmpty() {
		var message = s.Dequeue()
		msg1.Text = message.Text
		msg1.ID = message.ID
		go MessageThreadTimeoutBased(message.ID) // Start thread to check timeouts for message
	} else {
		Log("No messages in queue")
	}
	return nil
}

func (t *Queue) PushMessage(message string, reply *int) error {
	// output message received

	var msg Message
	msg.Text = message
	s.Enqueue(msg) // Push message in queue

	s.PrintMessages()

	*reply = 1
	return nil
}

func (t *Queue) ReceivedAck(messageId int, reply *int) error {
	// output message received

	if s.SetMessageReceived(messageId) {
		Log("ACK received for message id: " + strconv.Itoa(messageId))
		*reply = 0
	} else {
		Log("Error in receive ACK for message id: " + strconv.Itoa(messageId))
		*reply = -1
	}
	return nil
}

func (t *Queue) ElaboratedAck(messageId int, reply *int) error {
	// output message received
	Log("ACK elaborated for message id: " + strconv.Itoa(messageId))

	if s.SetMessageElaborated(messageId) {
		*reply = 0
	} else {
		Log("Error in ACK elaborated for message id: " + strconv.Itoa(messageId))
		*reply = -1
	}
	return nil
}

func main() {

	Log("Launching queue")

	queue := new(Queue)
	queue_server := rpc.NewServer()

	err := queue_server.RegisterName(QUEUE_NAME, queue)
	checkError(err)

	listener, err := net.Listen("tcp", ":" + PORT)
	checkError(err)

	Log("Queue started")
	go LogQueue()

	/* This works:
	   rpc.Accept(listener)
	*/
	/* and so does this:
	 */
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go queue_server.ServeConn(conn)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func LogQueue() {
	if LOG_ENABLED {
		for {
			time.Sleep(time.Second)
			l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [Queue] Messages in queue: " + strconv.Itoa(s.Size()))
			l.Print(s.PrintMessages())
		}
	}
}

func Log(s string) {
	if LOG_ENABLED {
		time.Sleep(time.Second)
		l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [Queue] ")
		l.Print(s + "\n")
	}
}
