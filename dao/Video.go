package dao

import (
	"log"
)

type Video struct {
	Id            int64
	PlayUrl       string
	CoverUrl      string
	CreateTime    int64
	AuthId        int64
	FavoriteCount int
	CommentCount  int
	Title         string
}

func (v Video) TableName() string {
	return "video"
}

func SaveVideo(video *Video) error {
	err := DB.Create(video).Error
	if err != nil {
		log.Println("insert user error", err)
		return err
	}
	return err
}

func QueryByTime(time int64) []Video {
	var video_list []Video
	result := DB.Where("create_time < ?", time).
		Order("create_time DESC").
		Limit(30).
		Find(&video_list)
	if result.Error != nil {
		log.Println("Error querying database:", result.Error)
		return video_list
	}
	return video_list
}
