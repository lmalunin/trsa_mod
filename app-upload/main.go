package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lmalunin/toolkit"
)

func main() {
	mux := routes()

	log.Println("starting server on :8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/upload", uploadFiles)
	mux.HandleFunc("/upload-one", uploadOneFile)

	return mux
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	t := toolkit.Tools{
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml"},
		MaxFileSize:      1024 * 1024 * 1024,
	}

	files, err := t.UploadedFiles(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := ""
	for _, x := range files {
		out += fmt.Sprintf("Uploaded file: %s (NewFileName: %s)\n", x.OriginalFileName, x.NewFileName)
	}

	_, _ = w.Write([]byte(out))
	w.WriteHeader(http.StatusOK)
}

func uploadOneFile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	t := toolkit.Tools{
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml"},
		MaxFileSize:      1024 * 1024 * 1024,
	}

	file, err := t.UploadOneFile(r, "./upload-one")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := fmt.Sprintf("Uploaded file: %s (NewFileName: %s)\n", file.OriginalFileName, file.NewFileName)

	_, _ = w.Write([]byte(out))
	w.WriteHeader(http.StatusOK)

}
