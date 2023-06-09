package http

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/BytemanD/easygo/pkg/fileutils"
	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/PuerkitoBio/goquery"
)

func GetHtml(url string) goquery.Document {
	resp, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	return *doc
}

func UrlJoin(url1 string, url2 string) string {
	re := regexp.MustCompile("^http(s)://")
	if re.FindString(url2) != "" {
		return url2
	} else {
		return strings.TrimRight(url1, "/") + "/" + strings.TrimLeft(url2, "/")
	}
}

func GetLinks(url string, regex string) list.List {
	doc := GetHtml(url)
	links := list.New()

	reg := regexp.MustCompile(regex)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if regex == "" || reg.FindString(href) != "" {
			links.PushBack(UrlJoin(url, href))
		}
	})
	return *links
}

func Download(url string, output string) error {
	_, fileName := filepath.Split(url)
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("下载 %s 失败, 原因: %s", url, err)
		return err
	}

	defer resp.Body.Close()

	fp := fileutils.FilePath{Path: output}
	if err := fp.MakeDirs(); err != nil {
		return err
	}

	outputPath := path.Join(output, fileName)
	outputFile, _ := os.Create(outputPath)
	defer outputFile.Close()

	wt := bufio.NewWriter(outputFile)
	io.Copy(wt, resp.Body)
	wt.Flush()
	return nil
}

func DownloadLinksInHtml(url string, regex string, output string) {
	links := GetLinks(url, regex)

	var (
		wg sync.WaitGroup
	)
	link := links.Front()
	wg.Add(links.Len())
	logging.Info("链接数量: %d", links.Len())
	if links.Len() > 0 {
		logging.Info("开始下载")
	}
	for i := 0; i < links.Len(); i++ {
		go func(aaa list.Element) {
			Download(fmt.Sprintf("%s", aaa.Value), output)
			logging.Debug("下载完成: %s ", aaa.Value)
			wg.Done()
		}(*link)
		link = link.Next()
		if link == nil {
			break
		}
	}
	wg.Wait()
	if links.Len() > 0 {
		logging.Info("下载完成")
	}
}
