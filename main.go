package main

import (
	"errors"
	"fmt"
	"log"

	//"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// "strconv"
	// "strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", handleUpload)

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

type FileDetails struct {
	FileName string  `json:"fileName"`
	FileSize float64 `json:"fileSize"`
}

func handleUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var uploadedFiles []FileDetails

	files := form.File["files"]
	for key, file := range files {
		err := saveFile(file, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate file size in MB
		fileSize := float64(file.Size) / (1024 * 1024)

		// Collect file details
		uploadedFiles = append(uploadedFiles, FileDetails{
			FileName: file.Filename,
			FileSize: fileSize,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully", "files": uploadedFiles})
}

func saveFile(fileHeader *multipart.FileHeader, key int) error {
	fmt.Println(fileHeader.Filename)
	ext := strings.Split(fileHeader.Filename, ".")
	if len(ext) < 2 || ext[1] != "pdf" {
		return errors.New("pass a valid pdf")
	}

	if fileHeader.Size > 2*1024*1024 {
		return errors.New("File size exceeds the 2MB limit")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	// Create a new file in the desired destination folder
	dstPath := filepath.Join("./Storage", ext[0]+strconv.Itoa(key)+"."+ext[1]) // //fileHeader.Filename

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// fmt.Printf("File saved to: %s\n", dstPath)
	return nil
}
