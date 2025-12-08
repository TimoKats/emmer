package server

import (
	"bytes"
	"context"
	"encoding/json"
	gio "io"
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

func (io S3Fs) formatS3Key(filename string) string {
	return strings.ReplaceAll(filename, "--", "/")
}

func (io S3Fs) Put(filename string, value any) error {
	// convert data to JSON
	filename = io.formatS3Key(filename)
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// upload JSON to S3
	_, err = io.client.PutObject(io.ctx, &s3.PutObjectInput{
		Bucket:      aws.String(io.bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(jsonBytes),
		ContentType: aws.String("application/json"),
	})
	return err
}

func (io S3Fs) Get(filename string) (map[string]any, error) {
	// fetch the object from S3
	data := make(map[string]any)
	filename = io.formatS3Key(filename)
	resp, err := io.client.GetObject(io.ctx, &s3.GetObjectInput{
		Bucket: aws.String(io.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return data, err
	}
	// read the object body
	defer resp.Body.Close()                  //nolint:errcheck
	bodyBytes, err := gio.ReadAll(resp.Body) //nolint:errcheck
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bodyBytes, &data)
	return data, err
}

func (io S3Fs) Del(filename string) error {
	filename = io.formatS3Key(filename)
	_, err := io.client.DeleteObject(io.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(io.bucket),
		Key:    aws.String(filename),
	})
	if err == nil {
		log.Printf("succesfully removed %s from %s", filename, io.bucket)
	}
	return err
}

func (io S3Fs) Ls() ([]string, error) {
	resp, err := io.client.ListObjectsV2(io.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(io.bucket),
	})
	if err == nil {
		keys := []string{}
		for _, item := range resp.Contents {
			keys = append(keys, *item.Key)
		}
		return keys, nil
	}
	return []string{}, err
}

func SetupS3() *S3Fs {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	bucket := os.Getenv("EM_S3_BUCKET")
	if err != nil || len(bucket) == 0 {
		log.Panicf("can't setup S3: %s", err.Error())
	}
	return &S3Fs{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		ctx:    ctx,
	}
}
