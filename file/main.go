package main

import (
	"flag"
	"path/filepath"

	"github.com/rinosukmandityo/autorestore/helper"
)

func main() {

	var configLoc string
	flag.StringVar(&configLoc, "config", filepath.Join(helper.WD, "configs", "configs.json"), "config file location")
	flag.Parse()

	configs := new(helper.FileConfig)
	helper.ReadJsonFile(configLoc, configs)

	helper.DownloadObjectsFromS3(configs.File, configs.S3)
}
