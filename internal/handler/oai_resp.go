package handler

import (
	"bufio"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func getContentType(header http.Header) string {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	return contentType
}

func HandleOAIResponse1(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := getContentType(respHeader)
	ctx.Writer.Header().Set("Content-Type", contentType)
	_, err := io.Copy(ctx.Writer, responseBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
	}
}

func HandleOAIResponse2(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := getContentType(respHeader)
	ctx.Writer.Header().Set("Content-Type", contentType)
	reader := bufio.NewReader(responseBody)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, writeErr := ctx.Writer.Write(buf[:n]); writeErr != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
				return
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
			return
		}
	}
}

func HandleOAIResponse3(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := getContentType(respHeader)
	ctx.Writer.Header().Set("Content-Type", contentType)
	scanner := bufio.NewScanner(responseBody)
	for scanner.Scan() {
		line := scanner.Bytes()
		if _, err := ctx.Writer.Write(line); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
			return
		}
		ctx.Writer.Write([]byte("\n"))
		ctx.Writer.Flush()
	}
	if err := scanner.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
}

func HandleOAIResponse4(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := getContentType(respHeader)
	ctx.Writer.Header().Set("Content-Type", contentType)
	ctx.Stream(func(w io.Writer) bool {
		scanner := bufio.NewScanner(responseBody)
		for scanner.Scan() {
			line := scanner.Bytes()
			if _, err := w.Write(line); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
				return false
			}
			w.Write([]byte("\n"))
			w.(http.Flusher).Flush()
		}
		if err := scanner.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
			return false
		}
		return false
	})

}

// HandleOAIResponse5 有问题，用不了
func HandleOAIResponse5(ctx *gin.Context, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()
	for k, v := range respHeader {
		ctx.Writer.Header().Set(k, v[0])
	}
	contentType := respHeader.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	if contentType != "text/event-stream" {
		HandleOAIResponse1(ctx, responseBody, respHeader)
		return
	}
	scanner := bufio.NewScanner(responseBody)
	for scanner.Scan() {
		line := scanner.Bytes()
		res := make(map[string]any)
		if err := sonic.Unmarshal(line, &res); err != nil {
			fmt.Println("line", string(line))
		}
		jsonStr, _ := sonic.Marshal(res["data"])
		ctx.SSEvent("data", jsonStr)
		ctx.Writer.Flush()
	}
	if err := scanner.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
}
