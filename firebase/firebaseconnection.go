package firebase

import (
	"context"
	"fmt"
	"log"

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

	opt := option.WithCredentialsFile("firebase/serviceAccountkey.json")
	var err error
	app, err = firebase.NewApp(ctx, nil, opt)
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
