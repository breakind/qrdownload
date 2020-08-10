package main

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const (
	endPoint        = "your oss endpoint"    // oss endpoint
	accessKeyID     = "your oss key"       // key
	accessKeySecret = "your oss secrect" // secrect
	bucketName      = "your oss bucket name"
	appStoreUrl     = "your iOS appstore url"
	apkUrl          = "your android apk url"
	autoUploadApk   = true
	qrHtmlName      = "qrdownload.html"
)

func main() {
	client, err := oss.New("http://"+endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	iOSUrl := appStoreUrl
	androidUrl := apkUrl

	//自动上传当前目录下的第一个apk文件
	if autoUploadApk {
		fileList := []string{}
		err := filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				fileList = append(fileList, path)
				return nil
			})

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}

		fileToUpload := ""
		for _, file := range fileList {
			if path.Ext(file) == ".apk" {
				fileToUpload = file
				break
			}
		}

		if fileToUpload == "" {
			fmt.Println("Error: Can't find apk to upload")
			os.Exit(-1)
		}

		fmt.Println("Start upload apk file...")
		err = bucket.PutObjectFromFile(path.Base(fileToUpload), fileToUpload)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}

		androidUrl = fmt.Sprintf("http://%s.%s/%s", bucketName, endPoint, path.Base(fileToUpload))
	}

	var b bytes.Buffer
	tpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	err = tpl.Execute(&b, map[string]interface{}{
		"android_url": androidUrl,
		"ios_url":     iOSUrl,
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("Start upload html file...")
	err = bucket.PutObject(qrHtmlName, &b)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("Use this address to generate QR image:")
	qrHtml := fmt.Sprintf("http://%s.%s/%s", bucketName, endPoint, qrHtmlName)
	fmt.Println(qrHtml)
}
