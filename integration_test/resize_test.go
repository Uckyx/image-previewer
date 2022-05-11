package integration_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	c       *http.Client
	timeout = 60 * time.Second
)

func TestGet(t *testing.T) {
	tests := []struct {
		Url    string
		Status int
	}{
		{
			Url:    "/resize/200/200/raw.githubusercontent.com/Uckyx/image-previewer/dev/integration_test/img_example/_gopher_original_1024x504.jpg",
			Status: http.StatusOK,
		},
		{
			Url:    "/resize/200/200/raw.123.jpg",
			Status: http.StatusBadGateway,
		},
		{
			Url:    "/resize/200/200/raw.githubusercontent.com/Uckyx/image-previewer/dev/integration_test/img_example/foo.jpg",
			Status: http.StatusBadGateway,
		},
		{
			Url:    "/resize/2000/2000/raw.githubusercontent.com/Uckyx/image-previewer/dev/integration_test/img_example/_gopher_original_1024x504.jpg",
			Status: http.StatusOK,
		},
		{
			Url:    "/resize/width/height/raw.githubusercontent.com/Uckyx/image-previewer/dev/integration_test/img_example/_gopher_original_1024x504.jpg",
			Status: http.StatusNotFound,
		},
		{
			Url:    "/resize/200/200/raw.githubusercontent.com/Uckyx/image-previewer/dev/integration_test/img_example/.env.dist",
			Status: http.StatusBadGateway,
		},
		{
			Url:    "/resize/200/200/awd2q3@DA:::L:L!@#!@/",
			Status: http.StatusBadRequest,
		},
	}

	for k, tt := range tests {
		q := tt
		t.Run(fmt.Sprintf("%s %d", q.Url, k), func(t *testing.T) {
			t.Parallel()
			resp, err := c.Get(buildUrl(q.Url))
			require.NoError(t, err)
			require.Equal(t, q.Status, resp.StatusCode)
			_, err = readResponse(resp)
			require.NoError(t, err)
		})
	}
}

func init() {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c = &http.Client{
		Timeout:   timeout,
		Transport: customTransport,
	}
}

func getBaseUrl() string {
	return strings.TrimRight("http://127.0.0.1", "/")
}

func buildUrl(uri string) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(), strings.TrimLeft(uri, "/"))
}

func readResponse(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
