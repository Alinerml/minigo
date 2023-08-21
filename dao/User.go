package dao

import (
	"log"
	"simple-demo/utils"
)

type User struct {
	Id              int64
	Name            string `gorm:"column:name"`
	Password        string `gorm:"column:password"`
	CreateTime      int64  `gorm:"column:create_time"`
	FollowCount     int    `gorm:"column:follow_count"`
	FollowerCount   int    `gorm:"column:follower_count"`
	Avatar          string `gorm:"column:avatar"`
	BackgroundImage string `gorm:"column:background_image"`
	Signature       string `gorm:"column:signature"`
	TotalFavorited  int    `gorm:"column:total_favorited"`
	WorkCount       int    `gorm:"column:work_count"`
	FavoriteCount   int    `gorm:"column:favorite_count"`
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

func UserLogin(username string, password string) (int64, int, error) {
	password = utils.Encode(password)
	var u User
	err := DB.Where("name=?", username).First(&u).Error
	if err != nil {
		log.Println("login error", err)
		return 0, -1, err //查询失败
	}
	if password != u.Password {
		return 0, 0, err //密码错误
	}
	return u.Id, 1, err //成功
}

func QueryById(user_id int64) User {
	var u User
	err := DB.First(&u, user_id).Error
	if err != nil {
		log.Println("query error", err)
		return u
	}
	return u
}
