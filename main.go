package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	securejoin "github.com/cyphar/filepath-securejoin"
)

type File struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (f *File) WriteTo(w io.Writer) (int64, error) {
	buf, err := base64.URLEncoding.DecodeString(f.Data)
	if err != nil {
		return 0, fmt.Errorf("could not decode as base64: %w", err)
	}

	n, err := w.Write(buf)
	return int64(n), err
}

func writeJSON(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"code":%d,"description":"%s"}`, code, http.StatusText(code))))
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("ERROR: init: TOKEN must be set")
	}

	root := "."
	if r := os.Getenv("ROOT"); r != "" {
		root = r
	}

	root, err := filepath.Abs(root)
	if err != nil {
		log.Fatalln("ERROR: init: could not get current working directory:", err)
	}

	listen := ":80"
	if l := os.Getenv("LISTENADDR"); l != "" {
		listen = l
	}

	http.HandleFunc("/api/1.0/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		//authorize
		if r.Method != http.MethodPost {
			log.Println("WARNING: http: bad request: bad method")
			writeJSON(w, http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") != "Bearer "+token {
			log.Println("WARNING: http: bad authorization")
			writeJSON(w, http.StatusUnauthorized)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			log.Println("WARNING: http: bad request: bad content-type")
			writeJSON(w, http.StatusBadRequest)
			return
		}

		//parse body
		var files []*File
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&files); err != nil {
			log.Println("ERROR: http: could not decode json:", err)
			writeJSON(w, http.StatusBadRequest)
			return
		}

		for _, file := range files {
			fp, err := securejoin.SecureJoin(root, file.Name)
			if err != nil {
				err = fmt.Errorf("could not find file path: %w", err)
				log.Printf("ERROR: writer: %s: %v\n", file.Name, err)
				writeJSON(w, http.StatusBadRequest)
				return
			}

			if info, err := os.Stat(fp); err == nil {
				if info.IsDir() {
					err = errors.New("directory exists with same name")
					log.Printf("ERROR: writer: %s: %v\n", file.Name, err)
					writeJSON(w, http.StatusBadRequest)
					return
				}
			}

			f, err := os.Create(fp)
			if err != nil {
				err = fmt.Errorf("could not open file: %w", err)
				log.Printf("ERROR: writer: %s: %v\n", fp, err)
				writeJSON(w, http.StatusBadRequest)
				return
			}

			n, err := file.WriteTo(f)
			if err != nil {
				err = fmt.Errorf("could not write file: %w", err)
				log.Printf("ERROR: writer: %s: %v\n", fp, err)
				writeJSON(w, http.StatusInternalServerError)
				return
			}
			log.Printf("INFO: writer: %s: transfer complete (%d bytes)\n", fp, n)
		}

		writeJSON(w, http.StatusOK)
	})

	log.Printf("INFO: server: listening on %s in %s\n", listen, root)
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Printf("ERROR: server: unexpected shutdown: %v\n", err)
	}
}
