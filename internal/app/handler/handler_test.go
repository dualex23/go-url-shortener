package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/dualex23/go-url-shortener/internal/app/storage"
	"github.com/dualex23/go-url-shortener/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getServerAddress() string {
	cfg := config.AppParseFlags()
	return cfg.ServerAddr
}

var baseURL string

func TestMain(m *testing.M) {
	logger.New()
	serverAddr := getServerAddress()
	baseURL = fmt.Sprintf("http://%s", serverAddr)

	os.Exit(m.Run())
}

func TestMainHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDataBaseInterface(ctrl)

	// Настраиваем мок для случая, когда URL уже существует
	mockDB.EXPECT().FindByOriginalURL(gomock.Any(), "https://practicum.yandex.ru/").Return("existing-id", fmt.Sprintf("%s/existing-id", baseURL), nil).AnyTimes()

	// Настраиваем мок для случая, когда URL не существует и нужно его создать
	mockDB.EXPECT().FindByOriginalURL(gomock.Any(), "https://newurl.yandex.ru/").Return("", "", sql.ErrNoRows).AnyTimes()
	mockDB.EXPECT().SaveUrls(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tempFile, err := os.CreateTemp("", "test-*.json")
	require.NoError(t, err, "Ошибка при создании временного файла")
	defer os.Remove(tempFile.Name())

	storage := storage.NewStorage(tempFile.Name(), "", mockDB)
	handler := NewShortenerHandler(baseURL, storage)

	tests := []struct {
		name   string
		method string
		body   io.Reader
		want   struct {
			status          int
			responsePattern *regexp.Regexp
		}
	}{
		{
			name:   "test existing URL",
			method: http.MethodPost,
			body:   strings.NewReader("https://practicum.yandex.ru/"),
			want: struct {
				status          int
				responsePattern *regexp.Regexp
			}{
				status:          http.StatusConflict,
				responsePattern: regexp.MustCompile(fmt.Sprintf(`^%s/existing-id$`, baseURL)),
			},
		},
		{
			name:   "test new URL",
			method: http.MethodPost,
			body:   strings.NewReader("https://newurl.yandex.ru/"),
			want: struct {
				status          int
				responsePattern *regexp.Regexp
			}{
				status:          http.StatusCreated,
				responsePattern: regexp.MustCompile(fmt.Sprintf(`^%s/[a-zA-Z0-9]{8}$`, baseURL)),
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

			assert.Equal(t, test.want.status, res.StatusCode)

			resBody, _ := io.ReadAll(res.Body)
			require.NoError(t, err, "Error reading response body")
			assert.Regexp(t, test.want.responsePattern, string(resBody))
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

	storage := &storage.Storage{
		UrlsMap: map[string]storage.URLData{
			"validID": {ID: "validID", OriginalURL: "https://practicum.yandex.ru/", ShortURL: baseURL + "/validID"},
		},
		StorageMode: "memory",
	}
	handler := NewShortenerHandler(baseURL, storage)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDataBaseInterface(ctrl)

	// Настроить возвращение существующего URL
	mockDB.EXPECT().
		FindByOriginalURL(gomock.Any(), "https://practicum.yandex.ru/").
		Return("existing-id", fmt.Sprintf("%s/existing-id", baseURL), nil).
		AnyTimes()

	// Настроить возвращение ошибки для несуществующего URL и сохранить его
	mockDB.EXPECT().FindByOriginalURL(gomock.Any(), "https://new-url.yandex.ru/").Return("", "", sql.ErrNoRows).AnyTimes()
	mockDB.EXPECT().SaveUrls(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tempFile, err := os.CreateTemp("", "test-*.json")
	require.NoError(t, err, "Couldn't create the file")
	defer os.Remove(tempFile.Name())

	storage := storage.NewStorage(tempFile.Name(), "memory", mockDB)
	handler := NewShortenerHandler(baseURL, storage)

	tests := []struct {
		name   string
		method string
		body   io.Reader
		want   struct {
			status          int
			responsePattern *regexp.Regexp
		}
	}{
		{
			name:   "positive test - new URL",
			method: http.MethodPost,
			body:   bytes.NewReader([]byte(`{"url":"https://new-url.yandex.ru/"}`)),
			want: struct {
				status          int
				responsePattern *regexp.Regexp
			}{
				status:          http.StatusCreated,
				responsePattern: regexp.MustCompile(`{"result":"` + regexp.QuoteMeta(baseURL) + `/[a-zA-Z0-9]{8}"}`),
			},
		},
		{
			name:   "negative test - existing URL",
			method: http.MethodPost,
			body:   bytes.NewReader([]byte(`{"url":"https://practicum.yandex.ru/"}`)),
			want: struct {
				status          int
				responsePattern *regexp.Regexp
			}{
				status:          http.StatusConflict,
				responsePattern: regexp.MustCompile(`{"result":"` + regexp.QuoteMeta(fmt.Sprintf("%s/existing-id", baseURL)) + `"}`),
			},
		},
		{
			name:   "negative test - empty URL",
			method: http.MethodPost,
			body:   bytes.NewReader([]byte(`{"url":""}`)),
			want: struct {
				status          int
				responsePattern *regexp.Regexp
			}{
				status:          http.StatusBadRequest,
				responsePattern: regexp.MustCompile(`URL field is required`),
			},
		},
		{
			name:   "negative test - wrong method",
			method: http.MethodGet,
			body:   nil,
			want: struct {
				status          int
				responsePattern *regexp.Regexp
			}{
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

			assert.Equal(t, test.want.status, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Error reading body")

			if test.want.responsePattern != nil {
				assert.Regexp(t, test.want.responsePattern, string(resBody), "Unexpected body")
			}
		})
	}
}
