package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	Session    *session.Session
	Downloader *s3manager.Downloader
	Bucket     string
	Region     string
	Key        string
}

type SeekParams struct {
	Position int
}

type CloudProvider interface{}

type SmartFile struct {
	CloudProvider CloudProvider
	SeekParams    SeekParams
}

func (s3 *S3) UpdateFileDetails(bucketKey string) {
	s := strings.SplitN(bucketKey, "/", 2)
	s3.Bucket = s[0]
	s3.Key = s[1]
	s3.Region = "us-east-1" // TODO: Figure out where to get this from
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3.Region),
	})
	s3.Session = sess
}

/*
func (sf *SmartFile) ReadFile(name string) ([]byte, error) {
	bucketKey := strings.Split(name, "://")[1]
	bucket := strings.Split(bucketKey, "/")[0]
	key := strings.Split(bucketKey, "/")[1]
	fmt.Println(bucket)
	buff := &aws.WriteAtBuffer{}

	downloader := s3manager.NewDownloader(sf.Session)
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	retBuf := buff.Bytes()
	return retBuf, err
}
*/

func Open(name string) (sf *SmartFile, err error) {
	var cp CloudProvider
	split := strings.Split(name, "://")
	remoteType := split[0]
	switch remoteType {
	case "s3":
		cp = &S3{}
		cp.(*S3).UpdateFileDetails(split[1])
		cp.(*S3).Downloader = s3manager.NewDownloader(cp.(*S3).Session)
		fmt.Println(reflect.TypeOf(cp))
	}
	sf = &SmartFile{
		CloudProvider: cp,
		SeekParams: SeekParams{
			Position: 0,
		},
	}
	return
}

func (sf *SmartFile) Read(b []byte) (n int64, err error) {
	//switch reflect.TypeOf(sf.CloudProvider) {
	//case *main.S3:
	inputObject := &s3.GetObjectInput{
		Bucket: aws.String(sf.CloudProvider.(*S3).Bucket),
		Key:    aws.String(sf.CloudProvider.(*S3).Key),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", sf.SeekParams.Position, sf.SeekParams.Position+cap(b)-1)),
	}
	fmt.Println(*inputObject.Range)
	buff := &aws.WriteAtBuffer{}
	n, err = sf.CloudProvider.(*S3).Downloader.Download(buff, inputObject)
	copy(b, buff.Bytes()[:n])
	sf.SeekParams.Position += int(n)
	return
	//}
	// Fix this to have a switch case for cloud providers

}

func (sf *SmartFile) Close() {
	sf = nil
}

func main() {
	item := "s3://ankit-storage/sumo_logic.txt"
	// Open
	sf, err := Open(item)
	defer sf.Close()
	fmt.Println(sf.CloudProvider, err)

	// Read
	b := make([]byte, 8)
	_, err = sf.Read(b)
	fmt.Println(b)
	b = make([]byte, 8)
	_, err = sf.Read(b)
	fmt.Println(b)
	b = make([]byte, 32)
	n, err := sf.Read(b)
	fmt.Println(b, n)
}
