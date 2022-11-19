package message

import "time"

type Message struct {
	Message_Id        string    `bson:"_id" json:"id"`
	Sender_UserName   string    `json:"sender_username"`
	Receiver_UserName string    `json:"receiver_username"`
	Message_Content   string    `json:"message_content"`
	Message_Time      time.Time `json:"message_time"`
}

type GroupMessage struct {
	Group_Id        string    `bson:"_id" json:"id"`
	Sender_Id       string    `json:"sender_id"`
	Message_Content string    `json:"message_content"`
	Message_Time    time.Time `json:"message_time"`
}
