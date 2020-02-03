package main

const QUEUE_NAME = "Queue"

const ADDRESS = "127.0.0.1"
const PORT = "34219"

const TIMEOUT_RETRASMIT = 10
const TIMEOUT_VISIBILITY = 15

const TIME_FOR_SEND = 2
const TIME_BEFORE_ACK_RESPOND = 2
const TIME_FOR_ELABORATION = 3

const LOG_ENABLED = false


// DON'T CHANGE !!!
const INQUEUE = 1000
const SENT = 1001
const RECEIVED = 1002
const ELABORATED = 1003

type Message struct {
	ID     int
	Text   string
	Status int
}