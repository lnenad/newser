package pkg

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/gocolly/colly/v2"
)

type WebsiteDefinition struct {
	Index                    string   `yaml:"index,omitempty"`
	IndexSelector            string   `yaml:"indexSelector,omitempty"`
	TitleSelector            string   `yaml:"titleSelector,omitempty"`
	LinkSelector             string   `yaml:"linkSelector,omitempty"`
	LinkAttr                 string   `yaml:"linkAttr,omitempty"`
	LinkPrefix               string   `yaml:"linkPrefix,omitempty"`
	ArticleContainerSelector string   `yaml:"articleContainerSelector,omitempty"`
	ArticleContentSelector   string   `yaml:"articleContentSelector,omitempty"`
	IgnoreString             string   `yaml:"ignoreString,omitempty"`
	RemoveElems              []string `yaml:"removeElems,omitempty"`
	CollectOnly              int      `yaml:"collectOnly,omitempty"`
	Disable                  int      `yaml:"disable,omitempty"`
}

func WriteArticlesFromWebsite(config Config, wd WebsiteDefinition, pdf fpdf.Pdf, numArticles *int) {
	articlesChan := make(chan Article)

	crawlWebsite(
		wd,
		articlesChan,
	)

	WriteHeader(pdf, wd.Index)

	messages := 0

	for {
		select {
		case article := <-articlesChan:
			if article.Content == "" {
				continue
			}
			messages++
			(*numArticles)++
			//registerImage(pdf, article.img)
			WriteArticle(config, pdf, article)

			if wd.CollectOnly > 0 && messages == wd.CollectOnly {
				return
			}
			//fmt.Println("received", article)
		case <-time.After(5 * time.Second):
			fmt.Println("Crawling completed, timeout after 5 seconds")
			return
		}
	}
}

func crawlWebsite(wd WebsiteDefinition, articlesChan chan Article) {
	c := colly.NewCollector()
	c.OnHTML(wd.IndexSelector, func(e *colly.HTMLElement) {
		imgs := e.ChildAttrs("a > picture > source", "srcset")
		title := e.ChildText(wd.TitleSelector)
		link := ""
		if wd.LinkSelector != "" {
			links := e.ChildAttrs(wd.LinkSelector, wd.LinkAttr)
			if len(links) > 0 {
				link = links[0]
			}
		} else {
			link = e.Attr("href")
		}
		if strings.Trim(title, " ") == "" {
			return
		}
		if link == "" {
			log.Fatalf("Invalid link selection, no links for fetching content found for website: %v, title: %v", wd.Index, title)
		}
		if wd.LinkPrefix != "" {
			link = wd.LinkPrefix + link
		}
		article := Article{
			Title: string(title),
			Img:   getImage(imgs),
			Link:  link,
		}

		contentChan := make(chan string, 1)
		go getContent(article, wd, contentChan)
		article.Content = <-contentChan
		articlesChan <- article
	})
	go c.Visit(wd.Index)
}

func getImage(imageSets []string) string {
	if len(imageSets) < 1 {
		return ""
	}

	images := strings.Split(imageSets[0], ",")

	return strings.Split(images[0], " ")[0]
}

var whitespaces = regexp.MustCompile(`\n\s+`)

func getContent(a Article, wd WebsiteDefinition, contentChan chan string) {
	c := colly.NewCollector()
	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Error while fetching article content for url: %v\nError: %v", a.Link, e)
	})
	c.OnHTML(wd.ArticleContainerSelector, func(e *colly.HTMLElement) {
		if len(wd.RemoveElems) > 0 {
			children := e.DOM.Children()
			for _, re := range wd.RemoveElems {
				children.Find(re).Remove()
			}
		}
		content := whitespaces.ReplaceAllString(e.ChildText(wd.ArticleContentSelector), "\n")
		if wd.IgnoreString != "" && strings.Contains(content, wd.IgnoreString) {
			contentChan <- ""
		}
		contentChan <- content
	})
	c.Visit(a.Link)
}
