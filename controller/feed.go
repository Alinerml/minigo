package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"simple-demo/conf"
	"simple-demo/dao"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) { //首页查询视频
	//判断是否登录
	token := c.Query("token")
	if token != "" { //已登录
		token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return conf.SecretKey, nil
		})
		claims, _ := token_p.Claims.(jwt.MapClaims)
		user_id := claims["sub"] //得到用户id
		authid := int64(user_id.(float64))
		latest_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
		if latest_time == 0 {
			latest_time = time.Now().Unix()
		}
		video_list := dao.QueryVideosByTime(latest_time)
		video_res := make([]Video, len(video_list))

		for i, source := range video_list {
			user := dao.QueryUserById(source.AuthId)
			islike := dao.IsLike(authid, source.Id)
			isfollow := dao.IsFollow(authid, source.AuthId)
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
				CoverUrl:      "http://rz2cue1gw.bkt.clouddn.com/photo/123%202023-08-11%20184530.png",
				FavoriteCount: source.FavoriteCount,
				CommentCount:  source.CommentCount,
				IsFavorite:    islike,
			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: video_res,
		})
	} else {

		latest_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
		if latest_time == 0 {
			latest_time = time.Now().Unix()
		}
		video_list := dao.QueryVideosByTime(latest_time)
		video_res := make([]Video, len(video_list))

		for i, source := range video_list {
			user := dao.QueryUserById(source.AuthId)
			video_res[i] = Video{
				Id: source.Id,
				Author: User{
					Id:              user.Id,
					Name:            user.Name,
					FollowCount:     user.FollowCount,
					FollowerCount:   user.FollowerCount,
					IsFollow:        false,
					Avatar:          user.Avatar,
					BackgroundImage: user.BackgroundImage,
					Signature:       user.Signature,
					TotalFavorited:  user.TotalFavorited,
					WorkCount:       user.WorkCount,
					FavoriteCount:   user.FavoriteCount,
				},
				PlayUrl:       source.PlayUrl,
				CoverUrl:      "http://rz2cue1gw.bkt.clouddn.com/photo/123%202023-08-11%20184530.png",
				FavoriteCount: source.FavoriteCount,
				CommentCount:  source.CommentCount,
				IsFavorite:    false,
			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: video_res,
		})
	}

}
