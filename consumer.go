package main

import (
    "net/rpc"
    "fmt"
    "log"
    "time"
)

func main() {

    client, err := rpc.Dial("tcp", ADDRESS + ":" + PORT) // connect to te queue (server) on port 
    if err != nil {
        log.Fatal("dialing:", err)
    }


	for{
	   
	    var param Message 
	    //var reply int
		msgReply := new(Message)

	    msgCall := client.Go(QUEUE_NAME + ".PullMessage", &param, &msgReply, nil) // call asynchronously RPC PullMessage
		msgCall = <- msgCall.Done // will be equal to divCall
		if msgCall.Error != nil {
			fmt.Println("Error in message retrieve")
		}
		

		if(msgReply.Text != ""){ // if message received is not empty

			var messageID = msgReply.ID
			fmt.Println("Message id: ", messageID)

			time.Sleep(TIME_BEFORE_ACK_RESPOND * time.Second) // Wait TIME_BEFORE_ACK_RESPOND seconds before send ack

			var response int
			msgCall = client.Go(QUEUE_NAME + ".ReceivedAck", &messageID, &response, nil) // call asynchronously RPC ReceivedAck
			msgCall = <- msgCall.Done // will be equal to divCall
			if msgCall.Error != nil {
				log.Fatal("Error in message retrieve")
			}
						

			if(response != 0){
				fmt.Println("error in ack received for message id: ", messageID)
			}else{
				fmt.Println("Ack received for message id: ", messageID)

				time.Sleep(TIME_FOR_ELABORATION * time.Second) // Wait TIME_FOR_ELABORATION seconds before send elaboration ack


				var response int
				msgCall = client.Go(QUEUE_NAME + ".ElaboratedAck", &messageID, &response, nil) // call asynchronously RPC ElaboratedAck
				msgCall = <- msgCall.Done // will be equal to divCall
				if msgCall.Error != nil {
					log.Fatal("Error in message retrieve")
				}

				fmt.Println("Message: ", msgReply.Text)
				if(response != 0){
					fmt.Println("error in ack elaborated for message id: ", messageID)
				}else{
					fmt.Println("Ack elaborated for message id: ", messageID)
				}
				
				
			}
			


		}else{
			fmt.Println("No messages to retrieve")
		}


	}


   

}