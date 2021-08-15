package concurrent

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	httpReq *http.Client
}

const (
	HeaderAcceptRanges      = "Accept-Ranges"
	HeaderAcceptRangesBytes = "bytes"

	HeaderRange = "Range"

	SplitRange    int64 = 10 * 1024 * 1024
	TimeFormatYmd       = "20060203"
)

type DownloadParam struct {
	Url      string
	FileName string
}

func NewClient() (cli *Client) {
	cli = &Client{}
	cli.httpReq = http.DefaultClient
	trans := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		MaxConnsPerHost:     20,
		IdleConnTimeout:     time.Minute,
	}
	cli.httpReq.Transport = trans
	return
}

func (c *Client) Download(ctx context.Context, param DownloadParam) (err error) {
	// 获取header
	headReq, err := http.NewRequest(http.MethodHead, param.Url, nil)
	if err != nil {
		return
	}
	resp, err := c.httpReq.Do(headReq)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("非成功状态码:" + strconv.Itoa(resp.StatusCode))
	}

	// 检查是否支持并发下载
	if resp.Header.Get(HeaderAcceptRanges) == HeaderAcceptRangesBytes {

		m := admin{
			id:      param.FileName,
			httpReq: c.httpReq,
			size:    resp.ContentLength,
			url:     param.Url,
		}
		err = m.execute(ctx)
		if err != nil {
			return err
		}
	}

	return
}
