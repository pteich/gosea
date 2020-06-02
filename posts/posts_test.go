package posts

import (
	"testing"
)

func TestPosts_loadPosts(t *testing.T) {
	p := NewWithSEA()

	posts, err := p.LoadPosts()
	if err != nil {
		t.Errorf("error should be nil")
	}

	t.Log(err)
	t.Log(posts)
}
