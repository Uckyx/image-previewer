package integration_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var defaultImgURL = "raw.githubusercontent.com/Uckyx/image-previewer/master/img_example/"

func Test_Resize(t *testing.T) {
	ctx := context.Background()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	c := &http.Client{
		Timeout:   60 * time.Second,
		Transport: customTransport,
	}

	t.Parallel()

	tests := []struct {
		name   string
		URL    string
		Status int
	}{
		{
			name:   "status_ok_case",
			URL:    "/resize/200/200/" + defaultImgURL + "_gopher_original_1024x504.jpg",
			Status: http.StatusOK,
		},
		{
			name:   "status_ok_case_large_resize",
			URL:    "/resize/2000/2000/" + defaultImgURL + "_gopher_original_1024x504.jpg",
			Status: http.StatusOK,
		},
		{
			name:   "status_bad_gateway_with_not_walid_url_case",
			URL:    "/resize/200/200/raw.554123.jpg",
			Status: http.StatusBadGateway,
		},
		{
			name:   "status_bad_gateway_with_not_found_img_case",
			URL:    "/resize/200/200/" + defaultImgURL + "foo.jpg",
			Status: http.StatusBadGateway,
		},
		{
			name:   "status_bad_gateway_with_not_support_file_case",
			URL:    "/resize/200/200/raw.githubusercontent.com/Uckyx/image-previewer/dev/.env.dist",
			Status: http.StatusBadGateway,
		},
		{
			name:   "status_not_found_case",
			URL:    "/resize/width/height/" + defaultImgURL + "_gopher_original_1024x504.jpg",
			Status: http.StatusNotFound,
		},
		{
			name:   "status_bad_request_with_not_correct_url_case",
			URL:    "/resize/200/200/awd2q3@DA:::L:L!@#!@/",
			Status: http.StatusBadRequest,
		},
	}

	for k, tt := range tests {
		q := tt
		t.Run(fmt.Sprintf("%s %d", q.name, k), func(t *testing.T) {
			t.Parallel()
			request, _ := http.NewRequestWithContext(ctx, http.MethodGet, buildURL(q.URL), nil)
			resp, err := c.Do(request)
			require.NoError(t, err)
			require.Equal(t, q.Status, resp.StatusCode)
			_, err = readResponse(resp)
			require.NoError(t, err)
		})
	}
}

func buildURL(uri string) string {
	return fmt.Sprintf("%s/%s", getBaseURL(), strings.TrimLeft(uri, "/"))
}

func getBaseURL() string {
	return strings.TrimRight("http://127.0.0.1", "/")
}

func readResponse(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
