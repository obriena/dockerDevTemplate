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

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("controllers.Update Receiving post data")

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
	postInput := domain.Post{
		Id:        jsonObj.Id,
		Content:   jsonObj.Content,
		Title:     jsonObj.Title,
		Language:  jsonObj.Language,
		OwnerId:   1,
		UpdatedBy: 1,
		UpdatedAt: time.Now(),
	}

	postInteractor.Update(ctx, postInput)

	elapsed := time.Since(start)
	log.Println("controllers.SavePost complete elapsed time. ", elapsed.Milliseconds())
}
