package go_postman

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type FetchOptions struct {
	Method  string
	Headers map[string]string
	Timeout int
}

func Fetch(url string, postBody map[string]interface{}, options FetchOptions) (response *http.Response, cancel context.CancelFunc, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)

	requestBodyBytes, _ := json.Marshal(postBody)
	requestBody := bytes.NewBuffer(requestBodyBytes)
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(options.Method), url, requestBody)
	if strings.ToUpper(options.Method) == "GET" {
		req, err = http.NewRequestWithContext(ctx, strings.ToUpper(options.Method), url, nil)
	}
	if err != nil {
		cancel() // отменяем контекст только при ошибке
		err = fmt.Errorf("could not create request: %s\n", err)
		return
	}
	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	response, err = http.DefaultClient.Do(req)
	if err != nil {
		cancel() // отменяем контекст только при ошибке
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = fmt.Errorf("Timeout abort: %s\n", err)
			return
		}
		err = fmt.Errorf("could not do request: %s\n", err)
	}
	return
}
