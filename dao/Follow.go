package dao

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type Follow struct {
	Id         int64
	FollowerId int64 `gorm:"column:follower_id"`
	FollowId   int64 `gorm:"column:follow_id"`
	CreateTime int64 `gorm:"column:create_time"`
}

func (follow Follow) TableName() string {
	return "follow"
}

func Addfollow(follower_id int64, follow_id int64) error {
	var follow Follow
	if err := DB.Where("follower_id=?", follower_id).Where("follow_id=?", follow_id).First(&follow).Error; err != gorm.ErrRecordNotFound { //判断是否已经有记录
		return err
	}
	follow = Follow{
		FollowerId: follower_id,
		FollowId:   follow_id,
		CreateTime: time.Now().Unix(),
	}
	err := DB.Create(&follow).Error
	if err != nil {
		log.Println("insert follow error", err)
		return err
	}

	//修改用户信息的关注数和粉丝数
	var follower User
	DB.First(&follower, follower_id)
	follower.FollowCount++
	err = DB.Save(&follower).Error
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	var followw User //被关注的
	DB.First(&followw, follow_id)
	followw.FollowerCount++
	err = DB.Save(&followw).Error
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	return err
}

func IsFollow(follower_id int64, follow_id int64) bool {
	var follow Follow
	if err := DB.Where("follower_id=?", follower_id).Where("follow_id=?", follow_id).First(&follow).Error; err == gorm.ErrRecordNotFound { //判断是否已经有记录
		return false
	}
	return true
}

func CancleFollow(follower_id int64, follow_id int64) error {
	var follow Follow
	result := DB.Where("follower_id = ?", follow_id).Where("follow_id = ?", follow_id).First(&follow)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("no record", result.Error)
		return result.Error
	}
	result = DB.Delete(&follow)
	if result.Error != nil {
		log.Println("delete follow error", result.Error)
		return result.Error
	}

	//修改用户信息
	var follower User
	DB.First(&follower, follower_id)
	follower.FollowCount++
	err := DB.Save(&follower).Error
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	var followw User //被关注的
	DB.First(&followw, follow_id)
	followw.FollowerCount++
	err = DB.Save(&followw).Error
	if err != nil {
		log.Println("update user error", err)
		return err
	}
	return err
}

func QueryFollowByUserId(user_id int64) (error, []User) { //关注列表
	var follow_list []User
	//先查所有关注的id
	var follows []int64
	result := DB.Model(&Follow{}).Where("follower_id = ?", user_id).Pluck("follow_id", &follows)
	if result.Error != nil {
		return result.Error, follow_list
	}
	//再根据id查用户
	for _, user_id := range follows {
		follow_list = append(follow_list, QueryUserById(user_id))
	}
	return result.Error, follow_list
}

func QueryFollowerByUserId(user_id int64) (error, []User) { //粉丝列表
	var follower_list []User
	//先查所有关注的id
	var followers []int64
	result := DB.Model(&Follow{}).Where("follow_id = ?", user_id).Pluck("follower_id", &followers)
	if result.Error != nil {
		return result.Error, follower_list
	}
	//再根据id查用户
	for _, user_id := range followers {
		follower_list = append(follower_list, QueryUserById(user_id))
	}
	return result.Error, follower_list
}
