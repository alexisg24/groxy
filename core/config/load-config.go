package config

import (
	"fmt"

	fileloader "github.com/alexisg24/groxy/core/file-loader"
)

var GlobalConfig fileloader.GlobalFileLoaderConfig

func Init() {
	err := GlobalConfig.Load("config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded config\n")
}
