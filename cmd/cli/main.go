package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DiscoreMe/SecureCloud/config"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

var address string
var storage string
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {
	cfg := config.NewConfig()
	address = cfg.Address

	rootCmd := &cobra.Command{
		Use:   "sc",
		Short: "Secure Cloud Manager",
	}

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload file",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				file, err := os.Open(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("file %s open error: %s\n", arg, err))
					continue
				}
				if url, err := uploadFile(cfg, file.Name(), file); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("file %s upload error: %s\n", arg, err))
					continue
				} else {
					fmt.Println(url)
				}
			}
		},
	}
	uploadCmd.Flags().StringVarP(&storage, "storage", "s", "local", "Type storage")

	rootCmd.AddCommand(uploadCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func uploadFile(cfg *config.Config, filename string, r io.ReadCloser) (string, error) {
	defer r.Close()

	var u url.URL
	u.Scheme = "http"
	u.Host = cfg.Address
	u.Path = "/api/upload"
	values := u.Query()
	values.Add("storage", storage)
	u.RawQuery = values.Encode()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	field, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(field, r); err != nil {
		return "", err
	}
	w.Close()

	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("token", cfg.Token)

	httpResp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("bad status: %s", httpResp.Status)
	}

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	type uploadResp struct {
		URL string `json:"url"`
	}

	var resp uploadResp
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}

	return resp.URL, nil
}
