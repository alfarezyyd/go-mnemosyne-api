package config

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
)

type GoogleCloudStorage struct {
	StorageClient *storage.Client
}

func InitializeGoogleCloudStorage() (*GoogleCloudStorage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("mnemosyne-labs.json"))
	if err != nil {
		return nil, err
	}
	return &GoogleCloudStorage{
		StorageClient: client,
	}, nil
}
