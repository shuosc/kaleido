package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

var Bucket *oss.Bucket

func init() {
	client, err := oss.New(os.Getenv("OSS_END_POINT"), os.Getenv("OSS_ACCESS_KEY"), os.Getenv("OSS_ACCESS_SECRET"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	bucketName := "kaleido-message"
	Bucket, err = client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
