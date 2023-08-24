package dao

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type Like struct {
	Id         int64
	UserId     int64
	VideoId    int64
	CreateTime int64
}

func (like Like) TableName() string {
	return "like"
}

func Addlike(user_id int64, video_id int64) error {
	var like Like
	if err := DB.Where("user_id=?", user_id).Where("video_id=?", video_id).First(&like).Error; err != gorm.ErrRecordNotFound { //判断是否已经有记录
		return err
	}
	like = Like{
		UserId:     user_id,
		VideoId:    video_id,
		CreateTime: time.Now().Unix(),
	}
	println(123)
	err := DB.Create(&like).Error
	if err != nil {
		log.Println("insert like error", err)
		return err
	}
	//修改video表点赞数
	var video Video
	DB.First(&video, video_id)
	video.FavoriteCount++
	err = DB.Save(&video).Error
	if err != nil {
		log.Println("update video error", err)
		return err
	}
	//修改用户信息的点赞数
	var user User
	DB.First(&user, user_id)
	user.FavoriteCount++
	err = DB.Save(&user).Error //修改喜欢数量
	if err != nil {
		log.Println("update user error", err)
		return err
	}

	var auth User
	DB.First(&auth, video.AuthId)
	auth.TotalFavorited++
	err = DB.Save(&user).Error //修改作者的获赞数量
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	return err
}

func IsLike(user_id int64, video_id int64) bool {
	var like Like
	if err := DB.Where("user_id=?", user_id).Where("video_id=?", video_id).First(&like).Error; err == gorm.ErrRecordNotFound { //判断是否已经有记录
		return false
	}
	return true
}

func Canclelike(user_id int64, video_id int64) error {
	var like Like
	result := DB.Where("user_id = ?", user_id).Where("video_id = ?", video_id).First(&like)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("no record", result.Error)
		return result.Error
	}
	result = DB.Delete(&like)
	if result.Error != nil {
		log.Println("delete like error", result.Error)
		return result.Error
	}

	//修改视频点赞数据
	var video Video
	DB.First(&video, video_id)
	video.FavoriteCount--
	err := DB.Save(&video).Error
	if err != nil {
		log.Println("update video error", err)
		return err
	}
	//修改用户信息的点赞数
	var user User
	DB.First(&user, user_id)
	user.FavoriteCount--
	err = DB.Save(&user).Error
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	var auth User
	DB.First(&auth, video.AuthId)
	auth.TotalFavorited--
	err = DB.Save(&user).Error //修改作者的获赞数量
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	return err
}

func QueryLikesByUserId(user_id int64) (error, []Video) {
	var videolist []Video
	//先查所有喜欢的视频id
	var videos []int64
	result := DB.Model(&Like{}).Where("user_id = ?", user_id).Pluck("video_id", &videos)
	if result.Error != nil {
		return result.Error, videolist
	}
	//再根据id查视频详情
	for _, video := range videos {
		videolist = append(videolist, QueryVideoById(video))
	}
	return result.Error, videolist
}
