package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"minigo/conf"
	"minigo/dao"
	"minigo/utils"
	"net/http"
	"strconv"
	"time"
)

type UserLoginResponse struct { //首字母必须大写才能反射，解析json
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) { //注册
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
		Name:            username,
		Password:        password,
		Avatar:          "http://s0ce3vuua.bkt.clouddn.com/photo/20230905140635.jpg",
		BackgroundImage: "http://s0ce3vuua.bkt.clouddn.com/photo/20230905140059.jpg",
		Signature:       "这是一个固定的个性签名",
		CreateTime:      time.Now().Unix(),
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
}

func Login(c *gin.Context) { //登录
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
				Response: Response{StatusCode: 0, StatusMsg: "登录成功"},
				UserId:   user_id,
				Token:    token,
			})
		}

	}
}

func UserInfo(c *gin.Context) { //查询别人
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusOK, Response{
			-1,
			"未登录",
		})
		return
	}
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	auth_id := claims["sub"] //得到用户id
	authid := int64(auth_id.(float64))
	id := c.Query("user_id")
	user_id, _ := strconv.ParseInt(id, 10, 64)
	user := dao.QueryUserById(user_id)
	isfollow := dao.IsFollow(authid, user_id)
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
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			IsFollow:        isfollow,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		},
	})

}
