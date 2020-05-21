package awss3

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type AWSS3 struct {
}

func NewAWSS3() *AWSS3 {
	mod := &AWSS3{}

	return mod
}

func (m *AWSS3) Name() string {
	return "AWS s3"
}

func (m *AWSS3) Slug() string {
	return "aws-s3"
}

func (m *AWSS3) Description() string {
	return "Test if there is a AWS S3 bucket available on the hostname. Test this module with 'flaws.cloud'."
}

func checkBucketName(bucketName string) (bool, error) {
	pattern := "(^(([a-z0-9]|[a-z0-9][a-z0-9\\-]*[a-z0-9])\\.)*([a-z0-9]|[a-z0-9][a-z0-9\\-]*[a-z0-9])$)"
	bucketNameLen := len(bucketName)
	if bucketNameLen < 3 || bucketNameLen > 63 {
		return false, nil
	}

	patternRegExp, err := regexp.Compile(pattern)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot compile bucket regexp.")
		logrus.Warn(err)
		return false, err
	}

	return patternRegExp.MatchString(bucketName), nil
}

func (m *AWSS3) checkBucketWithoutCreds(ctx *common.ModuleContext, bucketName string) (bool, error) {
	// Make a simple GET request and check the status code (http to prevent SSL errors)
	url := "http://" + bucketName + ".s3.amazonaws.com"
	statusCode, _, _, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot GET the bucket without credentials.")
		logrus.Warn(err)
		return false, err
	}

	switch *statusCode {
	case 200:
		// we found a bucket !
		return true, nil
	case 403:
		// bucket exists but we need creds to list it
		return false, nil
	case 404:
		// this is not a bucket
		return false, nil
	case 503:
		// there is an error, maybe try later...
		return false, nil
	default:
		err = fmt.Errorf("AWS returns an unknow status code.")
		return false, err
	}
}

func (m *AWSS3) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}
	bucketName := hostname

	if strings.Contains(bucketName, ".amazonaws.com") {
		index := strings.LastIndex(bucketName, ".s3")
		if index > -1 {
			bucketName = bucketName[:index]
		}
	}

	isValidBucketName, err := checkBucketName(bucketName)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot check bucket name.")
		logrus.Warn(err)
		return err
	}

	if !isValidBucketName {
		err = fmt.Errorf("Bucket name is invalid.")
		logrus.Warn(err)
		return err
	}

	// TODO: check bucket with creds
	isBucket, err := m.checkBucketWithoutCreds(ctx, bucketName)
	if err != nil {
		logrus.Debug(err)
		return err
	}

	if isBucket {
		err = ctx.CreateNewAssetAsRaw("AWS S3 bucket is open")
		if err != nil {
			err = fmt.Errorf("Could not create vulnerability.")
			logrus.Warn(err)
			return err
		}

		err = ctx.AddTag("awd-s3")
		if err != nil {
			err = fmt.Errorf("Could not add tag to asset.")
			logrus.Warn(err)
			return err
		}
	}

	return nil
}
