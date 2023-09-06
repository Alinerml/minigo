package utils

import (
	"context"
	"encoding/base64"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"minigo/conf"
)

// 封装上传图片到七牛云然后返回状态和图片的url
func UploadToQiNiu(data *multipart.FileHeader, finalname string, fileSize int64) (int, string, string) {
	file, _ := data.Open()
	var AccessKey = conf.AccessKey
	var SerectKey = conf.SerectKey
	var Bucket = conf.Bucket
	var ImgUrl = conf.QiniuServer
	entry := conf.Bucket + ":" + "photo/" + finalname + "jpg"
	encodedEntryURI := base64.StdEncoding.EncodeToString([]byte(entry))
	putPlicy := storage.PutPolicy{
		Scope:         Bucket,
		PersistentOps: "vframe/jpg/offset/1|saveas/" + encodedEntryURI,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadongZheJiang2,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, "video/"+finalname+"mp4", file, fileSize, &putExtra)
	if err != nil {
		return 0, err.Error(), err.Error()
	}

	playurl := "http://" + ImgUrl + "/" + ret.Key
	coverurl := "http://" + ImgUrl + "/" + "photo/" + finalname + "jpg"
	return 200, playurl, coverurl
}
