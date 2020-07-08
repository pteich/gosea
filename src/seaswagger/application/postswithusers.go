package application

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pteich/gosea/src/seaswagger/domain/entity"
	"github.com/pteich/gosea/src/seaswagger/domain/service"
	"github.com/pteich/gosea/src/seaswagger/infrastructure/models"
)

type postsWithUsersCfg struct {
	WorkerCount float64 `inject:"config:api.workerCount"`
}

type PostsWithUsers struct {
	seaBackend  service.SeaBackendLoader
	workerCount int
}

func (p *PostsWithUsers) Inject(
	seaBackend service.SeaBackendLoader,
	cfg *postsWithUsersCfg,
) {
	p.seaBackend = seaBackend
	p.workerCount = int(cfg.WorkerCount)
}

func (p *PostsWithUsers) RetrievePostsWithUsersFromBackend(ctx context.Context, filter string) ([]entity.Post, error) {
	responsePosts := make([]entity.Post, 0)

	remotePosts, err := p.seaBackend.LoadPosts(ctx)
	if err != nil {
		return responsePosts, err
	}

	remotePostsChan := make(chan *models.Post)
	responsePostsChan := make(chan entity.Post)
	loadUserFunc := func(workerId int, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for remotePost := range remotePostsChan {
			userId, ok := remotePost.UserID.(string)
			if !ok {
				userIdNumber, ok := remotePost.UserID.(json.Number)
				if !ok {
					return
				}
				userId = userIdNumber.String()
			}
			user, err := p.seaBackend.LoadUser(ctx, userId)
			if err != nil {
				continue
			}

			post := entity.Post{
				Title:       remotePost.Title,
				Body:        remotePost.Body,
				Username:    user.Username,
				CompanyName: user.Company.Name,
			}

			responsePostsChan <- post
		}
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < p.workerCount; i++ {
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
		/*
			if !remotePost.Contains(filter, entity.FieldTitle) {
				continue
			}
		*/
		remotePostsChan <- remotePost
	}
	close(remotePostsChan)

	wg.Wait()
	close(responsePostsChan)
	<-responsePostEnded

	return responsePosts, nil
}
