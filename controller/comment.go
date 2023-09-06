package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"minigo/conf"
	"minigo/dao"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	user_id := claims["sub"] //得到用户id
	authid := int64(user_id.(float64))
	user := dao.QueryUserById(authid)
	if actionType == "1" { //添加评论
		text := c.Query("comment_text")
		err, comment := dao.AddComment(videoId, authid, text)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error()})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0,
			StatusMsg: "添加评论成功"},
			Comment: Comment{
				Id: comment.Id,
				User: User{
					Id:              user.Id,
					Name:            user.Name,
					FollowCount:     user.FollowerCount,
					FollowerCount:   user.FollowCount,
					IsFollow:        true,
					Avatar:          user.Avatar,
					BackgroundImage: user.BackgroundImage,
					Signature:       user.Signature,
					TotalFavorited:  user.TotalFavorited,
					WorkCount:       user.WorkCount,
					FavoriteCount:   user.FavoriteCount,
				},
				Content:    comment.CommentText,
				CreateDate: comment.CommentTime,
			}})
		return
	} else { //删除评论
		comment_id, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		err := dao.DeleteComment(comment_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "删除评论成功"})
	}

}

//func CommentAction(c *gin.Context) {
//	token := c.Query("token")
//	actionType := c.Query("action_type")
//
//	if user, exist := usersLoginInfo[token]; exist {
//		if actionType == "1" {
//			text := c.Query("comment_text")
//			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
//				Comment: Comment{
//					Id:         1,
//					User:       user,
//					Content:    text,
//					CreateDate: "05-01",
//				}})
//			return
//		}
//		c.JSON(http.StatusOK, Response{StatusCode: 0})
//	} else {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
//	}
//}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) { //评论列表
	token := c.Query("token")
	token_p, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	claims, _ := token_p.Claims.(jwt.MapClaims)
	user_id := claims["sub"] //得到用户id
	authid := int64(user_id.(float64))

	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	err, comment_list := dao.QueryCommentByVideoId(video_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	comment_res := make([]Comment, len(comment_list))

	for i, source := range comment_list {
		user := dao.QueryUserById(source.UserId)
		isFollow := dao.IsFollow(authid, source.Id)
		comment_res[i] = Comment{
			Id: source.Id,
			User: User{
				Id:              user.Id,
				Name:            user.Name,
				FollowCount:     user.FollowCount,
				FollowerCount:   user.FollowerCount,
				IsFollow:        isFollow,
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				Signature:       user.Signature,
				TotalFavorited:  user.TotalFavorited,
				WorkCount:       user.WorkCount,
				FavoriteCount:   user.FavoriteCount,
			},
			Content:    source.CommentText,
			CreateDate: source.CommentTime,
		}
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询评论列表成功",
		},
		CommentList: comment_res,
	})

}
