package seabackend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosts_LoadPosts(t *testing.T) {
	tests := []struct {
		name          string
		jsonResponse  string
		status        int
		wrongEndpoint bool
		wantErr       bool
		wantResponse  []RemotePost
	}{
		{
			name:         "Normaler Response mit mehreren Werten",
			jsonResponse: `[{"userId": 1, "id":1, "title": "Title1", "body": "Body1"},{"userId": 2, "id":2, "title": "Title2", "body": "Body2"}]`,
			status:       http.StatusOK,
			wantErr:      false,
			wantResponse: []RemotePost{
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
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name:         "Response mit Zahlen als String",
			jsonResponse: `[{"userId": "1", "id":"1", "title": "Title1", "body": "Body1"},{"userId": 2, "id":2, "title": "Title2", "body": "Body2"}]`,
			wantErr:      false,
			status:       http.StatusOK,
			wantResponse: []RemotePost{
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
			wantErr:      true,
			wantResponse: nil,
		},
		{
			name:          "Falscher Endpunkt",
			wrongEndpoint: true,
			wantErr:       true,
			wantResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.jsonResponse)
			}))
			defer testSrv.Close()

			testPosts := &SeaBackend{
				endpoint:   testSrv.URL,
				httpClient: testSrv.Client(),
			}

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

func TestSeaBackend_LoadUsers(t *testing.T) {

	sb := NewWithSEA()

	users, err := sb.LoadUsers(context.TODO())
	assert.NoError(t, err)

	t.Log(users)

}
