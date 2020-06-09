package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pteich/gosea/posts"
)

type postsMock struct {
	remotePosts []posts.RemotePost
	err         error
}

func (pm *postsMock) LoadPosts() ([]posts.RemotePost, error) {
	return pm.remotePosts, pm.err
}

func TestApi_Posts(t *testing.T) {
	logBuf := &bytes.Buffer{}

	testApi := &Api{
		posts: &postsMock{
			remotePosts: []posts.RemotePost{
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
			err: nil,
		},
		logger: log.New(logBuf, "test", log.LstdFlags),
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost/posts", nil)
	w := httptest.NewRecorder()

	testApi.Posts(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("content-type"))

	var responsePosts []Post
	json.NewDecoder(w.Body).Decode(&responsePosts)

	expected := []Post{
		{
			Title: "Title1",
			Body:  "Body1",
		},
		{
			Title: "Title2",
			Body:  "Body2",
		},
	}

	assert.Equal(t, expected, responsePosts)
}
