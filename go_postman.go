package go_postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FetchOptions struct {
	Method  string
	Headers map[string]string
}

func Fetch(url string, postBody map[string]string, options FetchOptions) (response *http.Response, err error) {

	requestBodyBytes, _ := json.Marshal(postBody)
	requestBody := bytes.NewBuffer(requestBodyBytes)
	req, err := http.NewRequest(strings.ToUpper(options.Method), url, requestBody)

	if err != nil {
		err = fmt.Errorf("could not create request: %s\n", err)
		return
	}
	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	response, err = http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("could not do request: %s\n", err)
	}
	return
}
