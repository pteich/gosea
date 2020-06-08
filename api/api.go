package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pteich/gosea/posts"
)

type postsService interface {
	LoadPosts() ([]posts.RemotePost, error)
}

type Api struct {
	posts  postsService
	logger *log.Logger
}

func New(posts postsService, logger *log.Logger) *Api {
	return &Api{
		posts:  posts,
		logger: logger,
	}
}

// Posts returns a json response with remote posts
func (a *Api) Posts(w http.ResponseWriter, r *http.Request) {
	var err error

	a.logger.Printf("got request %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	remotePosts, err := a.posts.LoadPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := r.URL.Query().Get("filter")

	responsePosts := make([]Post, 0)
	for _, remotePost := range remotePosts {
		if filter != "" && !strings.Contains(strings.ToLower(remotePost.Title), strings.ToLower(filter)) {
			continue
		}

		post := Post{
			Title: remotePost.Title,
			Body:  remotePost.Body,
		}
		responsePosts = append(responsePosts, post)
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(responsePosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
