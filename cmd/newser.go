package main

import (
	"log"

	"github.com/lnenad/newser/pkg"
)

func main() {
	config := pkg.GetConfig()

	log.Printf("Loaded configuration: %+v", config)

	pdf := pkg.SetupPdf()
	counter := 0

	if len(config.Defs.Website) > 0 {
		for _, websiteDefinition := range config.Defs.Website {
			if websiteDefinition.Disable == 1 {
				continue
			}
			pkg.WriteArticlesFromWebsite(config, websiteDefinition, pdf, &counter)
		}
	}

	savePath := pkg.GetSavePath(config.Output.Directory, config.Output.Extension)
	err := pdf.OutputFileAndClose(savePath)
	if err != nil {
		log.Fatal("Error while outputting pdf file", err)
	}
	log.Printf("Written %v articles to pdf at %v", counter, savePath)
}
