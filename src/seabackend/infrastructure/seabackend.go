package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pteich/gosea/src/seabackend/domain/entity"
)

type Cacher interface {
	Get(key string, data interface{}) error
	Set(key string, data interface{}) error
}

// SeaBackend bundles all function to access external json endpoint
type SeaBackend struct {
	endpoint       string
	cache          Cacher
	defaultTimeout time.Duration
	httpClient     *http.Client
}

func (sb *SeaBackend) Inject(cache Cacher, cfg *struct {
	Endpoint       string  `inject:"config:seabackend.endpoint"`
	DefaultTimeout float64 `inject:"config:seabackend.defaultTimeout"`
}) {
	if cfg != nil {
		sb.endpoint = cfg.Endpoint
		sb.defaultTimeout = time.Duration(cfg.DefaultTimeout) * time.Second
		sb.httpClient = &http.Client{
			Timeout: time.Duration(cfg.DefaultTimeout) * time.Second,
		}
	}
	sb.cache = cache
}

// LoadPosts loads all existing posts from external endpoint
func (sb *SeaBackend) LoadPosts(ctx context.Context) ([]entity.RemotePost, error) {
	var remotePosts []entity.RemotePost

	err := sb.load(ctx, sb.endpoint+"/posts", &remotePosts)
	if err != nil {
		return remotePosts, fmt.Errorf("could not load posts: %w", err)
	}

	return remotePosts, nil
}

// LoadUsers loads all existing users from external endpoint
func (sb *SeaBackend) LoadUsers(ctx context.Context) ([]entity.RemoteUser, error) {
	var remoteUsers []entity.RemoteUser
	err := sb.load(ctx, sb.endpoint+"/users", &remoteUsers)
	if err != nil {
		return remoteUsers, fmt.Errorf("could not load users: %w", err)
	}

	return remoteUsers, nil
}

// LoadUsers loads all existing users from external endpoint
func (sb *SeaBackend) LoadUser(ctx context.Context, id string) (entity.RemoteUser, error) {
	var remoteUsers []entity.RemoteUser
	var user entity.RemoteUser

	err := sb.load(ctx, sb.endpoint+"/users?id="+id, &remoteUsers)
	if err != nil {
		return user, fmt.Errorf("could not load user: %w", err)
	}

	if len(remoteUsers) <= 0 {
		return user, fmt.Errorf("could not load user for id %s", id)
	}

	user = remoteUsers[0]

	return user, nil
}

func (sb *SeaBackend) load(ctx context.Context, requestUrl string, data interface{}) (err error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, sb.defaultTimeout)
	defer cancel()

	err = sb.cache.Get(requestUrl, data)
	if err == nil {
		return nil
	}

	req, err := http.NewRequestWithContext(ctxTimeout, http.MethodGet, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("accept-encoding", "application/json")

	res, err := sb.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed execute request: %w", err)
	}
	defer func() {
		bodyClosErr := res.Body.Close()
		if bodyClosErr != nil {
			err = bodyClosErr
		}
	}()

	if res.StatusCode >= 400 {
		return fmt.Errorf("remote server returned status %d", res.StatusCode)
	}

	respData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed load body: %w", err)
	}

	err = json.Unmarshal(respData, data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	err = sb.cache.Set(requestUrl, data)
	if err != nil {
		return fmt.Errorf("failed to save data to cache: %w", err)
	}

	return nil
}
