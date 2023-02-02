package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/imbradyboy/go-firestore-crud/pkg/config"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// STRUCTS

// shape of joke object created/passed in by user
type JokeDTO struct {
	Joke      string `json:"joke,omitempty" firestore:"joke,omitempty"`
	Punchline string `json:"punchline,omitempty" firestore:"punchline,omitempty"`
}

// shape of joke returned to user so useful metadata is shown
type JokeRO struct {
	JokeDTO
	ID         string    `json:"id"`
	CreateTime time.Time `json:"dateCreated"`
	UpdateTime time.Time `json:"dateUpdated"`
}

// CUSTOM ERRORS
type FirestoreError struct {
	When          time.Time
	CustomError   error
	OriginalError error
}

var (
	errDocNotFound = errors.New("joke not found")
	errBadRequest  = errors.New("something went wrong")
)

// VARIABLES
const collection_name string = "jokes"

// FUNCTIONS
func (e *FirestoreError) Error() string {
	// log out actual error message, but return only the code to user
	fmt.Printf("%v error %v -> returned as %v\n", e.When, e.OriginalError.Error(), status.Code(e.OriginalError))
	return status.Code(e.OriginalError).String()
}

func GetAllJokes(ctx context.Context) ([]JokeRO, error) {
	// get docs in the collection
	iter := config.FirestoreDb.Collection(collection_name).Documents(ctx)

	// slice to hold all doc data
	sliceOfJokes := []JokeRO{}

	// iterate over each doc
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// return custom error
			return nil, &FirestoreError{
				When:          time.Now(),
				OriginalError: err,
				CustomError:   errBadRequest,
			}
		}

		// read data into joke struct format
		var getJoke JokeRO
		doc.DataTo(&getJoke)
		getJoke.ID = doc.Ref.ID
		getJoke.CreateTime = doc.CreateTime
		getJoke.UpdateTime = doc.UpdateTime

		// add to slice of jokes
		sliceOfJokes = append(sliceOfJokes, getJoke)
	}

	return sliceOfJokes, nil
}

func GetJokeById(id string, ctx context.Context) (JokeRO, error) {
	// get doc
	dsnap, err := config.FirestoreDb.Collection(collection_name).Doc(id).Get(ctx)
	if err != nil {
		// custom error
		return JokeRO{}, &FirestoreError{
			When:          time.Now(),
			OriginalError: err,
			CustomError:   errBadRequest,
		}
	}

	// read doc data into struct
	var getJoke JokeRO
	dsnap.DataTo(&getJoke)
	getJoke.ID = dsnap.Ref.ID
	getJoke.CreateTime = dsnap.CreateTime
	getJoke.UpdateTime = dsnap.UpdateTime

	return getJoke, nil
}

func AddJoke(newJoke JokeDTO, ctx context.Context) (JokeRO, error) {
	// add doc to collection
	ref, result, err := config.FirestoreDb.Collection(collection_name).Add(ctx, newJoke)
	// custom error
	if err != nil {
		return JokeRO{}, &FirestoreError{
			When:          time.Now(),
			OriginalError: err,
			CustomError:   errBadRequest,
		}
	}

	// construct success response for user
	addedJoke := JokeRO{
		JokeDTO:    newJoke,
		ID:         ref.ID,
		CreateTime: result.UpdateTime,
		UpdateTime: result.UpdateTime,
	}

	return addedJoke, nil
}

func UpdateJoke(id string, updatedJoke JokeDTO, ctx context.Context) (JokeRO, error) {
	// fetch doc to check if it exists
	// the GO Firestore package does not accept maps or structs in an update, only field paths and values, which is more hassle and less of a "dynamic" update. So, instead we are doing a set with merge if the doc exists
	snap, err := config.FirestoreDb.Collection(collection_name).Doc(id).Get(ctx)
	if err != nil {
		// custom errors
		// check if doc exists
		if status.Code(err) == codes.NotFound {
			return JokeRO{}, &FirestoreError{
				When:          time.Now(),
				OriginalError: err,
				CustomError:   errDocNotFound,
			}
		} else {
			// generic error
			return JokeRO{}, &FirestoreError{
				When:          time.Now(),
				OriginalError: err,
				CustomError:   errBadRequest,
			}
		}
	}

	// merge all doesnt support structs, so we have to convert the data into a map first
	var updatedJokeMap map[string]interface{}
	data, _ := json.Marshal(updatedJoke)
	json.Unmarshal(data, &updatedJokeMap)

	// update (set with merge) the doc
	result, err := config.FirestoreDb.Collection(collection_name).Doc(snap.Ref.ID).Set(ctx, updatedJokeMap, firestore.MergeAll)

	if err != nil {
		// custom error
		return JokeRO{}, &FirestoreError{
			When:          time.Now(),
			OriginalError: err,
			CustomError:   errBadRequest,
		}
	}

	// construct success response for user
	updatedJokeResponse := JokeRO{
		JokeDTO:    updatedJoke,
		ID:         snap.Ref.ID,
		CreateTime: snap.CreateTime,
		UpdateTime: result.UpdateTime,
	}

	return updatedJokeResponse, nil
}

func DeleteJoke(id string, ctx context.Context) (map[string]string, error) {
	// delete doc -> if doc does not exist it completes successfully
	_, err := config.FirestoreDb.Collection(collection_name).Doc(id).Delete(ctx)
	if err != nil {
		return map[string]string{}, err
	}

	// construct success payload for user
	deletedSucess := map[string]string{
		"id": id,
	}

	return deletedSucess, nil
}
