package concurrent

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type admin struct {
	url     string
	workers []*worker
	id      string
	httpReq *http.Client
	size    int64
	wg      *sync.WaitGroup
	dir     string
	filePN  string
}

type worker struct {
	index    int
	start    int64
	end      int64
	fileName string
}

func (s *admin) execute(ctx context.Context) (err error) {
	s.wg = &sync.WaitGroup{}
	pwd, _ := os.Getwd()
	s.dir = pwd + "/" + time.Now().Format(TimeFormatYmd)
	// 创建目录
	_, err = os.Stat(s.dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(s.dir, os.ModePerm)
	}
	if err != nil {
		return
	}
	group := new(errgroup.Group)

	var i = 0
	var start int64 = 0
	for {
		s.wg.Add(1)
		w := &worker{
			index:    i,
			start:    start,
			end:      start + SplitRange,
			fileName: s.id + "_" + strconv.Itoa(i),
		}
		start = w.end + 1
		i++
		group.Go(func() error {
			return s.download(w)
		})
		s.workers = append(s.workers, w)

		if start >= s.size {
			break
		}
	}
	if err = group.Wait(); err != nil {
		return
	}

	// 组合文件块
	return s.merge(ctx)
}

func (s *admin) merge(ctx context.Context) (err error) {
	fileName := s.dir + "/" + s.id
	destFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer destFile.Close()

	for _, w := range s.workers {
		partFileName := s.dir + "/" + w.fileName
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		io.Copy(destFile, partFile)
		partFile.Close()
		os.Remove(partFileName)
	}
	s.filePN = fileName
	return
}

func (s *admin) download(worker *worker) (err error) {
	defer s.wg.Done()
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
	// 分组下载
	reqDl, _ := http.NewRequest(http.MethodGet, s.url, nil)
	reqDl.Header.Set(HeaderRange, fmt.Sprintf("bytes=%d-%d", worker.start, worker.end))
	respDl, err := s.httpReq.Do(reqDl)
	if err != nil {
		// 清理现场
		// 是否断点下载
		return err
	}
	body, err := ioutil.ReadAll(respDl.Body)
	if err != nil {
		return err
	}
	// 创建文件目录
	err = ioutil.WriteFile(s.dir+"/"+worker.fileName, body, os.FileMode(0666))
	return nil
}
