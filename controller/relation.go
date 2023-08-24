package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"simple-demo/conf"
	"simple-demo/dao"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) { //关注操作
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	auth_id := claims["sub"] //得到用户id
	follow_id := int64(auth_id.(float64))

	follower_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type, _ := strconv.Atoi(c.Query("action_type"))
	if action_type == 1 {
		err := dao.Addfollow(follow_id, follower_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
			)
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "关注成功",
			},
			)
		}
	} else {
		err := dao.CancleFollow(follow_id, follower_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
			)
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "取消关注成功",
			},
			)
		}
	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) { //关注列表
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	err, follow_list := dao.QueryFollowByUserId(user_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	follow_res := make([]User, len(follow_list))

	for i, source := range follow_res {

		follow_res[i] = User{
			Id:              source.Id,
			Name:            source.Name,
			FollowCount:     source.FollowCount,
			FollowerCount:   source.FollowerCount,
			IsFollow:        true,
			Avatar:          source.Avatar,
			BackgroundImage: source.BackgroundImage,
			Signature:       source.Signature,
			TotalFavorited:  source.TotalFavorited,
			WorkCount:       source.WorkCount,
			FavoriteCount:   source.FavoriteCount,
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: follow_res,
	})

}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) { //粉丝列表

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	err, follower_list := dao.QueryFollowerByUserId(user_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	follower_res := make([]User, len(follower_list))

	for i, source := range follower_res {
		isfollow := dao.IsFollow(user_id, source.Id)
		follower_res[i] = User{
			Id:              source.Id,
			Name:            source.Name,
			FollowCount:     source.FollowCount,
			FollowerCount:   source.FollowerCount,
			IsFollow:        isfollow,
			Avatar:          source.Avatar,
			BackgroundImage: source.BackgroundImage,
			Signature:       source.Signature,
			TotalFavorited:  source.TotalFavorited,
			WorkCount:       source.WorkCount,
			FavoriteCount:   source.FavoriteCount,
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: follower_res,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
