package main

import (
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type ImageInfo struct {
	URL    string
	Width  int
	Height int
	Size   string
}

type PostData struct {
	ID     string
	Images []ImageInfo
}

func main() {
	http.HandleFunc("/post/", postHandler)
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("./data"))))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/post/"):]

	dirPath := fmt.Sprintf("./data/post/%s/", id)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Directory not found", http.StatusNotFound)
		return
	}

	var images []ImageInfo
	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				filePath := filepath.Join(dirPath, file.Name())
				imgFile, err := os.Open(filePath)
				if err != nil {
					continue
				}
				defer imgFile.Close()

				img, _, err := image.DecodeConfig(imgFile)
				if err != nil {
					continue
				}

				fileInfo, err := imgFile.Stat()
				if err != nil {
					continue
				}

				sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
				sizeMBStr := fmt.Sprintf("%.3f", sizeMB)

				images = append(images, ImageInfo{
					URL:    fmt.Sprintf("/data/post/%s/%s", id, file.Name()),
					Width:  img.Width,
					Height: img.Height,
					Size:   sizeMBStr,
				})
			}
		}
	}

	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := PostData{
		ID:     id,
		Images: images,
	}

	tmpl.Execute(w, data)
}