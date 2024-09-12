package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timb418/url-shortener/internal/app/config"
)

func TestShorten(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		reqBody string
		want    want
	}{
		{
			name:    "positive test #1",
			reqBody: "http://mail.ru",
			want: want{
				code:        201,
				response:    "http://localhost:8080/sUmWMS4Q", // Ожидаемый сокращенный URL
				contentType: "text/plain",
			},
		},
		{
			name:    "negative test #1: empty request body",
			reqBody: "",
			want: want{
				code:        400,
				response:    "Wrong URL\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #2: no protocol in URL",
			reqBody: "ya.ru",
			want: want{
				code:        400,
				response:    "Wrong URL\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "negative test #3: duplicate URL",
			reqBody: "http://mail.ru", // Повторный запрос с тем же URL
			want: want{
				code:        400,
				response:    "Failed to read request body\n", // Ошибка из-за дубликата
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.reqBody))
			w := httptest.NewRecorder()
			ShortenGivenLink(w, request, config.NewConfig().BaseURL)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestGetFullLinkByShort(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		shortURL   string
		want       want
		predefined bool // Флаг, указывающий, нужно ли добавлять ссылку в хранилище
	}{
		{
			name:       "positive test #1",
			shortURL:   "OY8MpP6e",
			predefined: true, // Добавляем ссылку в хранилище перед тестом
			want: want{
				code:        307,
				response:    "", // Redirect не возвращает тело ответа
				contentType: "text/plain",
			},
		},
		{
			name:     "negative test: short url not found",
			shortURL: "nonexistent",
			want: want{
				code:        400,
				response:    "URL not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.predefined {
				originalURL := "http://ya.ru"
				err := linkStorage.StoreLink(originalURL, test.shortURL)
				require.NoError(t, err)
			}

			request := httptest.NewRequest(http.MethodGet, "/"+test.shortURL, nil)
			w := httptest.NewRecorder()
			GetFullLinkByShort(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			// Проверяем заголовок Location только для успешного редиректа
			if test.want.code == http.StatusTemporaryRedirect {
				originalURL := "http://ya.ru"
				assert.Equal(t, originalURL, res.Header.Get("Location"))
			}
		})
	}
}
