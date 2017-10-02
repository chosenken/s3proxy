package main

import (
	"context"
	"flag"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	log   = logrus.WithField("package", "main")
	port  = flag.String("port", "4041", "Port to listen on")
	debug = flag.Bool("debug", false, "Enable debug logs")
)

func main() {
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logger := log.WithField("function", "main")
	logger.Info("Starting S3 CMD")
	listenAndServe()
}

func listenAndServe() {
	router := gin.Default()
	router.GET("/:bucket/*key", getS3File)
	router.Run(":" + *port)
}

func getS3File(ctx *gin.Context) {
	logger := log.WithField("function", "getS3File")
	// File should be in the format of "/<bucket name>/<key to file>"
	bucket := ctx.Param("bucket")
	key := ctx.Param("key")
	keySplit := strings.SplitAfter(key, "/")
	fileName := keySplit[len(keySplit)-1]
	if len(bucket) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "S3 File must be specified"})
		return
	}
	numBytes, buff, err := downloadS3File(&bucket, &key, ctx)
	if err != nil {
		logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Length", strconv.FormatInt(numBytes, 10))
	logger.WithFields(logrus.Fields{"fileName": fileName, "numBytes": numBytes}).Debug("Returning file")
	ctx.Data(http.StatusOK, "application/octet-stream", buff)
}

func downloadS3File(bucket, key *string, ctx context.Context) (int64, []byte, error) {
	logger := log.WithField("function", "downloadS3File")
	logger.WithFields(logrus.Fields{"bucket": bucket, "key": key}).Debug("Downloading S3 File")
	s3Downloader := s3manager.NewDownloader(session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	buff := &aws.WriteAtBuffer{}
	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	numBytes, err := s3Downloader.DownloadWithContext(ctx, buff, input)
	if err != nil {
		return 0, nil, err
	}
	return numBytes, buff.Bytes(), nil
}
