package dao

import (
	"log"
	"time"
)

type Message struct {
	Id         int64
	ToUserId   int64
	FromUserId int64
	Content    string
	CreateTime int64
}

func (message Message) TableName() string {
	return "message"
}

func AddMessage(ToUserId int64, FromUserId int64, Content string) error { //添加评论 需要返回comment
	var message Message
	message = Message{
		ToUserId:   ToUserId,
		FromUserId: FromUserId,
		Content:    Content,
		CreateTime: time.Now().Unix(),
	}
	result := DB.Create(&message)

	if result.Error != nil {
		log.Println("insert message error", result.Error)
		return result.Error
	}
	return result.Error
}

func QueryMessage(FromUserId int64, ToUserId int64, pre_msg_time int64) (error, []Message) { //评论列表
	var message_list []Message

	result := DB.Where("(from_user_id = ? and to_user_id = ? and create_time > ?) or (from_user_id = ? and to_user_id = ? and create_time > ?)", FromUserId, ToUserId, pre_msg_time, ToUserId, FromUserId, pre_msg_time).Find(&message_list)
	if result.Error != nil {
		log.Println("Error querying database:", result.Error)
		return result.Error, message_list
	}

	return result.Error, message_list
}
