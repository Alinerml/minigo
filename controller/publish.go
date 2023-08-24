package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"simple-demo/conf"
	"simple-demo/dao"
	"simple-demo/utils"
	"strconv"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) { //发布视频
	token := c.PostForm("token")
	data, err := c.FormFile("data") //获取文件
	title := c.PostForm("title")
	file, _ := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//提取封面

	//vc, err := gocv.VideoCaptureFile(filePath)
	//if err != nil {
	//	return
	//}
	//
	//frames := vc.Get(gocv.VideoCaptureFrameCount)
	//fps := vc.Get(gocv.VideoCaptureFPS)
	//duration := frames / fps
	//
	//frames = (time / duration) * frames
	//
	//// Set Video frames
	//vc.Set(gocv.VideoCapturePosFrames, frames)
	//
	//img := gocv.NewMat()
	//
	//vc.Read(&img)
	//
	//imageObject, err := img.ToImage()
	//if err != nil {
	//	return i, err
	//}

	//
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	user_id := claims["sub"]                                                             //得到用户id
	authid := int64(user_id.(float64))                                                   //todo:蜜汁操作
	finalName := fmt.Sprintf("video/%d_%s_%d", authid, data.Filename, time.Now().Unix()) //保存格式为id_name
	//saveFile := filepath.Join("./public/", finalName) //保存文件到服务器
	code, res := utils.UploadToQiNiu(file, finalName, data.Size)
	if code == 0 { //上传失败
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  res,
		})
		return
	}

	//上传成功，向数据库插入数据
	video := &dao.Video{
		PlayUrl:    res,
		AuthId:     authid,
		CreateTime: time.Now().Unix(),
		Title:      title,
	}
	err = dao.SaveVideo(video)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  res,
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) { //发布列表
	token := c.Query("token")
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	auth_id := claims["sub"] //得到用户id
	authid := int64(auth_id.(float64))
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	err, video_list := dao.QueryVideosByUserId(user_id)

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
		)
	}
	video_res := make([]Video, len(video_list))

	for i, source := range video_list {
		user := dao.QueryUserById(user_id)
		//判断是否点赞
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: video_res,
	})
}
