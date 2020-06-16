package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pteich/gosea/seabackend"
)

type postsService interface {
	LoadPosts(ctx context.Context) ([]seabackend.RemotePost, error)
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

// SeaBackend returns a json response with remote posts
func (a *Api) Posts(w http.ResponseWriter, r *http.Request) {
	var err error

	a.logger.Printf("got request %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	ctxValue := context.WithValue(r.Context(), "id", 1)

	remotePosts, err := a.posts.LoadPosts(ctxValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := r.URL.Query().Get("filter")

	responsePosts := make([]Post, 0)
	for _, remotePost := range remotePosts {
		if !remotePost.Contains(filter, seabackend.FieldAll) {
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
