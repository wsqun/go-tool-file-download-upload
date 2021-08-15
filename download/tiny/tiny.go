package tiny

import (
	"io/ioutil"
	"net/http"
	"os"
)

type FileParam struct {
	Url string
	DestFile string
}

type Client struct {

}

func (s *Client) Download(param FileParam)  {
	// 下载
	resp, err := http.Get(param.Url)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return
	}
	println("content-length:", resp.ContentLength)
	println("content-type:", resp.Header.Get("Content-Type"))

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(param.DestFile, body, os.FileMode(0666))
	if err != nil {
		return
	}
	println("write ok")
}
