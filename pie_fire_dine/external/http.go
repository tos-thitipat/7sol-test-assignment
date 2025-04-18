package external

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type ExternalApiRequester interface {
	GetSourceText(ctx context.Context, url string) (int, []byte, error)
}

type httpRequest struct {
	httpRequestFunc func(context.Context, string, []byte, string, string) (int, []byte, error)
}

func NewHttpRequester(requestFunc func(context.Context, string, []byte, string, string) (int, []byte, error)) ExternalApiRequester {
	return &httpRequest{httpRequestFunc: RequestHttp}
}

func RequestHttp(ctx context.Context, url string, input []byte, method string, token string) (int, []byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(input))
	if err != nil {
		return 0, nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", token)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, responseBody, nil
}

func (h *httpRequest) GetSourceText(ctx context.Context, url string) (int, []byte, error) {
	return h.httpRequestFunc(ctx, url, nil, "GET", "")
}
