# Auto Restore
To Restore MongoDB Database and Files into AWS S3 Bucket

How to use
---
#### Restore Database
It will download your MongoDB dump file with archive from AWS S3 Bucket and copy it into your local storage and restore it into your database  
It's located in `database` directory and following is brief description about the config file:
1. Change config file inside `configs/configs.json` with appropriate information about your database and S3 Bucket
2. `uri` is connection string used to connect into your database
3.	`archivename` is your archive name that will be used as prefix to get your archive file list from S3 Bucket
4.	`destpath` is destination directory for your dump files from S3 Bucket
5.	`retentionday` key use to set maximum object fetched from S3 bucket
6. If you set “retentionday: 0” the maximum object fetched will be 10
7.	`region` is your S3 Bucket region
8.	`bucket` is your S3 Bucket name
9.	`folder` is your folder in S3 Bucket, if any, and it will be used as prefix before “archivename” key
10. To test it just run `go run main.go`.
11.	To run it in scheduler just build it and add parameter `-config=yourconfigfilelocation`.

#### Restore Files
It will download all your files in a given S3 Bucket into your local storage
It's located in `file` directory and following is brief description about the config file:
1. Change config file inside `configs/configs.json` with appropriate information about your files directory location and S3 Bucket
2. `dirpath` is your download directory location
3.	`region` is your S3 Bucket region
4.	`bucket` is your S3 Bucket name
5. `folder` is your folder in S3 Bucket, if any, and it will be used as a prefix to get your all files from S3 Bucket
6. To test it just run `go run main.go`.
7.	To run it in scheduler just build it and add parameter `-config=yourconfigfilelocation`.