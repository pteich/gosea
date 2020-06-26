package seabackend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"flamingo.me/flamingo/v3/framework/web"
)

const (
	seaEndpoint     = "http://sa-bonn.ddnss.de:3000"
	defaultTimeout  = 10 * time.Second
	defaultCacheTTL = 10 * time.Second
	workerCount     = 3
)

type Cache interface {
	Get(key string, data interface{}) error
	Set(key string, data interface{}) error
}

// SeaBackend bundles all function to access external json endpoint
type SeaBackend struct {
	responder  *web.Responder
	endpoint   string
	cache      Cache
	httpClient *http.Client
}

// New returns a new initialized SeaBackend struct for given
// endpoint
func New(endpoint string) *SeaBackend {
	return &SeaBackend{
		endpoint: endpoint,
		cache:    NewRequestCache(defaultCacheTTL),
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

func (sb *SeaBackend) Inject(responder *web.Responder, cache *RequestCache) *SeaBackend {
	sb.responder = responder
	sb.endpoint = seaEndpoint
	sb.cache = cache
	sb.httpClient = &http.Client{
		Timeout: defaultTimeout,
	}
	return sb
}

// LoadPosts loads all existing posts from external endpoint
func (sb *SeaBackend) LoadPosts(ctx context.Context) ([]RemotePost, error) {
	var remotePosts []RemotePost

	err := sb.load(ctx, sb.endpoint+"/posts", &remotePosts)
	if err != nil {
		return remotePosts, fmt.Errorf("could not load posts: %w", err)
	}

	return remotePosts, nil
}

// LoadUsers loads all existing users from external endpoint
func (sb *SeaBackend) LoadUsers(ctx context.Context) ([]RemoteUser, error) {
	var remoteUsers []RemoteUser
	err := sb.load(ctx, sb.endpoint+"/users", &remoteUsers)
	if err != nil {
		return remoteUsers, fmt.Errorf("could not load users: %w", err)
	}

	return remoteUsers, nil
}

// LoadUsers loads all existing users from external endpoint
func (sb *SeaBackend) LoadUser(ctx context.Context, id string) (RemoteUser, error) {
	var remoteUsers []RemoteUser
	var user RemoteUser

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
	ctxTimeout, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err = sb.cache.Get(requestUrl, data)
	if err == nil {
		return err
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

	err = sb.cache.Set(requestUrl, data)
	if err != nil {
		return fmt.Errorf("failed to save data to cache: %w", err)
	}

	return nil
}

func (sb *SeaBackend) Posts(ctx context.Context, req *web.Request) web.Result {
	var err error

	remotePosts, err := sb.LoadPosts(ctx)
	if err != nil {
		return sb.responder.ServerError(err)
	}

	filter, _ := req.Query1("filter")

	responsePosts := make([]Post, 0)
	remotePostsChan := make(chan RemotePost)
	responsePostsChan := make(chan Post)
	loadUserFunc := func(workerId int, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for remotePost := range remotePostsChan {
			user, err := sb.LoadUser(ctx, remotePost.UserID.String())
			if err != nil {
				continue
			}

			post := Post{
				Title:       remotePost.Title,
				Body:        remotePost.Body,
				Username:    user.Username,
				CompanyName: user.Company.Name,
			}

			responsePostsChan <- post
		}
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < workerCount; i++ {
		go loadUserFunc(i, wg)
	}

	responsePostEnded := make(chan struct{})
	go func() {
		for post := range responsePostsChan {
			responsePosts = append(responsePosts, post)
		}
		responsePostEnded <- struct{}{}
	}()

	for _, remotePost := range remotePosts {
		if !remotePost.Contains(filter, FieldTitle) {
			continue
		}
		remotePostsChan <- remotePost
	}
	close(remotePostsChan)

	wg.Wait()
	close(responsePostsChan)
	<-responsePostEnded

	return sb.responder.Data(responsePosts)
}
