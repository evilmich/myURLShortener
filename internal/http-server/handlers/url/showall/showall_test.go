package showall_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"myURLShortener/internal/http-server/handlers/url/showall"
	"myURLShortener/internal/http-server/handlers/url/showall/mocks"
	"myURLShortener/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowHandler(t *testing.T) {
	cases := []struct {
		name      string
		data      []*showall.AliasUrl
		respError string
		mockError error
	}{
		{
			name: "Success",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlShowerMock := mocks.NewURLShower(t)

			if tc.respError == "" || tc.mockError != nil {
				urlShowerMock.On("ShowURL").
					Return(tc.data, tc.mockError).Once()
			}

			handler := showall.New(slogdiscard.NewDiscardLogger(), urlShowerMock)

			req, err := http.NewRequest(http.MethodGet, "/show", bytes.NewReader([]byte("")))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp showall.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
