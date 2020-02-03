go-queue
============

Simple message queue with producer and consumer

## Limits ##

- Defined message structure to send message
- It does not recover messages if the queue crashes

## Installation ##

The recommended way to install *go-queue* is Makefile:

- Do to the project folder
- `make` in command line


	
## Examples ##

- Run a queue
	`go run queue.go`
- Run a producer (or more)
	`go run producer.go`
- Run a consumer (or more)
	`go run consumer.go`
	
## Parameters (file const.go) ##

| Parameter                 | Default       | Description   |	
| :------------------------ |:-------------:| :-------------|
| QUEUE_NAME 	       |	"Queue"          |queue name 
| ADDRESS 	       |	"127.0.0.1"          |queue addess 
| PORT          | 12345           |port number 
| TIMEOUT_RETRASMIT 	       |	10	            |time to receive the ACK before putting the message back in the queue (in seconds)
| TIMEOUT_VISIBILITY		       | 15	           | time to elaborate the message before putting the message back in the queue (in seconds)
| TIME_FOR_SEND 	        | 1         | time that producer waits before sending a new message (for testing) (in seconds)
| TIME_BEFORE_ACK_RESPOND         | 1             | time that consumer waits before sending ACK for received message (for testing) (in seconds)
| TIME_FOR_ELABORATION          | 1           | time that consumer takes to elaborate the message (for testing) (in seconds)



## TODOs ##

- Save unsent messages in a log to recover them in case of crash
- replicate the queue for greater tolerance

