package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"minigo/conf"
	"minigo/dao"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) { //点赞操作
	token := c.Query("token")
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	user_id := claims["sub"] //得到用户id
	authid := int64(user_id.(float64))

	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action_id, _ := strconv.Atoi(c.Query("action_type"))
	if action_id == 1 {
		err := dao.Addlike(authid, video_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
			)
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "点赞成功",
			},
			)
		}
	} else {
		err := dao.Canclelike(authid, video_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
			)
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "取消点赞成功",
			},
			)
		}
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) { //喜欢列表
	token := c.Query("token")
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	auth_id := claims["sub"] //得到用户id
	authid := int64(auth_id.(float64))

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	err, video_list := dao.QueryLikesByUserId(user_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
		)
	}
	video_res := make([]Video, len(video_list))

	for i, source := range video_list {
		user := dao.QueryUserById(source.AuthId)
		isfollow := dao.IsFollow(authid, source.AuthId) //判断当前用户与作者是否关注
		video_res[i] = Video{
			Id: source.Id,
			Author: User{
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
			PlayUrl:       source.PlayUrl,
			CoverUrl:      source.CoverUrl,
			FavoriteCount: source.FavoriteCount,
			CommentCount:  source.CommentCount,
			IsFavorite:    true,
			Title:         source.Title,
		}
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: video_res,
	})
}
