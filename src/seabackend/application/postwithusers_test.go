package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pteich/gosea/src/seabackend/domain/entity"
	"github.com/pteich/gosea/src/seabackend/domain/service/mocks"
)

func TestPostsWithUsers_RetrievePostsWithUsersFromBackend(t *testing.T) {
	testCases := []struct {
		name        string
		remotePosts []entity.RemotePost
		remoteUser  entity.RemoteUser
		filter      string
		wantError   bool
		wantPosts   []entity.Post
	}{
		{
			name:        "test for empty remote posts list",
			remotePosts: []entity.RemotePost{},
			remoteUser:  entity.RemoteUser{},
			filter:      "",
			wantError:   false,
			wantPosts:   []entity.Post{},
		},
		{
			name:        "test for empty remote posts list with filter",
			remotePosts: []entity.RemotePost{},
			remoteUser:  entity.RemoteUser{},
			filter:      "test",
			wantError:   false,
			wantPosts:   []entity.Post{},
		},
		{
			name: "test for 1 remote posts",
			remotePosts: []entity.RemotePost{
				{
					UserID: "1",
					ID:     "1",
					Title:  "post 1",
					Body:   "post body 1",
				},
			},
			remoteUser: entity.RemoteUser{
				ID:       1,
				Name:     "User 1",
				Username: "user1",
				Company: entity.RemoteCompany{
					Name: "Company 1",
				},
			},
			filter:    "",
			wantError: false,
			wantPosts: []entity.Post{
				entity.Post{
					Username:    "user1",
					CompanyName: "Company 1",
					Name:        "",
					Title:       "post 1",
					Body:        "post body 1",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pwu := PostsWithUsers{}
			seaBackendMock := &mocks.SeaBackendLoader{}

			seaBackendMock.On("LoadPosts", mock.Anything).Return(tc.remotePosts, nil)

			if len(tc.remotePosts) == 0 {
				seaBackendMock.On("LoadUser", mock.Anything, mock.Anything).Return(tc.remoteUser, nil)
			}
			for _, posts := range tc.remotePosts {
				seaBackendMock.On("LoadUser", mock.Anything, posts.UserID.String()).Return(tc.remoteUser, nil)
			}

			// emulate Dingo dependency injection
			pwu.Inject(seaBackendMock, &postsWithUsersCfg{WorkerCount: 1})

			postsList, err := pwu.RetrievePostsWithUsersFromBackend(context.TODO(), tc.filter)
			if tc.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tc.wantPosts, postsList)
		})
	}
}
