package repositories

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func InitFirestore() {
	ctx := context.Background()
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	
	if projectID == "" {
		projectID = "itstep-social" // Fallback
	}

	// Initialize Firebase App
	// In production, use GOOGLE_APPLICATION_CREDENTIALS env var
	// point to serviceAccountKey.json
	conf := &firebase.Config{ProjectID: projectID}
	
	var app *firebase.App
	var err error

	// Try to load service account if provided
	saPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if saPath != "" {
		opt := option.WithCredentialsFile(saPath)
		app, err = firebase.NewApp(ctx, conf, opt)
	} else {
		app, err = firebase.NewApp(ctx, conf)
	}

	if err != nil {
		log.Fatalf("Failed to initialize firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize firestore client: %v", err)
	}

	FirestoreClient = client
	log.Println("Firestore connection established")
}

var (
	Descending = firestore.Desc
	MergeAll   = firestore.MergeAll
)
