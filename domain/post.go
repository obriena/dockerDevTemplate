package domain

import (
	"context"
	"time"
)

type Post struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	OwnerId   int       `json:"ownerId"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt time.Time
	DeletedBy int
	UpdatedAt time.Time
	UpdatedBy int
}

type ListInfo struct {
	Total    int
	PageSize int
	Page     int
	NextPage int
}

type PostRepo interface {
	Create(context.Context, Post) (Post, error)
	ReadAll(context.Context) ([]Post, error)
	ReadById(context.Context, int) (Post, error)
	ReadByOwnerId(context.Context, int) ([]Post, error)
	Update(context.Context, Post) (Post, error)
	Delete(context.Context, int) error
}

type PostInteractor interface {
	Add(context.Context, PostInteractorAddInput) (PostInteractorAddOutput, error)
	ReadAll(context.Context) (PostOutput, error)
	ReadById(context.Context, PostInteractorReadByIdInput) (PostInteractorReadByIdOutput, error)
	ReadByCreatedById(context.Context, PostInteractorReadByCreatedIdInput) (PostInteractorReadByCreatedByIdOutput, error)
	Delete(context.Context, PostInteractorDeleteInput) error
	Update(context.Context, Post) (PostInteractorUpdateOutput, error)
}

type PostInteractorAddInput struct {
	Content   string
	Title     string
	OwnerId   int
	CreatedBy int
}

type PostInteractorAddOutput struct {
	Post Post
}

type PostInteractorReadByIdInput struct {
	Id int
}

type PostInteractorReadByCreatedIdInput struct {
	Id int
}

type PostOutput struct {
	Post     []Post
	ListInfo ListInfo
}
type PostInteractorReadByIdOutput struct {
	Post Post
}

type PostInteractorReadByCreatedByIdOutput struct {
	Post     []Post
	ListInfo ListInfo
}

type PostInteractorDeleteInput struct {
	Id int
}

type PostInteractorUpdateOutput struct {
	Post Post
}
