package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetMethod(t *testing.T) {
	type want struct {
		code        int
		contentType string
		realUrl     string
	}
	tests := []struct {
		name    string
		request string
		want    want
		service URLService
	}{
		{
			name:    "valid ID",
			request: "http://localhost:8080/dfaffdfg",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
				realUrl:     "https://practicum.yandex.ru/",
			},
			service: URLService{
				urls: map[string]string{
					"dfaffdfg": "https://practicum.yandex.ru/",
				},
			},
		},

		{
			name:    "unknown ID",
			request: "http://localhost:8080/ddddd",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				realUrl:     "",
			},
			service: URLService{
				urls: map[string]string{
					"dfaffdfg": "https://practicum.yandex.ru/",
				},
			},
		},

		{
			name:    "empty id '/'",
			request: "http://localhost:8080/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				realUrl:     "",
			},
			service: URLService{
				urls: map[string]string{
					"dfaffdfg": "https://practicum.yandex.ru/",
				},
			},
		},

		{
			name:    "several segments",
			request: "http://localhost:8080/abc/def",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				realUrl:     "",
			},
			service: URLService{
				urls: map[string]string{
					"dfaffdfg": "https://practicum.yandex.ru/",
				},
			},
		},

		{
			name:    "several items in map",
			request: "http://localhost:8080/dfaffdfg",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
				realUrl:     "https://practicum.yandex.ru/",
			},
			service: URLService{
				urls: map[string]string{
					"kjhkjkj":  "https://logo.com/",
					"dfaffdfg": "https://practicum.yandex.ru/",
					"dfjfja":   "https://wiki.com/",
				},
			},
		},

		{
			name:    "query parameters",
			request: "http://localhost:8080/dfaffdfg?test=123",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
				realUrl:     "https://practicum.yandex.ru/",
			},
			service: URLService{
				urls: map[string]string{
					"dfaffdfg": "https://practicum.yandex.ru/",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /{id}", test.service.getMethod)

			mux.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.realUrl, res.Header.Get("Location"))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

}



func TestPostMethod(t *testing.T) {
	type want struct {
		code int
		urls map[string]string
		body string
		contentType string
	}
	tests := []struct {
		name        string
		requestBody string
		want        want
		service     URLService
	}{
		{
			name:        "add new item in map",
			requestBody: "https://practicum.yandex.ru/",
			want: want{
				code: http.StatusCreated,
				contentType: "text/plain",
				urls: map[string]string{
					"abc123": "https://practicum.yandex.ru/",
				},
				body: "http://localhost:8080/abc123",
			},
			service: URLService{
				urls:        map[string]string{},
				idGenerator: func() (string, error) { return "abc123", nil },
			},
		},

		{
			name:        "several items in map",
			requestBody: "https://practicum.yandex.ru/",
			want: want{
				code: http.StatusCreated,
				contentType: "text/plain",
				urls: map[string]string{
					"sff":    "https://logo.com/",
					"abc123": "https://practicum.yandex.ru/",
					"nvvnv":  "https://wiki.com/",
				},
				body: "http://localhost:8080/abc123",
			},
			service: URLService{
				urls: map[string]string{
					"sff":   "https://logo.com/",
					"nvvnv": "https://wiki.com/"},
				idGenerator: func() (string, error) { return "abc123", nil },
			},
		},

		{
			name:        "error wrong generateId() function",
			requestBody: "https://practicum.yandex.ru/",
			want: want{
				code: http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
				urls: map[string]string{},
				body: "generation failed\n",
				
			},
			service: URLService{
				urls:        map[string]string{},
				idGenerator: func() (string, error) { return "abc123", errors.New("generation failed") },
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(test.requestBody))
			// создаём новый Recorder
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("POST /{$}", test.service.postMethod)

			mux.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.urls, test.service.urls)
			assert.Equal(t, test.want.body, w.Body.String())
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

}
