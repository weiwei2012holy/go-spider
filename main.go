package main

import (
    "fmt"
    "local.test/spider/spider"
    "log"
)

//colly 文档：http://go-colly.org/docs/introduction/start/
//HTML解析规则：https://github.com/PuerkitoBio/goquery

func main() {

    clientList := RegisterClient()
    //todo 交互式命令行获取需要运行的脚本，后续可以扩展更多参数
    clientName := spider.GameSkyClient.TaskName
    var client Client

    for _, c := range clientList {
        if c.Name == clientName {
            client = c
            break
        }
    }
    if client.Client == nil {
        log.Fatal("无效的Client")
    }

    fmt.Println("任务启动中。。。")
    go client.Client.Provider()
    client.Client.Consumer()
}

type Client struct {
    Name   string
    Client spider.Spider
}

// RegisterClient 注册服务
func RegisterClient() []Client {
    var client = []Client{
        {
            Name:   spider.StudyGolangClient.TaskName,
            Client: spider.StudyGolangClient,
        },
        {
            Name:   spider.GameSkyClient.TaskName,
            Client: spider.GameSkyClient,
        },
    }
    return client

}
