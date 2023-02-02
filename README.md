# Jokes API with GO and Firestore
This is a simple REST API to perform CRUD operations with jokes using the GO HTTP package and Firestore. The purpose is to provide a simple example of a CRUD app that uses the standard `net/http` library, interacts with Firestore, and follows a common API structure.

## Prerequisites
The only prerequisite to using this example is having GO installed on your machine. A working knowledge of Firestore and REST APIs is useful, but this project has been set up to work out of the box.

## Getting Started
1. Clone the repo with `git clone https://github.com/imbradyboy/go-firestore-crud.git`
2. Create a new Firebase project and download a service account credential. If you've never done this before, you can find official instructions [here](https://firebase.google.com/docs/admin/setup#set-up-project-and-service-account)
3. Create a `.env` file in the root of your project and add the absolute path to the service account you downloaded in the step above to a variable named `FB_ADMIN_SA_LOCATION` like below
```
FB_ADMIN_SA_LOCATION=path/to/file.json
```
4. Run `go mod download` to install dependencies
5. Open up PostMan and import the collection `go-firestore-jokes.postman_collection.json` at the root of the project
6. Run `go run cmd/main.go` from the root of the project


