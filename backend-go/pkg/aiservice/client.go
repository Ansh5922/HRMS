package aiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type AIServiceClient struct {
	BaseURL string
	Secret  string
	HTTP    *http.Client
}

type FaceVerifyResult struct {
	Matched    bool    `json:"matched"`
	EmpID      *string `json:"emp_id"`
	Confidence float64 `json:"confidence"`
	Message    string  `json:"message"`
}

func NewClient(baseURL, secret string) *AIServiceClient {
	return &AIServiceClient{
		BaseURL: baseURL,
		Secret:  secret,
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *AIServiceClient) VerifyFace(imageData []byte, filename string) (*FaceVerifyResult, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(imageData); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.BaseURL+"/ai/v1/face/verify", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Internal-Secret", c.Secret)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling face service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("face service error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res FaceVerifyResult
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}