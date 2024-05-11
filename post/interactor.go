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

func (i *Interactor) Add(context context.Context, input domain.PostInteractorAddInput) (domain.PostInteractorAddOutput, error) {
	log.Println("post.interactor.Add start")
	post := domain.Post{Content: input.Content,
		Title:   input.Title,
		OwnerId: input.OwnerId,
	}
	i.repo.Create(context, post)
	log.Println("post.interactor.Add complete")
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

func (i *Interactor) ReadById(context context.Context, readInput domain.PostInteractorReadByIdInput) (domain.PostInteractorReadByIdOutput, error) {
	post, err := i.repo.ReadById(context, readInput.Id)
	if err != nil {
		return domain.PostInteractorReadByIdOutput{}, err
	}
	output := domain.PostInteractorReadByIdOutput{Post: post}
	return output, nil
}

func (i *Interactor) ReadByCreatedById(context context.Context, readInput domain.PostInteractorReadByCreatedIdInput) (domain.PostInteractorReadByCreatedByIdOutput, error) {
	log.Println("post.Interactor: Reading ReadByCreatedById: ", readInput.Id)
	post, err := i.repo.ReadByOwnerId(context, readInput.Id)
	if err != nil {
		log.Println("post/Interactor: Error reading all posts")
		return domain.PostInteractorReadByCreatedByIdOutput{}, err
	}

	log.Println("post.Interactor: retrieved posts by creator id", len(post))
	li := domain.ListInfo{Total: len(post), PageSize: 1, Page: 1, NextPage: 1}
	po := domain.PostInteractorReadByCreatedByIdOutput{Post: post, ListInfo: li}
	return po, nil
}

func (i *Interactor) Delete(context context.Context, deleteInput domain.PostInteractorDeleteInput) error {
	log.Println("post.Interactor: Delete: ", deleteInput.Id)
	err := i.repo.Delete(context, deleteInput.Id)
	if err != nil {
		log.Println("post/Interactor: Error deleting post ", deleteInput.Id)
		return err
	}

	log.Println("post.Interactor: Deleted post", deleteInput.Id)
	return nil
}

func (i *Interactor) Update(context context.Context, post domain.Post) (domain.PostInteractorUpdateOutput, error) {
	log.Println("post.Interactor: Update: ", post.Id)
	newPost, err := i.repo.Update(context, post)

	if err != nil {
		log.Println("post/Interactor: Error deleting post ", post.Id)
		return domain.PostInteractorUpdateOutput{}, err
	}
	output := domain.PostInteractorUpdateOutput{Post: newPost}
	log.Println("post.Interactor: Updated post")
	return output, nil
}
