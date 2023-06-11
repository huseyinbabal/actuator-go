package actuator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	InfoEndpoint     = "/info"
	HeapDumpEndpoint = "/heapdump"
	HealthEndpoint   = "/health"
)

type Actuator struct {
	ApiURL string
	Client *http.Client
}

func New(apiUrl string) *Actuator {
	return NewWithHttpClient(apiUrl, &http.Client{})
}

func NewWithHttpClient(apiUrl string, client *http.Client) *Actuator {
	return &Actuator{ApiURL: apiUrl, Client: client}
}

func (a *Actuator) Health() (*HealthResponse, error) {
	url := fmt.Sprintf("%s%s", a.ApiURL, HealthEndpoint)
	res, err := Get[HealthResponse](context.Background(), a.Client, url, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *Actuator) Info() (*InfoResponse, error) {
	url := fmt.Sprintf("%s%s", a.ApiURL, InfoEndpoint)
	res, err := Get[InfoResponse](context.Background(), a.Client, url, false)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *Actuator) HeapDump() (*HeapDumpResponse, error) {
	url := fmt.Sprintf("%s%s", a.ApiURL, HeapDumpEndpoint)
	res, err := Get[HeapDumpResponse](context.Background(), a.Client, url, true)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func Get[T any](ctx context.Context, httpClient *http.Client, url string, binary bool) (T, error) {
	var m T
	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return m, fmt.Errorf("failed to create request %w", err)
	}
	res, err := httpClient.Do(r)
	if err != nil {
		return m, fmt.Errorf("failed to do http request %w", err)
	}
	if binary {
		f, err := downloadFile(ctx, httpClient, url, "hprof.")
		if err != nil {
			return m, fmt.Errorf("failed to download file %w", err)
		}
		hd := HeapDumpResponse{Location: f.Name()}
		marshal, err := json.Marshal(hd)
		if err != nil {
			return m, fmt.Errorf("failed to marshall heapdump response %v", err)
		}
		return parseJSON[T](marshal)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, fmt.Errorf("failed to unmarshall http response %w", err)
	}
	return parseJSON[T](body)
}
func parseJSON[T any](s []byte) (T, error) {
	var r T
	if err := json.Unmarshal(s, &r); err != nil {
		return r, err
	}
	return r, nil
}

func downloadFile(ctx context.Context, httpClient *http.Client, url, filePattern string) (*os.File, error) {
	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request %w", err)
	}
	res, err := httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request %w", err)
	}
	temp, err := os.CreateTemp(os.TempDir(), filePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file %w", err)
	}
	io.Copy(temp, res.Body)
	return temp, nil
}
