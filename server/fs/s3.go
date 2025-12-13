package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Fs struct {
	client *s3.Client
	bucket string
	ctx    context.Context
}

// format filename to s3 key (replace -- with /) to deal with s3 folders
func (fs S3Fs) formatS3Key(filename string) string {
	return strings.ReplaceAll(filename, "--", "/")
}

// creates empty (or prefilled) JSON file at path
func (fs S3Fs) Put(filename string, value any) error {
	// convert data to JSON
	filename = fs.formatS3Key(filename)
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// upload JSON to S3
	_, err = fs.client.PutObject(fs.ctx, &s3.PutObjectInput{
		Bucket:      aws.String(fs.bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(jsonBytes),
		ContentType: aws.String("application/json"),
	})
	return err
}

// reads JSON file into map[string]any variable
func (fs S3Fs) Get(filename string) (map[string]any, error) {
	// fetch the object from S3
	data := make(map[string]any)
	filename = fs.formatS3Key(filename)
	resp, err := fs.client.GetObject(fs.ctx, &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return data, err
	}
	// read the object body
	defer resp.Body.Close()                 //nolint:errcheck
	bodyBytes, err := io.ReadAll(resp.Body) //nolint:errcheck
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bodyBytes, &data)
	return data, err
}

// removes entire JSON file
func (fs S3Fs) Del(filename string) error {
	filename = fs.formatS3Key(filename)
	_, err := fs.client.DeleteObject(fs.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(filename),
	})
	if err == nil {
		log.Printf("succesfully removed %s from %s", filename, fs.bucket)
	}
	return err
}

// list files in s3 bucket
func (fs S3Fs) Ls() ([]string, error) {
	// get all objects
	resp, err := fs.client.ListObjectsV2(fs.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(fs.bucket),
	})
	// format as list of strings
	if err == nil {
		keys := []string{}
		for _, item := range resp.Contents {
			keys = append(keys, *item.Key)
		}
		return keys, nil
	}
	return []string{}, err
}

// creates new S3Fs instance with settings applied
func SetupS3() *S3Fs {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	bucket := os.Getenv("EM_S3_BUCKET")
	if err != nil || len(bucket) == 0 {
		log.Panicf("can't setup S3: %s", err.Error())
	}
	log.Printf("selected S3 fs in %s", bucket)
	return &S3Fs{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		ctx:    ctx,
	}
}
