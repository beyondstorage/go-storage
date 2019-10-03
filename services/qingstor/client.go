package qingstor

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yunify/qingstor-sdk-go/v3/config"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/define"
)

// Client is the qingstor object sotrage client.
//
//go:generate go run ../../internal/cmd/meta_gen/main.go
type Client struct {
	config  *config.Config
	service *service.Service
	bucket  *service.Bucket
}

// SetupBucket will setup bucket for client.
func (c *Client) SetupBucket(bucketName, zoneName string) (err error) {
	errorMessage := "setup qingstor bucket failed: %w"

	if zoneName != "" {
		bucket, err := c.service.Bucket(bucketName, zoneName)
		if err != nil {
			return fmt.Errorf(errorMessage, err)
		}
		c.bucket = bucket
		return nil
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s://%s.%s:%d", c.config.Protocol, bucketName, c.config.Host, c.config.Port)

	r, err := client.Head(url)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	if r.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("head status is %d instead of %d", r.StatusCode, http.StatusTemporaryRedirect)
		return fmt.Errorf(errorMessage, err)
	}

	// Example URL: https://bucket.zone.qingstor.com
	zoneName = strings.Split(r.Header.Get("Location"), ".")[1]
	bucket, err := c.service.Bucket(bucketName, zoneName)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	c.bucket = bucket
	return
}

// NewFromConfig will create a new client from config.
func NewFromConfig(cfg *config.Config) (*Client, error) {
	errorMessage := "create new qingstor client from config failed: %w"

	srv, err := service.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return &Client{service: srv}, nil
}

// NewFromHomeConfigFile will create a new client from default home config file.
func NewFromHomeConfigFile() (*Client, error) {
	errorMessage := "create new qingstor client from home config file failed: %w"

	cfg, err := config.NewDefault()
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return NewFromConfig(cfg)
}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, option ...define.Option) (i define.Informer, err error) {
	panic("implement me")
}

// Delete implements Storager.Delete
func (c *Client) Delete(path string, option ...define.Option) (err error) {
	panic("implement me")
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, option ...define.Option) (err error) {
	panic("implement me")
}

// Move implements Storager.Move
func (c *Client) Move(src, dst string, option ...define.Option) (err error) {
	panic("implement me")
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, option ...define.Option) (dir chan *define.Dir, file chan *define.File, stream chan *define.Stream, err error) {
	panic("implement me")
}

// ReadFile implements Storager.ReadFile
func (c *Client) ReadFile(path string, option ...define.Option) (r io.ReadCloser, err error) {
	errorMessage := "qingstor ReadFile failed: %w"

	_ = parseOptionReadFile(option...)
	input := &service.GetObjectInput{}

	output, err := c.bucket.GetObject(path, input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return output.Body, nil
}

// WriteFile implements Storager.WriteFile
func (c *Client) WriteFile(path string, size int64, r io.ReadCloser, option ...define.Option) (err error) {
	errorMessage := "qingstor WriteFile failed: %w"

	defer r.Close()

	opts := parseOptionWriteFile(option...)
	input := &service.PutObjectInput{
		ContentLength: &size,
		Body:          r,
	}
	if opts.HasMd5 {
		input.ContentMD5 = &opts.Md5
	}
	if opts.HasStorageClass {
		input.XQSStorageClass = &opts.StorageClass
	}

	_, err = c.bucket.PutObject(path, input)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}
	return nil
}

// ReadStream implements Storager.ReadStream
func (c *Client) ReadStream(path string, option ...define.Option) (r io.ReadCloser, err error) {
	panic("not supported")
}

// WriteStream implements Storager.WriteStream
func (c *Client) WriteStream(path string, r io.ReadCloser, option ...define.Option) (err error) {
	panic("not supported")
}

// InitSegment implements Storager.InitSegment
func (c *Client) InitSegment(path string, size int64, option ...define.Option) (err error) {
	panic("implement me")
}

// ReadSegment implements Storager.ReadSegment
func (c *Client) ReadSegment(path string, offset, size int64, option ...define.Option) (r io.ReadCloser, err error) {
	panic("implement me")
}

// WriteSegment implements Storager.WriteSegment
func (c *Client) WriteSegment(path string, offset, size int64, r io.ReadCloser, option ...define.Option) (err error) {
	panic("implement me")
}

// CompleteSegment implements Storager.CompleteSegment
func (c *Client) CompleteSegment(path string, option ...define.Option) (err error) {
	panic("implement me")
}

// AbortSegment implements Storager.AbortSegment
func (c *Client) AbortSegment(path string, option ...define.Option) (err error) {
	panic("implement me")
}
