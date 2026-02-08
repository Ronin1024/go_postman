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
}

func Fetch(url string, postBody map[string]string, options FetchOptions) (response *http.Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	requestBodyBytes, _ := json.Marshal(postBody)
	requestBody := bytes.NewBuffer(requestBodyBytes)
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(options.Method), url, requestBody)

	if err != nil {
		err = fmt.Errorf("could not create request: %s\n", err)
		return
	}
	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	response, err = http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = fmt.Errorf("Timeout abort: %s\n", err)
			return
		}
		err = fmt.Errorf("could not do request: %s\n", err)
	}
	return
}
