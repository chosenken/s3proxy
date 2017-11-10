package main

import (
	"flag"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
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
	// File should be in the format of "/<bucket name>/<key to file>"
	bucket := ctx.Param("bucket")
	key := ctx.Param("key")
	if len(bucket) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "S3 File must be specified"})
		return
	}
	sess, err := session.NewSession()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	s3svc := s3.New(sess)
	req, _ := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	})

	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
	}

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
