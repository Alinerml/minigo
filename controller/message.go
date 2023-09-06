package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"minigo/conf"
	"minigo/dao"
	"net/http"
	"strconv"
)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) { //发送消息
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	user_id := claims["sub"] //得到用户id
	authid := int64(user_id.(float64))
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	println(userIdB, authid)
	err := dao.AddMessage(userIdB, authid, content)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "发送消息失败"})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0,
		StatusMsg: "消息发送成功"})

}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) { //查询消息
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	pre_msg_time := c.Query("pre_msg_time")
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	user_id := claims["sub"] //得到用户id
	authid := int64(user_id.(float64))
	pre_time, _ := strconv.ParseInt(pre_msg_time, 10, 64)
	userIdB, _ := strconv.ParseInt(toUserId, 10, 64)
	err, message_list := dao.QueryMessage(authid, userIdB, pre_time)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "查询消息记录失败",
		})
		return
	} else {
		message_res := make([]Message, len(message_list))

		for i, source := range message_list {
			message_res[i] = Message{
				Id:         source.Id,
				Content:    source.Content,
				CreateTime: source.CreateTime,
				ToUserId:   source.ToUserId,
				FromUserId: source.FromUserId,
			}
		}
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "查询消息记录成功"}, MessageList: message_res})
	}

}
