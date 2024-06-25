package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

var storageClient *storage.Client

// InitializeFirebaseApp initializes the Firebase app and storage client
func InitializeFirebaseApp(credentialFilePath string) error {
	opt := option.WithCredentialsFile(credentialFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return err
	}

	storageClient, err = app.Storage(context.Background())
	if err != nil {
		log.Fatalf("error getting Storage client: %v\n", err)
		return err
	}

	return nil
}

// UploadFileToFirebase uploads a file to Firebase Storage and returns the file URL
func UploadFileToFirebase(ctx context.Context, bucketName string, file multipart.File, fileName string) (string, error) {
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("error getting bucket: %v", err)
	}

	object := bucket.Object(fileName)
	wc := object.NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return url, nil
}
