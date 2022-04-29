package spider

import "github.com/gocolly/colly"

type Spider interface {
    Provider()
    Consumer()
}

type TaskInfo struct {
    TaskName string
}

func GetClient() *colly.Collector {
    c := colly.NewCollector()

    return c
}
