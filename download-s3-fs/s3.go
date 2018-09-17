package main

import (
	"strings"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"context"
	"time"
)

const defaultTimeOut = 20 * time.Second

type ScS3 struct {
	ScS3Conn *s3.S3
}

func S3Conn(blc string) (ScS3) {
	/*
	创建s3连接
	*/
	var scs3 ScS3
	sess, _ := session.NewSession()
	svc := s3.New(sess, &aws.Config{Region:aws.String(blc)})
	scs3.ScS3Conn = svc

	return scs3
}

func (r *ScS3) SCCreateBucket(bucketName string) error {
	/*
	创建桶
	 */
	params := &s3.CreateBucketInput{
		Bucket:aws.String(bucketName),
		//CreateBucketConfiguration:&s3.CreateBucketConfiguration{
		//	LocationConstraint:aws.String(r.BLC),
		//},
	}
	_, err := r.ScS3Conn.CreateBucket(params)
	return err
}

func (r *ScS3) SCDeleteBucket(bucketName string) error {
	/*
	删除桶
	 */
	params := &s3.DeleteBucketInput{
		Bucket:aws.String(bucketName),
	}
	_, err := r.ScS3Conn.DeleteBucket(params)
	return err
}

func (r *ScS3) SCListObjects(bucketName string, keyRoute string) (objLst *s3.ListObjectsOutput, err error) {
	/*
	列出桶路径下的所有object
	 */
	params := &s3.ListObjectsInput{
		Bucket:aws.String(bucketName),
		Prefix:aws.String(keyRoute),
	}
	objLst, err = r.ScS3Conn.ListObjects(params)
	return objLst, err
}

func (r *ScS3) SCPutObject(bucketName string, keyRoute string, metadataValue string, timeOut ...time.Duration) error {
	/*
	上传文件
	 */
	expireTime := append(timeOut, defaultTimeOut)[0]
	params := &s3.PutObjectInput{
		Bucket:aws.String(bucketName),
		Key:aws.String(keyRoute),
		Body:strings.NewReader(metadataValue),
	}
	ctx, cancle := context.WithTimeout(context.Background(), expireTime)
	defer cancle()
	_, err := r.ScS3Conn.PutObjectWithContext(ctx, params)
	return err
}

func (r *ScS3) SCGetObject(bucketName string, keyRoute string, timeOut ...time.Duration) (objDetail []byte, err error) {
	/*
	下载文件
	 */
	expireTime := append(timeOut, defaultTimeOut)[0]
	params := &s3.GetObjectInput{
		Bucket:aws.String(bucketName),
		Key:aws.String(keyRoute),
	}
	ctx, cancle := context.WithTimeout(context.Background(), expireTime)
	defer cancle()
	resp, err := r.ScS3Conn.GetObjectWithContext(ctx, params)
	if err != nil {
		return objDetail, err
	}
	body, err1 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body, err1
}
