// --- internal/handlers/upload.go ---
package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("data", 0755)
	path := filepath.Join("data", handler.Filename)
	out, _ := os.Create(path)
	defer out.Close()
	io.Copy(out, file)

	w.Write([]byte("File uploaded: " + handler.Filename))
}