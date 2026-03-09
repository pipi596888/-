package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Mysql struct {
		DataSource string
	}
	RabbitMQ struct {
		Address string
	}
	Oss struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
	}
	Aliyun struct {
		AccessKeyId     string
		AccessKeySecret string
		AppKey          string
	}
	Tengxun struct {
		SecretId  string
		SecretKey string
		AppId     string
	}
}
