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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	log.Println("controllers.DeletePost Start")
	split := strings.Split(r.URL.Path, "/")
	postId := split[2]
	i, err := strconv.Atoi(postId)
	if err != nil {
		log.Println(r.URL.Query())
		log.Println("controllers.DeletePost: Unable to handle requested id: ", postId)
	}

	ctx := r.Context()
	postInteractor := ctx.Value(infra.CtxPostInteractorKey).(domain.PostInteractor)

	start := time.Now()
	postInput := domain.PostInteractorDeleteInput{Id: i}
	error := postInteractor.Delete(ctx, postInput)
	if error != nil {
		log.Println("Error retrieving posts: ", error)
		infra.RespondJSON(w, r, domain.Post{})
	} else {
		message := make(map[string]string)
		message["message"] = "Post Deleted"
		message["success"] = "true"
		infra.RespondJSON(w, r, message)
	}
	elapsed := time.Since(start)

	log.Println("controllers.RetrievePostByID end elapsed time: ", elapsed.Milliseconds())
}
