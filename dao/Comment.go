package dao

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Comment struct {
	Id          int64
	VideoId     int64
	UserId      int64
	ToUserId    int64
	CommentText string
	CommentTime string
}

func (comment Comment) TableName() string {
	return "comment"
}

func AddComment(video_id int64, user_id int64, text string) (error, Comment) { //添加评论 需要返回comment
	var comment Comment
	comment = Comment{
		VideoId:     video_id,
		UserId:      user_id,
		CommentTime: fmt.Sprintf("%d-%d", time.Now().Month(), time.Now().Day()),
		CommentText: text,
	}
	result := DB.Create(&comment)

	if result.Error != nil {
		log.Println("insert comment error", result.Error)

		return result.Error, comment
	}
	//更改视频的评论数
	var video Video
	DB.First(&video, comment.VideoId)
	video.CommentCount++
	err := DB.Save(&video).Error
	if err != nil {
		log.Println("update video error", err)
		return err, comment
	}
	return result.Error, comment
}

func DeleteComment(comment_id int64) error {
	var comment Comment
	result := DB.Where("id = ?", comment_id).First(&comment)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("no record", result.Error)
		return result.Error
	}
	result = DB.Delete(&comment)
	if result.Error != nil {
		log.Println("delete comment error", result.Error)
		return result.Error
	}
	//更改视频的评论数
	var video Video
	DB.First(&video, comment.VideoId)
	video.CommentCount++
	err := DB.Save(&video).Error
	if err != nil {
		log.Println("update video error", err)
		return err
	}
	return result.Error
}

func QueryCommentByVideoId(video_id int64) (error, []Comment) { //评论列表
	var comment_list []Comment

	result := DB.Where("video_id = ?", video_id).Find(&comment_list)
	if result.Error != nil {
		log.Println("Error querying database:", result.Error)
		return result.Error, comment_list
	}

	return result.Error, comment_list
}
