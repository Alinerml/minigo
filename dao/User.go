package dao

import (
	"github.com/RaymondCode/simple-demo/utils"
	"log"
	"time"
)

type User struct {
	ID              int64
	UserName        string    `gorm:"column:user_name"`
	Password        string    `gorm:"column:password"`
	CreateTime      time.Time `gorm:"column:create_time"`
	FollowCount     int       `gorm:"column:follow_count"`
	FollowerCount   int       `gorm:"column:follower_count"`
	Avatar          string    `gorm:"column:avatar"`
	BackgroundImage string    `gorm:"column:background_image"`
	Signature       string    `gorm:"column:signature"`
	TotalFavorited  int       `gorm:"column:total_favorited"`
	WorkCount       int       `gorm:"column:work_count"`
	FavoriteCount   int       `gorm:"column:favorite_count"`
}

func (u User) TableName() string {
	return "user"
}

func SaveUser(user *User) error {
	err := DB.Create(user).Error
	if err != nil {
		log.Println("insert user error", err)
		return err
	}
	return err
}

func UserLogin(username string, password string) (int64, int) {
	password = utils.Encode(password)
	var u User
	err := DB.Where("user_name=?", username).First(&u).Error
	if err != nil {
		log.Println("login error", err)
		return 0, -1
	}
	if password != u.Password {
		return 0, 0
	}
	return u.ID, 1
}
