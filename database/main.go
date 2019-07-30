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

	configs := new(helper.DBConfig)
	helper.ReadJsonFile(configLoc, configs)

	helper.RestoreDBFromS3(configs.Database, configs.S3)

}
