package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/obriena/dockerdevtemplate/domain"
	"github.com/obriena/dockerdevtemplate/infra"
)

func RetrievePostById(w http.ResponseWriter, r *http.Request) {
	log.Println("controllers.RetrievePostByID Start")
	split := strings.Split(r.URL.Path, "/")
	postId := split[2]
	i, err := strconv.Atoi(postId)
	if err != nil {
		log.Println(r.URL.Query())
		log.Println("controllers.RetrievePostById: Unable to handle requested id: ", postId)
	}

	ctx := r.Context()
	postInteractor := ctx.Value(infra.CtxPostInteractorKey).(domain.PostInteractor)

	start := time.Now()
	postInput := domain.PostInteractorReadByIdInput{Id: i}
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
