package infrastructure

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeaBackend_LoadPosts(t *testing.T) {
	sb := SeaBackend{}
	sb.Inject(&config{
		Endpoint: "sa-bonn.ddnss.de:3000",
	})

	posts, err := sb.LoadPosts(context.TODO())
	assert.NoError(t, err)

	for _, post := range posts {
		userid, ok := post.UserID.(json.Number)
		if ok {
			t.Log(userid.String())
		} else {
			userid, ok := post.UserID.(string)
			if ok {
				t.Log(userid)
			} else {
				t.Log("not supported")
			}
		}

	}
}
