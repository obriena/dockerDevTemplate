package domain

import (
	"context"
	"time"
)

type Post struct {
	Id        int
	Content   string
	Title     string
	OwnerId   int
	Deleted   bool
	CreatedAt time.Time
	CreatedBy int
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
	Update(context.Context, Post) (Post, error)
}

type PostInteractor interface {
	Add(context.Context, PostInteractorAddInput) (PostInteractorAddOutput, error)
	ReadAll(context.Context) (PostOutput, error)
	ReadById(context.Context, PostInteractorReadByIdInput) (PostInteractorReadByIdOutput, error)
	ReadByCreatedById(context.Context, PostInteractorReadByCreatedIdInput) (PostInteractorReadByCreatedByIdOutput, error)
}

type PostInteractorAddInput struct {
	Content string
	Title   string
	OwnerId int
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
	Post Post
}
