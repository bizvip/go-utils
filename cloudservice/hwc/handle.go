/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package hwc

import (
	"errors"
	"fmt"
	"os"

	"github.com/bizvip/go-utils/cloudservice/hwc/obs"
)

// HWC_USR=hid_sn0f9rjvmjulkp2
// HWC_ACCESS_KEY=
// HWC_SECRET_KEY=
// HWC_OBS_BUCKET=
// HWC_OBS_ID=062bff97380e4690a7bd99984e030267
// HWC_OBS_ENDPOINT=
// HWC_OBS_DOMAIN=rc001.obs.ap-southeast-1.myhuaweicloud.com
// HWC_OBS_REFERER=https://www.icbc.com.cn
// HWC_OBS_SCHEMA=https://

var (
	obsClient  *obs.ObsClient
	endpoint   string
	ak         string
	sk         string
	bucketName string
	metadata   map[string]string
)

func getObsClient() *obs.ObsClient {
	// endpoint = config.All.Hwc.Endpoint
	// ak = config.All.Hwc.Ak
	// sk = config.All.Hwc.Sk

	var err error
	if obsClient == nil {
		obsClient, err = obs.New(ak, sk, endpoint)
		if err != nil {
			panic(err)
		}
	}

	return obsClient
}

func init() {
	defer obs.CloseLog()
	_ = obs.InitLog(
		"runtime/quicklog/obs-sdk.log",
		1024*1024*100,
		5,
		obs.LEVEL_WARN,
		false,
	)
	metadata = make(map[string]string)
}

func getBucketName() string {
	if bucketName == "" {
		bucketName = ""
	}
	return bucketName
}

func getMetaData() map[string]string {
	return metadata
}

// func putObjectWithCallback(key string) {
// 	input := &obs.PutObjectInput{}
// 	input.Bucket = bucketName
// 	input.Key = key
//
// 	callbackMap := map[string]string{}
// 	callbackMap["callbackUrl"] = "http://example.com:80"
// 	// callbackMap["callbackHost"] = "example.com"
// 	callbackMap["callbackBody"] = "key=$(key)&size=$(size)&bucket=$(bucket)&etag=$(etag)"
// 	// callbackMap["callbackBodyType"] = "application/x-www-form-urlencoded"
// 	callbackBuffer := bytes.NewBuffer([]byte{})
// 	callbackEncoder := json.NewEncoder(callbackBuffer)
// 	// do not encode '&' to "\u0026"
// 	callbackEncoder.SetEscapeHTML(false)
// 	err := callbackEncoder.Encode(callbackMap)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	callbackVal := base64.StdEncoding.EncodeToString(callbackBuffer.Bytes())
//
// 	input.Body = strings.NewReader("Hello OBS")
//
// 	output, err := getObsClient().PutObject(input, obs.WithCallbackHeader(callbackVal))
// 	if err == nil {
// 		defer output.CloseCallbackBody()
//
// 		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
// 		fmt.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
// 		p := make([]byte, 1024)
// 		var readErr error
// 		var readCount int
// 		for {
// 			readCount, readErr = output.ReadCallbackBody(p)
// 			if readCount > 0 {
// 				fmt.Printf("%s", p[:readCount])
// 			}
// 			if readErr != nil {
// 				break
// 			}
// 		}
// 	} else {
// 		if obsError, ok := err.(obs.ObsError); ok {
// 			fmt.Println(obsError.StatusCode)
// 			fmt.Println(obsError.Code)
// 			fmt.Println(obsError.Message)
// 		} else {
// 			fmt.Println(err)
// 		}
// 	}
// }

func PutObject(filePath string, key string) bool {
	input := &obs.PutObjectInput{}
	input.Bucket = getBucketName()
	input.Key = key
	input.Metadata = getMetaData()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	input.Body = file

	_, err = getObsClient().PutObject(input)

	if err == nil {
		// fmt.Printf("%+v", resp)
		// fmt.Printf("StatusCode:%d, RequestId:%s\n", resp.StatusCode, resp.RequestId)
		// fmt.Printf("ETag:%s, StorageClass:%s\n", resp.ETag, resp.StorageClass)
		return true
	}

	var obsError obs.ObsError
	if errors.As(err, &obsError) {
		fmt.Println(obsError.StatusCode)
		fmt.Println(obsError.Code)
		fmt.Println(obsError.Message)
	}

	return false
}
