package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	serverAddr  = "http://localhost:8080"
	messagesUrl = serverAddr + "/v1/messages"
	filesUrl    = serverAddr + "/v1/files"
)

func GetRequestWithTextBody(ctx context.Context, text string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		messagesUrl,
		bytes.NewReader([]byte(text)),
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	request.Header.Add("Content-Type", "text/html")
	request.Header.Set("Authorization", "Basic "+creds.base64())

	return request, nil
}

func GetRequestWithFile(ctx context.Context, filePath string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("writer.CreateFormFile: %w", err)
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("writer.Close: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		filesUrl,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", "Basic "+creds.base64())

	return request, nil
}
