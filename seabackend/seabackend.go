package seabackend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	seaEndpoint    = "http://sa-bonn.ddnss.de:3000"
	defaultTimeout = 10 * time.Second
)

// SeaBackend bundles all function to access external json endpoint
type SeaBackend struct {
	endpoint   string
	httpClient *http.Client
}

// New returns a new initialized SeaBackend struct for given
// endpoint
func New(endpoint string) *SeaBackend {
	return &SeaBackend{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// NewWithSEA returns a new initialized SeaBackend struct pointing
// to SEA json server endpoint
func NewWithSEA() *SeaBackend {
	return New(seaEndpoint)
}

// LoadPosts loads all existing posts from external endpoint
func (p *SeaBackend) LoadPosts(ctx context.Context) ([]RemotePost, error) {
	var remotePosts []RemotePost

	err := p.load(ctx, p.endpoint+"/posts", &remotePosts)
	if err != nil {
		return remotePosts, fmt.Errorf("could not load posts: %w", err)
	}

	return remotePosts, nil
}

// LoadUsers loads all existing users from external endpoint
func (p *SeaBackend) LoadUsers(ctx context.Context) ([]RemoteUser, error) {
	var remoteUsers []RemoteUser
	err := p.load(ctx, p.endpoint+"/users", &remoteUsers)
	if err != nil {
		return remoteUsers, fmt.Errorf("could not load users: %w", err)
	}

	return remoteUsers, nil
}

// LoadUsers loads all existing users from external endpoint
func (p *SeaBackend) LoadUser(ctx context.Context, id string) (RemoteUser, error) {
	var remoteUsers []RemoteUser
	var user RemoteUser

	err := p.load(ctx, p.endpoint+"/users?id="+id, &remoteUsers)
	if err != nil {
		return user, fmt.Errorf("could not load user: %w", err)
	}

	if len(remoteUsers) <= 0 {
		return user, fmt.Errorf("could not load user for id %s", id)
	}

	user = remoteUsers[0]

	return user, nil
}

func (p *SeaBackend) load(ctx context.Context, requestUrl string, data interface{}) (err error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxTimeout, http.MethodGet, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("accept-encoding", "application/json")

	res, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed execute request: %w", err)
	}
	defer func() {
		err = res.Body.Close()
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

	return nil
}
