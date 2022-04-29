package spider

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "github.com/gocolly/colly"
    "log"
    url2 "net/url"
    "path"
    "regexp"
    "strings"
    "sync"
    "time"
)

type GameSky struct {
    Spider
    TaskInfo
    ListChan  chan GameSkyListItem
    Collector *colly.Collector
    ExistIds  map[int]bool
    Mu        *sync.Mutex
    //ListItemsChan chan string
}

var GameSkyClient = new(GameSky)

type GameSkyJsonData struct {
    Status     string `json:"status"`
    TotalPages int    `json:"totalPages"`
    Body       string `json:"body"`
}

type GameSkyListItem struct {
    Id       int
    Category string
    Name     string
    Link     string
    Cover    string
    Summary  string
    Date     string
    Author   string
    Comment  string
    Images   []string
}

func init() {
    GameSkyClient.TaskInfo.TaskName = "游民"
    //GameSkyClient.ListItemsChan = make(chan string)
    GameSkyClient.ListChan = make(chan GameSkyListItem, 1)
    GameSkyClient.Collector = GetClient()
    GameSkyClient.ExistIds = make(map[int]bool)
    GameSkyClient.Mu = &sync.Mutex{}
}

func (s *GameSky) Provider() {
    //todo ...
    var page = 1
    total := s.ParseListUrl(page)
    if total == 0 {
        log.Fatal("总数为空")
    }
    for i := page; i <= total; i++ {
        s.ParseListUrl(i)
    }
}

func (s *GameSky) Consumer() {
    //todo ...
    for item := range s.ListChan {
        newItem := item
        time.Sleep(time.Second)
        go func() {
            fmt.Println("Consumer=>>>", newItem)
            s.ParseDetailUrl(newItem)
        }()
    }
}

func (s *GameSky) ParseListUrlWap(page int) int {
    var data = make(map[string]interface{})

    data["type"] = "getwaplabelpage"
    data["isCache"] = true
    data["cacheTime"] = 60
    data["templatekey"] = "newshot"
    data["id"] = "1479167"
    data["nodeId"] = "20107"
    data["page"] = page

    jsonStr := JsonMarshal(data)

    jsonStr2 := "{\"type\":\"getwaplabelpage\",\"isCache\":true,\"cacheTime\":60,\"templatekey\":\"newshot\",\"id\":\"1479167\",\"nodeId\":\"21037\",\"page\":1}"
    fmt.Println(jsonStr)
    fmt.Println(jsonStr2)
    jqUrl := "https://db2.gamersky.com/LabelJsonpAjax.aspx"
    url := fmt.Sprintf("%s?jsondata=%s&_=%d", jqUrl, jsonStr, time.Now().UnixNano())

    var totalPage int
    s.Collector.OnResponse(func(response *colly.Response) {
        respJson := string(response.Body)
        respJson = strings.Trim(respJson, ";")
        respJson = strings.Trim(respJson, ")")
        respJson = strings.Trim(respJson, "(")

        var jsonData GameSkyJsonData
        JsonUnmarshal(respJson, &jsonData)
        if jsonData.Status != "ok" {
            log.Fatal("解析数据失败:", jsonData)
        }
        totalPage = jsonData.TotalPages

        doc, err := goquery.NewDocumentFromReader(strings.NewReader(jsonData.Body))
        if err != nil {
            log.Fatal(err)
        }
        var gameSkyList = make([]GameSkyListItem, 0)
        doc.Find("li").Each(func(i int, selection *goquery.Selection) {
            item := GameSkyListItem{}

            item.Category = selection.Find(".tit>.dh").First().Text()
            item.Name = CompressStr(selection.Find(".tit>.tt").First().Text())
            item.Link, _ = selection.Find(".tit>.tt").First().Attr("href")

            pu, _ := url2.Parse(item.Link)
            fmt.Println(pu)
            return
            //fmt.Println(pu)

            item.Cover, _ = selection.Find(".img> img").First().Attr("src")

            item.Summary = CompressStr(selection.Find(".con>.txt").First().Text())
            item.Date = selection.Find(".con>.time").First().Text()
            item.Comment = selection.Find(".con>.cy_comment").First().Text()

            //todo 写入数据到chan
            gameSkyList = append(gameSkyList, item)

            s.ListChan <- item

        })

        //fmt.Println(gameSkyList)

    })
    err := s.Collector.Visit(url)
    if err != nil {
        log.Fatal(err)
    }
    return totalPage
}

func (s *GameSky) CheckExist(id int) bool {
    defer s.Mu.Unlock()
    s.Mu.Lock()
    exist := s.ExistIds[id]
    if exist == false {
        s.ExistIds[id] = true
    }
    return exist
}

func (s *GameSky) ParseListUrl(page int) int {

    var data = make(map[string]interface{})
    //WEB
    data["type"] = "updatenodelabel"
    data["isCache"] = "true"
    data["cacheTime"] = 60
    data["nodeId"] = 20107
    data["page"] = page
    jsonStr := JsonMarshal(data)
    jqUrl := "https://db2.gamersky.com/LabelJsonpAjax.aspx"
    url := fmt.Sprintf("%s?jsondata=%s&_=%d", jqUrl, jsonStr, time.Now().UnixNano())

    var totalPage int
    s.Collector.OnResponse(func(response *colly.Response) {
        respJson := string(response.Body)
        respJson = strings.Trim(respJson, ";")
        respJson = strings.Trim(respJson, ")")
        respJson = strings.Trim(respJson, "(")

        var jsonData GameSkyJsonData
        JsonUnmarshal(respJson, &jsonData)
        if jsonData.Status != "ok" {
            log.Fatal("解析数据失败:", jsonData)
        }
        totalPage = jsonData.TotalPages

        doc, err := goquery.NewDocumentFromReader(strings.NewReader(jsonData.Body))
        if err != nil {
            log.Fatal(err)
        }
        var gameSkyList = make([]GameSkyListItem, 0)
        doc.Find("li").Each(func(i int, selection *goquery.Selection) {
            item := GameSkyListItem{}

            item.Category = selection.Find(".tit>.dh").First().Text()
            if item.Category != "趣味" {
                return
            }
            //过滤名称
            item.Name = CompressStr(selection.Find(".tit>.tt").First().Text())
            if !strings.Contains(item.Name, "动态图") && !strings.Contains(item.Name, "囧图") {
                return
            }
            item.Link, _ = selection.Find(".tit>.tt").First().Attr("href")
            //过滤已经处理或者存在的数据
            linkInfo := FindNumbers(item.Link)
            if len(linkInfo) > 1 {
                item.Id = linkInfo[len(linkInfo)-1]
            }
            if s.CheckExist(item.Id) {
                return
            }

            //item.Cover, _ = selection.Find(".img > * > img").First().Attr("src")
            item.Cover, _ = selection.Find(".pe_u_thumb").First().Attr("src")

            item.Summary = CompressStr(selection.Find(".txt").First().Text())
            item.Date = selection.Find(".time").First().Text()
            item.Comment = selection.Find(".pls").First().Text()

            //todo 写入数据到chan
            gameSkyList = append(gameSkyList, item)
            s.ListChan <- item

            //s.ParseDetailUrl(item)
        })

        //fmt.Println(gameSkyList)

    })
    err := s.Collector.Visit(url)
    if err != nil {
        log.Fatal(err)
    }
    return totalPage
}

func (s *GameSky) ParseDetailUrl(item GameSkyListItem) {

    c := GetClient()

    //s.Collector.OnHTML("", func(element *colly.HTMLElement) {
    //
    //})
    max := 100
    maxErrorCount := 5
    errorCount := 0
    //fmt.Println(item)

    item.Images = make([]string, 0)

    dirName := fmt.Sprintf("[%s][%d]%s", item.Date, item.Id, item.Name)
    for i := 2; i < max; i++ {
        pageUrl := s.BuildDetailPageUrl(item.Link, i)
        fmt.Println(pageUrl)

        c.OnError(func(response *colly.Response, err error) {
            fmt.Println(err)
            errorCount++
            if errorCount >= maxErrorCount {
                max = i
            }
        })
        c.OnHTML(".Mid2L_con", func(element *colly.HTMLElement) {

            element.ForEach(".GsImageLabel", func(i int, element *colly.HTMLElement) {

                imgName := element.Text
                imgList := element.ChildAttrs(".picact", "src")

                for i2, s2 := range imgList {
                    imgNameNew := fmt.Sprintf("%s-%d-%s-%s", imgName, i2+1, Md5(s2), path.Base(s2))
                    err := DownloadImage(s2, dirName, imgNameNew)
                    if err != nil {
                        fmt.Println(fmt.Sprintf("下载失败(%s) ==> %s", err, s2))
                    }
                }

                //fmt.Println(imgName, imgList)

                //if img != "" {
                //    item.Images = append(item.Images, img)
                //}

            })
            //fmt.Println(element.Attr("src"))
            //element.ForEach(".picact", func(i int, element *colly.HTMLElement) {
            //    img := element.Attr("src")
            //    if img != "" {
            //        item.Images = append(item.Images, img)
            //    }
            //})
            //fmt.Println(item.Images)

            //time.Sleep(time.Second * 5)

        })

        c.Visit(pageUrl)
        //fmt.Println(pageUrl)
    }

    //err := s.Collector.Visit(item.Link)
    //if err != nil {
    //    log.Fatal(err)
    //}
}

func (s GameSky) BuildDetailPageUrl(url string, page int) string {
    if page == 1 {
        return url
    }
    reg := regexp.MustCompile(".*/(\\d+).shtml")

    regRes := reg.FindStringSubmatch(url)
    //fmt.Println(regRes)

    if len(regRes) > 0 {
        id := regRes[len(regRes)-1]
        path := regRes[len(regRes)-2]
        return strings.Replace(path, id, fmt.Sprintf("%s_%d", id, page), -1)
    }
    return ""
}
