package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/obriena/dockerdevtemplate/domain"
	"github.com/obriena/dockerdevtemplate/infra"
)

func SavePost(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("controllers.SavePost Receiving post data")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Create a struct to unmarshal the JSON into
	var jsonObj domain.Post

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	postInteractor := ctx.Value(infra.CtxPostInteractorKey).(domain.PostInteractor)
	postInput := domain.PostInteractorAddInput{
		Content: jsonObj.Content,
		Title:   jsonObj.Title,
		OwnerId: 1,
	}

	postInteractor.Add(ctx, postInput)

	elapsed := time.Since(start)
	log.Println("controllers.SavePost complete elapsed time. ", elapsed.Milliseconds())
}
