package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	type want struct {
		status          int
		response        *string
		responsePattern *regexp.Regexp
		contentType     string
	}
	tests := []struct {
        name    string
        method  string
		body    io.Reader
        want    want
    }{
        {
            name:   "positive test",
            method: http.MethodPost,
			body: strings.NewReader("https://practicum.yandex.ru/"),
            want: want{
                status:          http.StatusCreated,
                responsePattern: regexp.MustCompile(`^http://localhost:8080/[a-zA-Z0-9]{8}$`),
                contentType:     "text/plain",
            },
        },
        {
            name:   "unsupported method",
            method: http.MethodGet,
			body: nil,
            want: want{
                status:      http.StatusMethodNotAllowed,
                contentType: "text/plain; charset=utf-8",
            },
        },
		{
			name:   "empty request body",
			method: http.MethodPost,
			body:   strings.NewReader(""),
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "incorrect content-type",
			method: http.MethodPost,
			body:   strings.NewReader("https://practicum.yandex.ru/"),
			want: want{
				status:      http.StatusBadRequest,
				contentType: "application/json",
			},
		},
    }
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", test.body)
			if test.name == "positive test" || test.name == "empty request body" || test.name == "incorrect content-type" {
				request.Header.Set("Content-Type", "text/plain")
			}

			if test.name == "incorrect content-type" {
				request.Header.Set("Content-Type", "application/json")
			}


			w := httptest.NewRecorder()
			MainHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			if assert.NotNil(t,res) {
				assert.Equal(t, test.want.status, res.StatusCode)
			}

			if test.want.contentType != "" {
				assert.Regexp(t, regexp.MustCompile(`^text/plain;?.*`), res.Header.Get("Content-Type"))
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

	// Инициализируем mapUrls перед тестированием
	mapUrls = map[string]string{
		"validID": "https://practicum.yandex.ru/",
	}

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
		{
			name:   "GET request with non-existing ID",
			method: http.MethodGet,
			path:   "/nonExistingID",
			want: want{
				status: http.StatusNotFound,
			},
		},
		{
			name:   "unsupported method",
			method: http.MethodPost,
			path:   "/validID",
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "missing ID",
			method: http.MethodGet,
			path:   "/",
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			GetHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			// Проверяем статус ответа
			assert.Equal(t, tc.want.status, res.StatusCode)

			// Проверяем Location, если он указан в ожидаемом результате
			if tc.want.location != "" {
				assert.Equal(t, tc.want.location, res.Header.Get("Location"))
			}

			// Проверяем Content-Type, если он указан в ожидаемом результате
			if tc.want.contentType != "" {
				assert.Equal(t, tc.want.contentType, res.Header.Get("Content-Type"))
			}

			// Если есть ожидаемый ответ, проверяем его
			if tc.want.response != nil {
				body, _ := io.ReadAll(res.Body)
				assert.Equal(t, *tc.want.response, string(body))
			}

			// Если есть паттерн для проверки ответа, используем его
			if tc.want.responsePattern != nil {
				body, _ := io.ReadAll(res.Body)
				assert.Regexp(t, tc.want.responsePattern, string(body))
			}
		})
	}
}

