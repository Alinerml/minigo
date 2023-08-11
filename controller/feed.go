package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latest_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if latest_time == 0 {
		latest_time = time.Now().Unix()
	}
	video_list := dao.QueryByTime(latest_time)
	video_res := make([]Video, len(video_list))

	for i, source := range video_list {
		user := dao.QueryById(source.AuthId)
		video_res[i] = Video{
			Id: source.Id,
			Author: User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false,
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
