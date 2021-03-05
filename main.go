package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gobuffalo/packr"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//假设文件上传为本地的服务器，上传的基础路径


// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// 获取当前执行文件的路径
		path, err := os.Executable()
		if err != nil {
			fmt.Println(err)
			return
		}
		uploadPath := filepath.Dir(path)
		fmt.Println(uploadPath)

		f, err := os.OpenFile(uploadPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func handleDownload (w http.ResponseWriter, request *http.Request) {
	//文件上传只允许GET方法
	if request.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}

	// 获取当前执行文件的路径
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	BaseUploadPath := filepath.Dir(path)
	fmt.Println(BaseUploadPath)

	filename := "down.zip"
	log.Println("filename: " + filename)

	//打开文件
	file, err := os.Open(BaseUploadPath + "/download" + "/" + filename)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	//结束后关闭文件
	defer file.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}

//func LoadFile(w http.ResponseWriter, r *http.Request) {
//	file, _:= os.Open("D:\\dongqingyang\\Desktop\\work\\" + "1.txt")
//	defer file.Close()
//	//设置响应的header头
//	w.Header().Add("Content-type", "application/octet-stream")
//	w.Header().Add("content-disposition", "attachment; filename=\""+"1.txt"+"\"")
//	buff, _ := ioutil.ReadAll(file)
//	w.Write(buff)
//}


//https://github.com/zupzup/golang-http-file-upload-download
func main() {
	box := packr.NewBox("./template")

	http.HandleFunc("/fileUpload", upload)

	http.HandleFunc("/fileDownload", handleDownload)

	http.Handle("/", http.FileServer(box))

	error := http.ListenAndServe(":3000", nil)
	if error != nil  {
		log.Fatal("Server run fail")
	}
}






