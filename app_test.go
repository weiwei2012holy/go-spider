package main

import (
    "github.com/stretchr/testify/assert"
    "local.test/spider/spider"
    "os"
    "testing"
)

func TestHelper(t *testing.T) {

    t.Run("DownloadImage", func(t *testing.T) {
        path, err := spider.ImageDownloadPath("")
        assert.NoError(t, err)
        err = os.RemoveAll(path)
        assert.NoError(t, err)

        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "", "test-1.gif")
        assert.NoError(t, err)
        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "/", "test-2.gif")
        assert.NoError(t, err)
        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "test", "test-3.gif")
        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "test/", "test-3-1.gif")
        assert.NoError(t, err)
        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "/test", "test-4.gif")
        assert.NoError(t, err)
        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "/test/", "test-5.gif")
        assert.NoError(t, err)
        err = spider.DownloadImage("http://pic.rmb.bdstatic.com/bjh/down/3271c2cd9d95d7bea755c8aa9b3639bc.gif", "/test/test2/test3", "test-5-1.gif")
        assert.NoError(t, err)
    })

}
