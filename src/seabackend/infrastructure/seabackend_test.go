package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pteich/gosea/src/seabackend/domain/entity"
	"github.com/pteich/gosea/src/seabackend/infrastructure/mocks"
)

// manual mocking Cacher interface
type CacheMock struct {
}

func (cm *CacheMock) Get(key string, data interface{}) error {
	return errors.New("not found")
}

func (cm *CacheMock) Set(key string, data interface{}) error {
	return nil
}

func TestPosts_LoadPosts(t *testing.T) {
	tests := []struct {
		name          string
		jsonResponse  string
		status        int
		wrongEndpoint bool
		cacheGetErr   error
		wantErr       bool
		wantResponse  []entity.RemotePost
	}{
		{
			name:         "Normaler Response mit mehreren Werten",
			jsonResponse: `[{"userId": 1, "id":1, "title": "Title1", "body": "Body1"},{"userId": 2, "id":2, "title": "Title2", "body": "Body2"}]`,
			status:       http.StatusOK,
			wantErr:      false,
			cacheGetErr:  errors.New("not found"),
			wantResponse: []entity.RemotePost{
				{
					UserID: json.Number("1"),
					ID:     json.Number("1"),
					Title:  "Title1",
					Body:   "Body1",
				},
				{
					UserID: json.Number("2"),
					ID:     json.Number("2"),
					Title:  "Title2",
					Body:   "Body2",
				},
			},
		},
		{
			name:         "Leerer Response",
			jsonResponse: ``,
			status:       http.StatusOK,
			cacheGetErr:  errors.New("not found"),
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name:         "Response mit Zahlen als String",
			jsonResponse: `[{"userId": "1", "id":"1", "title": "Title1", "body": "Body1"},{"userId": 2, "id":2, "title": "Title2", "body": "Body2"}]`,
			wantErr:      false,
			status:       http.StatusOK,
			cacheGetErr:  errors.New("not found"),
			wantResponse: []entity.RemotePost{
				{
					UserID: json.Number("1"),
					ID:     json.Number("1"),
					Title:  "Title1",
					Body:   "Body1",
				},
				{
					UserID: json.Number("2"),
					ID:     json.Number("2"),
					Title:  "Title2",
					Body:   "Body2",
				},
			},
		},
		{
			name:         "Falscher Status",
			jsonResponse: `[{"userId": "1", "id":"1", "title": "Title1", "body": "Body1"},{"userId": 2, "id":2, "title": "Title2", "body": "Body2"}]`,
			status:       http.StatusInternalServerError,
			cacheGetErr:  errors.New("not found"),
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name:          "Falscher Endpunkt",
			wrongEndpoint: true,
			wantErr:       true,
			wantResponse:  nil,
			cacheGetErr:   errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, err := fmt.Fprint(w, tt.jsonResponse)
				assert.NoError(t, err)
			}))
			defer testSrv.Close()

			// manual mocking Cacher interface
			// cacheMock := CacheMock{}
			cacheMock := mocks.Cacher{}
			cacheMock.On("Get", mock.Anything, mock.Anything).Return(tt.cacheGetErr).Once()
			cacheMock.On("Set", mock.Anything, mock.Anything).Return(nil).Once()

			testPosts := &SeaBackend{}
			testPosts.Inject(&cacheMock, &struct {
				Endpoint       string  `inject:"config:seabackend.endpoint"`
				DefaultTimeout float64 `inject:"config:seabackend.defaultTimeout"`
			}{
				Endpoint:       testSrv.URL,
				DefaultTimeout: 100,
			})

			if tt.wrongEndpoint {
				testPosts.endpoint = ""
			}

			rp, err := testPosts.LoadPosts(context.TODO())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantResponse, rp)
		})
	}
}
