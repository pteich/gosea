package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pteich/gosea/posts"
)

type postsService interface {
	LoadPosts() ([]posts.RemotePost, error)
}

type Api struct {
	posts postsService
}

func New(posts postsService) *Api {
	return &Api{
		posts: posts,
	}
}

// Posts returns a json response with remote posts
func (a *Api) Posts(w http.ResponseWriter, r *http.Request) {
	var err error

	remotePosts, err := a.posts.LoadPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := r.URL.Query().Get("filter")

	responsePosts := make([]Post, 0)
	for _, remotePost := range remotePosts {
		if !strings.Contains(strings.ToLower(remotePost.Title), strings.ToLower(filter)) {
			continue
		}

		post := Post{
			Title: remotePost.Title,
			Body:  remotePost.Body,
		}
		responsePosts = append(responsePosts, post)
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(responsePosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
