package main

import (
	_ "UploadAndDownload/implement/download"
	_ "UploadAndDownload/implement/upload"
	r "UploadAndDownload/utils/gin"
)

func main() {
	r.Router.Run()
}
