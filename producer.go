package main

import (
    "net/rpc"
    "fmt"
    "log"
    //"bufio"
    //"os"
    //"strings"
    "math/rand"
    "time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))


func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func String(length int) string {
  return StringWithCharset(length, charset)
}



func main() {

    client, err := rpc.Dial("tcp", ADDRESS + ":" + PORT)
    if err != nil {
        log.Fatal("dialing:", err)
    }


    for{
      var reply int
      text := String(15)
      msgCall := client.Go(QUEUE_NAME + ".PushMessage", text, &reply, nil)
    replyCall := <-msgCall.Done // will be equal to divCall
    if(replyCall.Error != nil){
      fmt.Printf("Error: %d\n", replyCall.Error.Error())
    }else{
      // no error
      fmt.Printf("Result: %d\n", reply)
    }

    time.Sleep(TIME_FOR_SEND * time.Second)
  }
  


}



