package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-demo/dao"
	"simple-demo/utils"
	"strconv"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct { //首字母必须大写才能反射，解析json
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1},
		})
		return
	}

	//密码加密存储
	password = utils.Encode(password)
	user := &dao.User{
		Name:       username,
		Password:   password,
		CreateTime: time.Now().Unix(),
	}
	err := dao.SaveUser(user)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	} else {
		token, err := utils.GenerateJWTToken(user.Id)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}

	//token := username + password
	//
	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user_id, err_message, err := dao.UserLogin(username, password)
	//查username

	//判断password是否相等
	if err_message == -1 || err_message == 0 { //用户不存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		token, err := utils.GenerateJWTToken(user_id)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user_id,
				Token:    token,
			})
		}

	}
}

func UserInfo(c *gin.Context) { //查询别人，还没有已关注这个字段
	id := c.Query("user_id")
	user_id, _ := strconv.ParseInt(id, 10, 64)
	user := dao.QueryById(user_id)
	if user.Id == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1,
				StatusMsg: "User doesn't exist",
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: User{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      false,
		},
	})

}
