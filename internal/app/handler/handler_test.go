package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	handler := NewShortenerHandler("http://localhost:8080")

	type want struct {
		status          int
		responsePattern *regexp.Regexp
	}
	tests := []struct {
		name   string
		method string
		body   io.Reader
		want   want
	}{
		{
			name:   "positive test",
			method: http.MethodPost,
			body:   strings.NewReader("https://practicum.yandex.ru/"),
			want: want{
				status:          http.StatusCreated,
				responsePattern: regexp.MustCompile(`^http://localhost:8080/[a-zA-Z0-9]{8}$`),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", test.body)

			w := httptest.NewRecorder()
			handler.MainHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			if assert.NotNil(t, res) {
				assert.Equal(t, test.want.status, res.StatusCode)
			}

			if test.want.responsePattern != nil {
				resBody, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Regexp(t, test.want.responsePattern, string(resBody))
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	type want struct {
		status          int
		response        *string
		responsePattern *regexp.Regexp
		contentType     string
		location        string
	}

	handler := NewShortenerHandler("http://localhost:8080")
	handler.MapURLs["validID"] = "https://practicum.yandex.ru/"

	tests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			name:   "valid GET request with existing ID",
			method: http.MethodGet,
			path:   "/validID",
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "https://practicum.yandex.ru/",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			handler.GetHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.want.status, res.StatusCode)

			if tc.want.location != "" {
				assert.Equal(t, tc.want.location, res.Header.Get("Location"))
			}

			if tc.want.contentType != "" {
				assert.Equal(t, tc.want.contentType, res.Header.Get("Content-Type"))
			}

			if tc.want.response != nil {
				body, _ := io.ReadAll(res.Body)
				assert.Equal(t, *tc.want.response, string(body))
			}

			if tc.want.responsePattern != nil {
				body, _ := io.ReadAll(res.Body)
				assert.Regexp(t, tc.want.responsePattern, string(body))
			}
		})
	}
}

func TestApiHandler(t *testing.T) {
	handler := NewShortenerHandler("http://localhost:8080")

	type want struct {
		status          int
		responsePattern *regexp.Regexp
	}
	tests := []struct {
		name   string
		method string
		body   io.Reader
		want   want
	}{
		{
			name:   "positive test",
			method: http.MethodPost,
			body:   bytes.NewReader([]byte(`{"url":"https://practicum.yandex.ru/"}`)),
			want: want{
				status:          http.StatusCreated,
				responsePattern: regexp.MustCompile(`{"result":"http://localhost:8080/[a-zA-Z0-9]{8}"}`),
			},
		},
		{
			name:   "negative test - empty URL",
			method: http.MethodPost,
			body:   bytes.NewReader([]byte(`{"url":""}`)),
			want: want{
				status:          http.StatusBadRequest,
				responsePattern: regexp.MustCompile(`URL field is required`),
			},
		},
		{
			name:   "negative test - wrong method",
			method: http.MethodGet,
			body:   nil,
			want: want{
				status:          http.StatusMethodNotAllowed,
				responsePattern: regexp.MustCompile(`Only POST request is allowed!`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/api/shorten", test.body)
			request.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.APIHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			if assert.NotNil(t, res) {
				assert.Equal(t, test.want.status, res.StatusCode)
			}

			if test.want.responsePattern != nil {
				resBody, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Regexp(t, test.want.responsePattern, string(resBody))
			}
		})
	}
}
