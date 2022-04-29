package spider

import (
    "crypto/md5"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
)

const (
    ImageDownloadDir = "images"
)

func JsonMarshal(v interface{}) string {
    bytes, err := json.Marshal(v)
    if err != nil {
        fmt.Println("json序列化：", err)
    }
    return string(bytes)
}

func JsonUnmarshal(str string, v interface{}) {
    err := json.Unmarshal([]byte(str), v)
    if err != nil {
        fmt.Println("json反序列化：", err)
    }
    return
}

func CompressStr(str string) string {
    if str == "" {
        return ""
    }
    //匹配一个或多个空白符的正则表达式
    reg := regexp.MustCompile("\\s+")
    return reg.ReplaceAllString(str, "")
}

func FindNumbers(str string) []int {
    var numbers = make([]int, 0)
    reg := regexp.MustCompile("[0-9]+")
    regRes := reg.FindAllString(str, -1)
    for _, m := range regRes {
        nm, _ := strconv.Atoi(m)
        numbers = append(numbers, nm)
    }
    return numbers
}

func ImageDownloadPath(path string) (string, error) {
    appPath, err := os.Getwd() ///Users/yuwei/go/src/local.test/spider
    if err != nil {
        return "", err
    }
    imgDir := filepath.Join(appPath, ImageDownloadDir, path)
    err = MakeDir(imgDir)
    if err != nil {
        return "", err
    }
    return imgDir, nil
}

func DownloadImage(url string, dir string, name string) error {
    imgPath, err := ImageDownloadPath(dir)
    if err != nil {
        return errors.New(fmt.Sprintf("初始化下载目录失败:%s", err))
    }
    filePath := filepath.Join(imgPath, name)

    if FileExists(filePath) {
        return errors.New("文件已经存在")
    }

    urlRes, err := http.Get(url)
    if err != nil {
        return errors.New(fmt.Sprintf("下载图片失败:%s", err))
    }
    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            fmt.Println(fmt.Sprintf("io.ReadCloser ERROR:%s", err))
        }
    }(urlRes.Body)
    b, err := io.ReadAll(urlRes.Body)
    if err != nil {
        return errors.New(fmt.Sprintf("读取图片失败:%s", err))
    }
    f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
    if err != nil {
        return errors.New(fmt.Sprintf("初始化文件失败:%s", err))
    }
    _, err = f.Write(b)
    if err != nil {
        return errors.New(fmt.Sprintf("写入文件失败:%s", err))
    }
    return nil
}

func MakeDir(dir string) error {
    if FileExists(dir) {
        return nil
    }
    err := os.MkdirAll(dir, 0711)
    if err != nil {
        return err
    }
    return nil
}

func FileExists(file string) bool {
    if _, err := os.Stat(file); os.IsNotExist(err) {
        return false
    } else {
        return true
    }
}

func Md5(str string) string {
    h := md5.New()
    io.WriteString(h, str)
    return fmt.Sprintf("%x", h.Sum(nil))
}
