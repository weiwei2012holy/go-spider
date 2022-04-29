package spider

import (
    "fmt"
    "github.com/gocolly/colly"
    "log"
    "time"
)

type StudyGolang struct {
    Spider
    TaskInfo
    ListChan      chan string
    ListItemsChan chan string
}

type Article struct {
    Name    string
    Link    string
    Date    string
    Author  string
    View    string
    Comment string
    Like    string
    Tags    []string
}

var StudyGolangClient = new(StudyGolang)

func init() {
    StudyGolangClient.TaskInfo.TaskName = "GO语言中文网"
}

func (s *StudyGolang) Provider() {
    var url = "https://studygolang.com/articles"
    for i := 1; i < 20; i++ {
        time.Sleep(time.Second * 1)
        newUrl := fmt.Sprintf("%s?p=%d", url, i)
        s.ListChan <- newUrl
    }
}

func (s *StudyGolang) Consumer() {
    for url := range s.ListItemsChan {
        //url := <-ch
        c := GetClient()
        var articleList = make([]Article, 0)
        c.OnHTML(".article", func(element *colly.HTMLElement) {
            article := Article{}
            element.ForEach(".row", func(i int, element *colly.HTMLElement) {
                if i == 0 {
                    article.Name = element.ChildText("h2>a")
                    article.Link = element.ChildAttr("h2>a", "href")
                }
                if i == 1 {
                    article.Date = element.ChildText(".date")
                    article.Author = element.ChildText(".author")
                    article.View = element.ChildText(".view > span")
                    article.Comment = element.ChildText(".cmt > span")
                    article.Like = element.ChildText(".likenum")
                    article.Tags = element.ChildTexts("li>a")
                    //article.Tags = element.ChildAttrs("li>a", "title")
                }
            })
            articleList = append(articleList, article)

        })
        err := c.Visit(url)
        if err != nil {
            log.Fatal(err)
        }
        //log.Println(articleList)
        log.Println("处理完毕:>>", url)
    }

}
