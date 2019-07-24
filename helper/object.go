package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GetListObjectsV2WithContext(s3config S3Detail, maxkeys int64) (result *s3.ListObjectsV2Output, e error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config.Region),
	}))
	svc, ctx := s3.New(sess), context.Background()

	if maxkeys < 1 {
		maxkeys = 10
	} else {
		maxkeys++
	}

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3config.Bucket),
		MaxKeys: aws.Int64(int64(maxkeys)),
		Prefix:  aws.String(s3config.Folder + "/" + s3config.KeyPrefix),
	}

	result, e = svc.ListObjectsV2WithContext(ctx, input)
	if e != nil {
		if aerr, ok := e.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(e.Error())
		}
		return
	}
	return
}

func GetObjectWithContext(s3config S3Detail) (result *s3.GetObjectOutput, e error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config.Region),
	}))
	svc, ctx := s3.New(sess), context.Background()

	input := &s3.GetObjectInput{
		Bucket: aws.String(GetBucketPathFromConfig(s3config)),
		Key:    aws.String(s3config.Key),
	}

	result, e = svc.GetObjectWithContext(ctx, input)
	if e != nil {
		if aerr, ok := e.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(e.Error())
		}
		return
	}
	return
}

func DownloadObjectsFromS3(fileconfig FileDetail, s3config S3Detail) {
	sess := session.New(&aws.Config{
		Region: aws.String(s3config.Region),
	})
	downloader := s3manager.NewDownloader(sess)

	iter := new(SyncBucketIterator)
	iter.bucket = s3config.Bucket

	query := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3config.Bucket),
		Prefix: aws.String(s3config.Folder),
	}
	svc := s3.New(sess)

	// Flag used to check if we need to go further
	truncatedListing := true

	for truncatedListing {
		resp, err := svc.ListObjectsV2WithContext(context.Background(), query)

		if err != nil {
			// Print the error.
			fmt.Println(err.Error())
			return
		}
		// Get all files
		NewSyncWalBucket(fileconfig, resp, iter)
		// Set continuation token
		query.ContinuationToken = resp.NextContinuationToken
		truncatedListing = *resp.IsTruncated
	}

	if err := downloader.DownloadWithIterator(aws.BackgroundContext(), iter); err != nil {
		log.Printf("unexpected error has occurred: %v", err)
	}

	if err := iter.Err(); err != nil {
		log.Printf("unexpected error occurred during file walking: %v", err)
	}

	log.Println("Download Objects from S3 Bucket success")
}
