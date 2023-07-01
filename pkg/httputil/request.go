package httputil

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJsonBody(r *http.Request, dst interface{}) error {
	d := json.NewDecoder(r.Body)

	return d.Decode(dst)
}

func SendGetRequest(url string, queryParams map[string]string, headers map[string][]string) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Set query params
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Set request headers
	req.Header = headers

	return sendRequest(req)
}

func sendRequest(req *http.Request) (int, []byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return resp.StatusCode, respBody, nil
}
