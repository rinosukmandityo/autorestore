package helper

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func RestoreDBFromS3(dbconfig DBDetail, s3config S3Detail) {
	s3config.KeyPrefix = dbconfig.ArchiveName
	result, e := GetListObjectsV2WithContext(s3config, int64(dbconfig.RetentionDay))
	if e != nil {
		log.Println(e.Error())
	}

	latestModified, latestKey := time.Time{}, ""
	for _, res := range result.Contents {
		key := strings.TrimLeft(*res.Key, s3config.Folder+"/")
		if !strings.HasSuffix(key, "/") && res.LastModified.After(latestModified) {
			latestKey = key
			latestModified = *res.LastModified
		}
	}
	s3config.Key = latestKey
	obj, _ := GetObjectWithContext(s3config)
	bodyData := make([]byte, *obj.ContentLength)
	obj.Body.Read(bodyData)
	if e = WriteFile(bodyData, dbconfig.DestinationPath, latestKey); e != nil {
		log.Println(e.Error())
	}
	RestoreDB(dbconfig, filepath.Join(dbconfig.DestinationPath, latestKey))
	log.Println("database has been restored")
}

func RestoreDB(dbconfig DBDetail, archiveName string) {
	ExecCommand([]string{"/C", "mongorestore", "--uri", dbconfig.URI, "--archive=" + archiveName, "--drop"})
}

func GetBucketPathFromConfig(s3config S3Detail) (bucket string) {
	bucket = strings.Trim(s3config.Bucket, "/")
	folder := strings.Trim(s3config.Folder, "/")
	if folder != "" {
		bucket = fmt.Sprintf("%s/%s/", bucket, folder)
	}
	return
}
