package file

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

type FileStore struct {
	bucket *storage.BucketHandle
}

func NewFileStore(bucket *storage.BucketHandle) *FileStore {
	return &FileStore{
		bucket: bucket,
	}
}

func generateRandomFileName(userId uint) string {
	const fileNameLength = 16
	timestamp := time.Now().UnixNano()                   // Using UnixNano for high precision
	randomString := generateRandomString(fileNameLength) // Generate a random string of 8 characters
	return fmt.Sprintf("%d-%s-%d", userId, randomString, timestamp)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomBytes)
}

func (fileStore *FileStore) UploadFile(userId uint, file multipart.File) (string, string, error) {
	randomFileName := generateRandomFileName(userId)

	object := fileStore.bucket.Object(randomFileName)
	writer := object.NewWriter(context.Background())

	defer file.Close()
	_, err := io.Copy(writer, file)
	if err != nil {
		return "", "", err
	}
	err = writer.Close()
	if err != nil {
		return "", "", err
	}

	attrs, err := object.Attrs(context.Background())
	if err != nil {
		return "", "", err
	}
	fileURL := attrs.MediaLink
	objectName := object.ObjectName()

	return objectName, fileURL, nil
}

func (fileStore *FileStore) DeleteFile(objectName string) error {
	object := fileStore.bucket.Object(objectName)
	return object.Delete(context.Background())
}
