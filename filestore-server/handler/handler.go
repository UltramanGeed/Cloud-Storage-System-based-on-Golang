package handler

import (
	"net/http"
	"io"
	"io/ioutil"
	"ioutil"
)

//处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}
	else if r.Method == "POST"{
		//接收文件流及存储到本地目录
	}
}