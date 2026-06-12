package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	if err := os.MkdirAll("./static", 0755); err != nil {
		fmt.Println("Failed to create static directory:", err)
		panic(err)
	}
	// 单个文件上传
	r.POST("/upload", uploadFile)

	// 多个文件上传
	r.POST("/uploads", uploadFiles)

	// 文件下载
	r.GET("/download/:filename", downloadfile)

	// 静态文件
	r.Static("/static", "./static")

	// 文件系统("url路径", "服务器目录")
	r.StaticFS("/files", http.Dir("./static"))
	r.Run(":8080")
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Uploaded File: %#v\n", file.Filename)
	dst := "./static/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "File uploaded successfully", "size": file.Size})
}

func uploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	if len(files) < 1 {
		c.JSON(400, gin.H{"error": "no files uploaded"})
		return
	}
	var filenames []string
	for _, file := range files {
		dst := "./static/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		filenames = append(filenames, file.Filename)
	}
	c.JSON(200, gin.H{"filenames": filenames, "count": len(filenames)})
}

func downloadfile(c *gin.Context) {
	/****/
	filename := c.Param("filename")
	dst := "static/" + filename
	fmt.Println("dst=", dst)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "file not found"})
		return
	}
	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	// c.Header("Content-Type", "application/octet-stream")
	if mtype, err := mimetype.DetectFile(dst); err == nil {
		c.Header("Content-Type", mtype.String())
		fmt.Println("mtype：", mtype.String())
	}
	c.File(dst)
	// */
}
