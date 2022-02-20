package pkg

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type Defs struct {
	Website []WebsiteDefinition `yaml: "website"`
}

type Output struct {
	Extension string `yaml: "extension"`
	Directory string `yaml: "directory"`
}

type Font struct {
	Title   int `yaml: "title"`
	Content int `yaml: "content"`
}

type Config struct {
	Defs   Defs   `yaml: "defs"`
	Font   Font   `yaml: "font"`
	Output Output `yaml: "output"`
}

func GetConfig() Config {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("No config.yaml found in root directory")
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error while parsing config: %v", err)
	}
	return config
}

func GetSavePath(directory, ext string) string {
	return filepath.Join(
		directory,
		time.Now().Format("2006-01-02-150405")+ext,
	)
}
