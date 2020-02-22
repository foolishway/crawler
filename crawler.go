package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gocolly/colly"
)

var once sync.Once

func main() {

	//创建收集器
	c := colly.NewCollector()

	//获取分页的url
	c.OnHTML(".toc a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnResponse(func(r *colly.Response) {
		_, filename := filepath.Split(r.Request.URL.Path)
		go SaveLocalFile(r.Body, filename)
	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	//开始访问https://research.swtch.com/
	c.Visit("https://research.swtch.com/")

}
func SaveLocalFile(body []byte, fileName string) error {
	once.Do(func() { Mkdir() })

	err := ioutil.WriteFile("./blogs/"+fileName+".html", body, 0777)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func Mkdir() error {
	dirName := "./blogs"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		log.Printf("Create dir error: %s\n", err)
		return err
	}
	err = os.Chmod(dirName, 0777)

	if err != nil {
		log.Printf("Change mod error: %s\n", err)
		return err
	}
	return nil
}
