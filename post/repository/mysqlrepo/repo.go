package mysqlrepo

import (
	"context"
	"log"
	"time"

	"github.com/obriena/dockerdevtemplate/domain"
	"gorm.io/gorm"
)

const table = "post"

type postModel struct {
	gorm.Model
	Id        int
	Content   string
	Title     string
	OwnerId   int
	CreatedAt time.Time
	CreatedBy int
	UpdatedAt time.Time
	UpdatedBy int
	deletedAt time.Time
	deletedBy int
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(context.Context, domain.Post) (domain.Post, error) {
	return domain.Post{}, nil
}

func (r *Repository) ReadAll(context context.Context) ([]domain.Post, error) {
	log.Println("post.repo: Reading all posts")
	allPosts := []postModel{}

	err := r.db.Table(table).Find(&allPosts).Error
	if err != nil {
		log.Println("Error retrieving data: ", err)
		return nil, err
	}
	postList := toDomainList(allPosts)
	log.Println("post.repo: Done Reading all posts", postList)

	return postList, nil
}

func (r *Repository) ReadById(context context.Context, postId int) (domain.Post, error) {
	return domain.Post{}, nil
}

func (r *Repository) Update(context.Context, domain.Post) (domain.Post, error) {
	return domain.Post{}, nil
}

func toDomainList(p []postModel) []domain.Post {
	var domainPostList = make([]domain.Post, len(p))
	for i := 0; i < len(p); i++ {
		domainPostList[i] = p[i].toDomain()
	}
	return domainPostList
}
func (p *postModel) toDomain() domain.Post {
	return domain.Post{
		Id:        p.Id,
		Content:   p.Content,
		Title:     p.Title,
		OwnerId:   p.OwnerId,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		DeletedAt: p.deletedAt,
		DeletedBy: p.deletedBy,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy,
	}
}
