package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/obriena/dockerdevtemplate/domain"
	"github.com/obriena/dockerdevtemplate/infra"
)

func RetrievePostById(w http.ResponseWriter, r *http.Request) {
	log.Println("controllers.RetrievePostByID Start")

	ctx := r.Context()
	postInteractor := ctx.Value(infra.CtxPostInteractorKey).(domain.PostInteractor)
	postId := ctx.Value(infra.CtxPostIdKey).(int)

	start := time.Now()
	postInput := domain.PostInteractorReadByIdInput{Id: postId}
	posts, err := postInteractor.ReadById(ctx, postInput)
	if err != nil {
		log.Println("Error retrieving posts: ", err)
		infra.RespondJSON(w, r, domain.Post{})
	} else {
		infra.RespondJSON(w, r, posts)
	}
	elapsed := time.Since(start)

	log.Println("controllers.RetrievePostByID end elapsed time: ", elapsed.Milliseconds())
}

func RetrievePostsByOwnerId(w http.ResponseWriter, r *http.Request) {
	log.Println("controllers.RetrievePostOwnerById Start")
	// split := strings.Split(r.URL.Path, "/")
	// postId := split[3]
	// i, err := strconv.Atoi(postId)
	// if err != nil {
	// 	log.Println(r.URL.Query())
	// 	log.Println("controllers.RetrievePostOwnerById: Unable to handle requested id: ", postId)
	// }

	ctx := r.Context()
	postInteractor := ctx.Value(infra.CtxPostInteractorKey).(domain.PostInteractor)
	postId := ctx.Value(infra.CtxPostIdKey).(int)

	start := time.Now()
	postInput := domain.PostInteractorReadByCreatedIdInput{Id: postId}

	posts, err := postInteractor.ReadByCreatedById(ctx, postInput)
	if err != nil {
		log.Println("Error retrieving posts: ", err)
		infra.RespondJSON(w, r, domain.Post{})
	} else {
		infra.RespondJSON(w, r, posts)
	}
	elapsed := time.Since(start)

	log.Println("controllers.RetrievePostByID end elapsed time: ", elapsed.Milliseconds())
}
