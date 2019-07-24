package main

import (
	"flag"

	"github.com/rinosukmandityo/autorestore/helper"
)

func main() {

	var configLoc string
	flag.StringVar(&configLoc, "config", "configs/configs.json", "config file location")
	flag.Parse()

	configs := new(helper.DBConfig)
	helper.ReadJsonFile(configLoc, configs)

	helper.RestoreDBFromS3(configs.Database, configs.S3)

}
