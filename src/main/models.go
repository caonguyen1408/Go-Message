package main

type Message struct {
	Message string `json:"message"`
}
type QueueMessage struct {
	Message string `json:"message"`
	TaskId  int64  `json:"task_id"`
}
