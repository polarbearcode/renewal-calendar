package handlers

import (
    "bytes"
    "fmt"
    "net/http"

    "myapp/internal/services"
	"mime/multipart"
)

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // or "*" for all origins
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == http.MethodOptions {
        // Preflight request, return OK without further handling
        w.WriteHeader(http.StatusNoContent)
        return false
    }
    return true
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	 if !enableCORS(w, r) {
        return
    }


    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(w, "No files found in form", http.StatusBadRequest)
		return
	}

    files := r.MultipartForm.File["file"] // slice of *multipart.FileHeader
    if len(files) == 0 {
        http.Error(w, "No files uploaded", http.StatusBadRequest)
        return
    }

	uploadToS3(w, r, files)
}

func uploadToS3(w http.ResponseWriter, r *http.Request, files []*multipart.FileHeader) {

    for _, header := range files {
		fmt.Println("Uploading:", header.Filename)
        file, err := header.Open()
        if err != nil {
            http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
            return
        }

        buf := new(bytes.Buffer)
        _, err = buf.ReadFrom(file)
        file.Close()
        if err != nil {
            http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Upload to S3 (using your service layer)
        err = services.UploadToS3(r.Context(), header.Filename, buf.Bytes())
        if err != nil {
            http.Error(w, "Failed to upload "+header.Filename+": "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

    fmt.Fprintf(w, "✅ Uploaded %d files", len(files))

}
