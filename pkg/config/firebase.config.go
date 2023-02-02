package config

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// global instance of db
var FirestoreDb *firestore.Client

func InitializeFirebaseApp() {
	ctx := context.Background()

	// create Firebase app with service account
	sa := option.WithCredentialsFile(os.Getenv("FB_ADMIN_SA_LOCATION"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		panic("Firebase app did not initialize properly")
	}

	// create Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		panic("Firestore did not initialize properly")
	}

	FirestoreDb = client

	// do not need to close client if required for entire app as stopping the app will cleanup automatically
	// defer client.Close()
}
