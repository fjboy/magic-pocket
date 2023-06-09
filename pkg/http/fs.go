package http

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

var HTML = `
<!DOCTYPE html>
<html lang="zh-CN">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<link rel="stylesheet" href="https://cdn.staticfile.org/font-awesome/4.7.0/css/font-awesome.css">
		<title>SimpleHttpFS</title>
	</head>
	<body>
		<div style="margin-left: 50px">
			<table>
				<tr></th><th>类型</th><th>名称</th></tr>
				{{ range $index, $webFile := . }}
					<tr>
						{{ if $webFile.IsDir }}
						<td><i class="fa fa-folder" style="color: orange;"></i></td>
						<td><a href="{{$webFile.WebPath}}" > {{ $webFile.Name }} </a></li></td>
						{{ else }}
						<td><i class="fa fa-file" style="color: green;"></i></td>
						<td><a target="view_window" href="{{$webFile.WebPath}}" > {{ $webFile.Name }} </a></li></td>
						{{ end }}
					</tr>
				{{ end}}
			<table>
		</div>
</body>
<html>
`
var FSConfig HTTPFSConfig

type HTTPFSConfig struct {
	Port int16
	Root string
}

type WebFile struct {
	Dir  string
	Name string
}

func (webFile *WebFile) LogicPath() string {
	return filepath.Join(webFile.Dir, webFile.Name)
}
func (webFile *WebFile) WebPath() string {
	return strings.ReplaceAll(webFile.LogicPath(), "\\", "/")
}

func (webFile *WebFile) IsDir() bool {
	filePath := filepath.Join(FSConfig.Root, webFile.LogicPath())
	fi, _ := os.Stat(filePath)
	return fi.IsDir()
}

func handleError(err error, respWriter http.ResponseWriter, request *http.Request) {
	logging.Error("目录处理异常, %s", err)
	switch {
	case os.IsNotExist(err):
		fmt.Fprintf(respWriter, "路径不存在")
	case os.IsPermission(err):
		respWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(respWriter, "目录无访问权限")
	default:
		respWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(respWriter, "访问目录失败 %s", request.URL.Path)
	}
}
func handleFileDownload(respWriter http.ResponseWriter, request *http.Request) {
	filePath := filepath.Join(FSConfig.Root, request.URL.Path)
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	file, _ := os.Open(filePath)
	defer file.Close()

	respWriter.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	logging.Info("下载文件: %s", file.Name())
	// TODO: 优化文件名
	http.ServeFile(respWriter, request, filePath)
}
func FilePathHandler(respWriter http.ResponseWriter, request *http.Request) {
	logging.Info("请求地址 %s", request.URL.Path)
	dirPath := filepath.Join(FSConfig.Root, request.URL.Path)
	dir, err := os.Stat(dirPath)
	if err != nil {
		handleError(err, respWriter, request)
		return
	}
	webFiles := []WebFile{}
	if dir.IsDir() {
		rd, err := ioutil.ReadDir(dirPath)
		if err != nil {
			handleError(err, respWriter, request)
			return
		}
		for _, fi := range rd {
			webFile := WebFile{
				Dir:  filepath.Join(request.URL.Path),
				Name: fi.Name(),
			}
			webFiles = append(webFiles, webFile)
		}
		tmpl, _ := template.New("dirPage").Parse(HTML)
		tmpl.Execute(respWriter, webFiles)
	} else {
		handleFileDownload(respWriter, request)
	}
}

func SimpleHttpFS() error {
	http.HandleFunc("/", FilePathHandler) //设置访问的路由
	logging.Info("启动文件服务器 端口:%d, 目录: %s",
		FSConfig.Port, FSConfig.Root)
	if _, err := os.Stat(FSConfig.Root); err != nil {
		return err
	}
	return http.ListenAndServe(
		fmt.Sprintf(":%d", FSConfig.Port),
		nil)
}
