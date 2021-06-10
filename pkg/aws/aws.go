package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	Sess   *session.Session
	Config *aws.Config
}

// NewClient will create new aws client
// credential will be generated automaticly by AWS SDK using one of credentials provider (environment variables, shared credentials file or IAM role)
// but in this case credential will not be refresh automatically by AWS if the session is expired.
// Need to refresh the credential manually
// learn more about AWS SDK Credentials : https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
// DO NOT USE THIS ON PROD. USE THIS ONLY ON LOCAL.
func NewClient(region string) *Client {

	// return no config and sesssion for empty string input
	if region == "" {
		return nil
	}

	sess := session.Must(session.NewSession())

	// create new config
	config := aws.NewConfig().
		WithRegion(region)

	return &Client{
		Sess:   sess,
		Config: config,
	}
}

// NewClientWithAssumedRole will create new aws client
// using credentials attached to EC2 instance (IAM role).
// credential will be refresh automatically by AWS
// if the session is expired.
// to check client session and config, use (*Client).Validate(bucket string)
func NewClientWithAssumedRole(
	region string,
	awsAccID string,
	role string) *Client {

	// return no config and sesssion for empty string inputs
	if awsAccID == "" || role == "" || region == "" {
		return nil
	}

	sess := session.Must(session.NewSession())

	// arn: Amazon Resource Name
	// more information about arn : https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
	arn := fmt.Sprintf(
		"arn:aws:iam::%s:role/%s",
		awsAccID,
		role,
	)

	// create new creds
	creds := stscreds.NewCredentials(sess, arn)

	// create new config
	config := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(region)

	return &Client{
		Sess:   sess,
		Config: config,
	}
}

// Validate client session and config credentials
// using list item on a bucket
func (c *Client) Validate(bucket string) bool {
	svc := s3.New(c.Sess, c.Config)

	_, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		log.Printf("Error test aws credentials : %v", err)
		return false
	}
	return true
}
