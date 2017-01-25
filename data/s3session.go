package data

import (
	"fmt"
	"io"

	"bytes"

	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Session represents a connection to the project's AWS S3 bucket. Implements
// data.FileSystem interface.
type S3Session struct {
	session  *s3.S3
	protocol string
	hostname string
	bucket   string
}

// NewS3Session returns a new S3Session object connected to the project's S3
// bucket. If connection fails, then returns non-nil error.
func NewS3Session() (*S3Session, error) {
	// Attempt to establish connection to AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AWS S3")
	}
	return &S3Session{
		session:  s3.New(sess),
		protocol: "https://",
		hostname: "s3.amazonaws.com/",
		bucket:   "skilldirectory/",
	}, nil // successful connection
}

// Read returns an io.Reader that reads from the resource located at the
// specified path within the project's S3 bucket.
func (s *S3Session) Read(path string) (resource io.Reader, err error) {
	// Setup object to read from AWS
	params := &s3.GetObjectInput{
		Bucket: aws.String("skilldirectory"),
		Key:    aws.String(path),
	}

	// Try to read the resource/file from AWS
	result, err := s.session.GetObject(params)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource from AWS S3: %q: %s", path, err)
	}

	// Extract data from response body and return it
	bodyBytes, _ := ioutil.ReadAll(result.Body)
	result.Body.Close()
	return bytes.NewReader(bodyBytes), nil
}

// Write saves the specified resource to the project's S3 bucket under the
// specifed path.
func (s *S3Session) Write(path string, resource io.ReadSeeker) (url string,
	err error) {
	// Setup object to save in AWS
	params := &s3.PutObjectInput{
		Bucket: aws.String("skilldirectory"), // Required
		Key:    aws.String(path),             // Required
		Body:   resource,
	}

	// Try to save the resource/file to AWS
	_, err = s.session.PutObject(params)
	if err != nil {
		return "", fmt.Errorf("failed to save resource to AWS S3: %q: %s", path, err)
	}
	url = s.protocol + s.hostname + s.bucket + path
	return url, nil // Succesfully saved resource to S3 instance
}

// Delete removes the resource located at the specified path from the project's
// S3 bucket.
func (s *S3Session) Delete(path string) (err error) {
	// Setup object to delete from AWS
	params := &s3.DeleteObjectInput{
		Bucket: aws.String("skilldirectory"), // Required
		Key:    aws.String(path),             // Required
	}

	// Try to delete the resource/file from AWS
	_, err = s.session.DeleteObject(params)
	if err != nil {
		return fmt.Errorf("failed to delete resource from AWS S3: %q: %s", path, err)
	}
	return nil
}
