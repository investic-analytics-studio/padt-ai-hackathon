package infra

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"google.golang.org/api/option"
)

var CloudStorage *storage.Client

func InitCloudStorage() {
	var err error
	CloudStorage, err = storage.NewClient(context.Background(), option.WithCredentialsFile("gcloud.json"))
	if err != nil {
		logger.Fatalf("InitCloudStorage(): %v", err)
	}
}
