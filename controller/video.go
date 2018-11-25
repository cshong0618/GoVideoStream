package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func VideoStream(writer http.ResponseWriter, request *http.Request) {
	videoPath := chi.URLParam(request, "*")
	videoPath, _ = url.QueryUnescape(videoPath)
	_, err := ioutil.ReadFile(videoPath)

	if err != nil {
		errorResponse := ErrorResponse{err.Error()}

		response, _ := json.Marshal(errorResponse)
		writer.Write(response)
		return
	}

	f, err := os.Open(videoPath)
	file, err := f.Stat()

	fileSize := file.Size()
	fileRange := request.Header.Get("range")

	if len(fileRange) > 0 {
		parts := strings.Split(strings.Replace(fileRange, "bytes=", "", -1), "-")
		start, _ := strconv.ParseInt(parts[0], 10, 8)
		end, _ := strconv.ParseInt(parts[1], 10, 8)

		if end == 0 {
			end = start + 1024000000

			if end > fileSize {
				end = fileSize - 1
			}
		}

		chunkSize := (end - start) + 1

		fileResponse := make([]byte, chunkSize)
		_, _ = f.Seek(start, 0)
		_, _ = f.Read(fileResponse)

		contentRange := fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize)

		writer.Header().Add("Content-Range", contentRange)
		writer.Header().Add("Accept-Ranges", "bytes")
		writer.Header().Add("Content-Length", strconv.FormatInt(chunkSize, 10))
		writer.Header().Add("Content-Type", "video/mkv")

		writer.WriteHeader(206)
		writer.Write(fileResponse)
	} else {
		writer.Header().Add("Content-Length", string(fileSize))
		writer.Header().Add("Content-Type", "video/mkv")

		fileResponse := make([]byte, fileSize)
		f.Read(fileResponse)

		writer.WriteHeader(200)
		writer.Write(fileResponse)
	}
}
