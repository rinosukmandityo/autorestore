package helper

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// SyncBucketIterator is used to download a given bucket to local storage.
type SyncBucketIterator struct {
	bucket    string
	fileInfos []fileInfo
	err       error
}

type fileInfo struct {
	key      string
	fullpath string
}

// NewSyncWalkBucket will walk the path, and store the key and full path
// of the object to be downloaded. This will return a new SyncBucketIterator
// with the data provided from walking the path.
func NewSyncWalkBucket(fileconfig FileDetail, objects *s3.ListObjectsV2Output, iter *SyncBucketIterator) {
	dirpath := fileconfig.DirectoryPath
	fInfos := []fileInfo{}

	for _, obj := range objects.Contents {
		key := *obj.Key
		if !strings.HasSuffix(key, "/") {
			pathJoin := filepath.Join(dirpath, key)          // fullpath in local, combination of full directory path and file name
			fpath, _ := filepath.Split(pathJoin)             // split between directory and file
			os.MkdirAll(fpath, 0777)                         // create directory on local
			fInfos = append(fInfos, fileInfo{key, pathJoin}) // put key (with its directory) and fullpath in local
		}
	}
	iter.fileInfos = append(iter.fileInfos, fInfos...)
}

// Next will determine whether or not there is any remaining files to
// be uploaded.
func (iter *SyncBucketIterator) Next() bool {
	return len(iter.fileInfos) > 0
}

// Err returns any error when os.Open is called.
func (iter *SyncBucketIterator) Err() error {
	return iter.err
}

// UploadObject will prep the new download object by open that file and constructing a new
// s3manager.DownloadInput.
func (iter *SyncBucketIterator) DownloadObject() s3manager.BatchDownloadObject {
	fi := iter.fileInfos[0]                                            // get first element
	iter.fileInfos = iter.fileInfos[1:]                                // put the rest of element
	body, err := os.OpenFile(fi.fullpath, os.O_RDWR|os.O_CREATE, 0777) // open file with full path to provide write operation
	if err != nil {
		iter.err = err
	}

	input := s3.GetObjectInput{
		Bucket: &iter.bucket,
		Key:    &fi.key,
	}

	return s3manager.BatchDownloadObject{
		Object: &input,
		Writer: body,
	}
}
