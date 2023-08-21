package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"simple-demo/conf"
	"simple-demo/dao"
	"simple-demo/utils"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
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
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
