package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// File is an uploaded file
type File struct {
	Name string
	Data []byte
}

type file struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Post uploads files to an API endpoint
func Post(url, token string, files []*File) error {
	var encoded []*file
	for _, f := range files {
		encoded = append(encoded, &file{Name: f.Name, Data: base64.URLEncoding.EncodeToString(f.Data)})
	}

	j, err := json.Marshal(encoded)
	if err != nil {
		return fmt.Errorf("could not marshal files: %w", err)
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("could not complete request: %w", err)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if string(buf) != `{"code":200,"description":"OK"}` {
		return fmt.Errorf("unexpected response: %v", string(buf))
	}

	return nil
}
