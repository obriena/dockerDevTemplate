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
	ID        int `gorm:"primarykey;size:16"`
	Content   string
	Title     string
	Language  string
	Deleted   bool
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

func (r *Repository) Create(context context.Context, post domain.Post) (domain.Post, error) {
	dbModel := postModel{
		Content: post.Content,
		Title:   post.Title,
		OwnerId: post.OwnerId,
	}
	tx := r.db.Table(table).Create(&dbModel)
	if tx.Error != nil {
		log.Println("Error creatint post", tx.Error)
	}
	log.Println("New Post ID: ", dbModel.ID)
	return domain.Post{Id: dbModel.ID}, nil
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
	log.Println("post.repo: Reading posts for id: ", postId)
	aPost := postModel{}

	err := r.db.Table(table).Where("id = ?", postId).First(&aPost).Error
	if err != nil {
		log.Println("Error retrieving data: ", err)
		return domain.Post{}, err
	}
	dPost := aPost.toDomain()
	log.Println("post.repo: Done Reading all posts", dPost)

	return dPost, nil
}
func (r *Repository) ReadByOwnerId(context context.Context, ownerId int) ([]domain.Post, error) {
	log.Println("post.repo: ReadByOwnerId posts: ", ownerId)
	allPosts := []postModel{}

	err := r.db.Table(table).Where("owner_id = ?", ownerId).Find(&allPosts).Error
	if err != nil {
		log.Println("Error retrieving data: ", err)
		return nil, err
	}
	postList := toDomainList(allPosts)
	log.Println("post.repo: Done Reading post by owner id", postList)

	return postList, nil
}

func (r *Repository) Update(context context.Context, post domain.Post) (domain.Post, error) {
	log.Println("post.repo: Updating post with id: ", post.Id)

	postModel := postModel{
		ID:        post.Id,
		Content:   post.Content,
		Title:     post.Title,
		Language:  post.Language,
		UpdatedAt: post.UpdatedAt,
		UpdatedBy: post.UpdatedBy,
	}

	res := r.db.Table("post").Model(&postModel).Updates(postModel)
	log.Println("update response: ", res.Error)
	log.Println("update response rows affected: ", res.RowsAffected)

	return post, nil
}

func (r *Repository) Delete(ctx context.Context, postId int) error {
	log.Println("post.repo: Deleting post with id: ", postId)

	res := r.db.Table(table).Where("id = ?", postId).Update("deleted", true)
	log.Println("delete response: ", res.Error)
	log.Println("delete response rows affected: ", res.RowsAffected)

	return nil
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
		Id:        p.ID,
		Content:   p.Content,
		Title:     p.Title,
		Language:  p.Language,
		OwnerId:   p.OwnerId,
		Deleted:   p.Deleted,
		CreatedAt: p.CreatedAt,
		DeletedAt: p.deletedAt,
		DeletedBy: p.deletedBy,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy,
	}
}
