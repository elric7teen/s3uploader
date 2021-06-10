package main

import (
	"os"
	"time"

	pkgAws "ghozi.com/prototype/s3uploader/pkg/aws"
	pkgS3 "ghozi.com/prototype/s3uploader/pkg/s3"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	logrus.SetReportCaller(true)
}

func main() {
	bucket := "bucket-name"
	filename := "test.txt"
	key := "folder/" + filename

	var client *pkgAws.Client
	if os.Getenv("MARIA_SVC_ENV") == "" {
		logrus.Info("set credentials using assumed role")
		client = pkgAws.NewClientWithAssumedRole(endpoints.ApSoutheast1RegionID, "aws account id", "role name")
	} else {
		logrus.Info("set credentials default credentials provider")
		client = pkgAws.NewClient(endpoints.ApSoutheast1RegionID)
	}
	// client := NewClient(nil, nil, nil)
	if !client.Validate(bucket) {
		panic("invalid client")
	}

	exp, err := client.Config.Credentials.ExpiresAt()
	if err != nil {
		panic(err)
	}

	logrus.Info("expired at :", exp)

	repo := pkgS3.NewRepo(client)

	for i := 1; i <= 30; i++ {
		if err := repo.Upload(filename, bucket, key); err != nil {
			logrus.Error(err)
			panic(err)
		}

		logrus.Info("success upload!")

		time.Sleep(1 * time.Minute)

		exp, err := client.Config.Credentials.ExpiresAt()
		if err != nil {
			panic(err)
		}
		logrus.Info("expired at :", exp)
	}
}
