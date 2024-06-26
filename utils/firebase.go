package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var storageClient *storage.Client

// InitializeFirebaseApp initializes the Firebase app and storage client
func InitializeFirebaseApp(credentialFilePath string) error {
	opt := option.WithCredentialsFile(credentialFilePath)
	_, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return err
	}

	storageClient, err = storage.NewClient(context.Background(), opt)
	if err != nil {
		log.Fatalf("error getting Storage client: %v\n", err)
		return err
	}

	return nil
}

// UploadFileToFirebase uploads a file to Firebase Storage and returns the file URL
func UploadFileToFirebase(ctx context.Context, bucketName string, file multipart.File, fileName string) (string, error) {
	bucket := storageClient.Bucket(bucketName)

	object := bucket.Object(fileName)
	wc := object.NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	// Make the object publicly accessible
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("error setting object ACL: %v", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return url, nil
}

func DeleteFileFromFirebase(ctx context.Context, bucketName string, fileURL string) error {
	bucket := storageClient.Bucket(bucketName)

	objectName := fileURL[len(fmt.Sprintf("https://storage.googleapis.com/%s/", bucketName)):]
	object := bucket.Object(objectName)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}
