package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/imbradyboy/go-firestore-crud/pkg/models"
	"github.com/imbradyboy/go-firestore-crud/pkg/utils"
	"google.golang.org/grpc/codes"
)

func GetAllJokes(w http.ResponseWriter, r *http.Request) {
	// fetch all jokes from model
	sliceOfJokes, err := models.GetAllJokes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// encode slice to json and send response
	res, _ := json.Marshal(sliceOfJokes)
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetJokeById(w http.ResponseWriter, r *http.Request) {
	// get id from path params
	id := utils.GetIdFromParams(r.URL.Path)

	// get joke from model
	joke, err := models.GetJokeById(id, r.Context())
	if err != nil {
		// check if returned not found error, otherwise send back generic error
		httpStatus := http.StatusBadRequest
		if err.Error() == codes.NotFound.String() {
			httpStatus = http.StatusNotFound
		}

		http.Error(w, err.Error(), httpStatus)
		return
	}

	// encode to json and send response back
	res, _ := json.Marshal(joke)
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func AddJoke(w http.ResponseWriter, r *http.Request) {
	// decode json body into joke struct
	var newJoke models.JokeDTO
	err := json.NewDecoder(r.Body).Decode(&newJoke)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add joke to db
	addedJoke, err := models.AddJoke(newJoke, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// encode response to json and send back
	res, _ := json.Marshal(addedJoke)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateJoke(w http.ResponseWriter, r *http.Request) {
	// get id from path params
	id := utils.GetIdFromParams(r.URL.Path)

	// decode json body into joke struct
	var updatedJoke models.JokeDTO
	err := json.NewDecoder(r.Body).Decode(&updatedJoke)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update specified joke in db
	updatedJokeResponse, err := models.UpdateJoke(id, updatedJoke, r.Context())
	if err != nil {
		// check for not found error on update, otherwise return generic error
		httpStatus := http.StatusBadRequest
		if err.Error() == codes.NotFound.String() {
			httpStatus = http.StatusNotFound
		}
		http.Error(w, err.Error(), httpStatus)
		return
	}

	// encode response to json and send back
	res, _ := json.Marshal(updatedJokeResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteJoke(w http.ResponseWriter, r *http.Request) {
	// get id from path params
	id := utils.GetIdFromParams(r.URL.Path)

	// delete joke from db
	deleteResult, err := models.DeleteJoke(id, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// encode respones to json and send back
	res, _ := json.Marshal(deleteResult)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
