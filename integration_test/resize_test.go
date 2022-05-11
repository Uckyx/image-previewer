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
		URL    string
		Status int
	}{
		{
			URL:    "/resize/200/200/" + defaultImgURL + "_gopher_original_1024x504.jpg",
			Status: http.StatusOK,
		},
		{
			URL:    "/resize/200/200/raw.554123.jpg",
			Status: http.StatusBadGateway,
		},
		{
			URL:    "/resize/200/200/" + defaultImgURL + "foo.jpg",
			Status: http.StatusBadGateway,
		},
		{
			URL:    "/resize/2000/2000/" + defaultImgURL + "_gopher_original_1024x504.jpg",
			Status: http.StatusOK,
		},
		{
			URL:    "/resize/width/height/" + defaultImgURL + "_gopher_original_1024x504.jpg",
			Status: http.StatusNotFound,
		},
		{
			URL:    "/resize/200/200/raw.githubusercontent.com/Uckyx/image-previewer/dev/.env.dist",
			Status: http.StatusBadGateway,
		},
		{
			URL:    "/resize/200/200/awd2q3@DA:::L:L!@#!@/",
			Status: http.StatusBadRequest,
		},
	}

	for k, tt := range tests {
		q := tt
		t.Run(fmt.Sprintf("%s %d", q.URL, k), func(t *testing.T) {
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