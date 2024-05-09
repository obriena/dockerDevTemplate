package post

import (
	"context"
	"log"

	"github.com/obriena/dockerdevtemplate/domain"
)

type Interactor struct {
	repo domain.PostRepo
}

func NewInteractor(repo domain.PostRepo) *Interactor {
	return &Interactor{
		repo: repo,
	}
}

func (i *Interactor) Add(context.Context, domain.PostInteractorAddInput) (domain.PostInteractorAddOutput, error) {
	return domain.PostInteractorAddOutput{}, nil
}

func (i *Interactor) ReadAll(context context.Context) (domain.PostOutput, error) {
	log.Println("post.Interactor: Reading all posts")
	post, err := i.repo.ReadAll(context)
	if err != nil {
		log.Println("post/Interactor: Error reading all posts")
		return domain.PostOutput{}, err
	}

	log.Println("post.Interactor: retrieved posts", len(post))
	log.Println("post.Interactor: posts", post)
	li := domain.ListInfo{Total: len(post), PageSize: 1, Page: 1, NextPage: 1}
	po := domain.PostOutput{Post: post, ListInfo: li}
	return po, nil
}

func (i *Interactor) ReadById(context.Context, domain.PostInteractorReadByIdInput) (domain.PostInteractorReadByIdOutput, error) {
	return domain.PostInteractorReadByIdOutput{}, nil
}

func (i *Interactor) ReadByCreatedById(context.Context, domain.PostInteractorReadByCreatedIdInput) (domain.PostInteractorReadByCreatedByIdOutput, error) {
	return domain.PostInteractorReadByCreatedByIdOutput{}, nil
}
