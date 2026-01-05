package firebase

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var app *firebase.App
var ctx = context.Background()

func Initailize() *firebase.App {
	if app != nil {
		return app
	}

	//need to change path for test and production
	// remember it barunn for testing "GOOGLE_APPLICATION_CREDENTIALS" for development "Location"
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	var err error
	app, err = firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_ID")}, opt)
	if err != nil {
		log.Fatalf("failed to intailize firebase %v", err)
	} else {
		fmt.Printf("firebase intailized")
	}
	return app
}

// firestore
func Firestore() *firestore.Client {
	if app == nil {
		Initailize()
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatalf("failed to create firestore client :%v", err)
	}
	return client
}
