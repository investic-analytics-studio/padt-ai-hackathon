package infra

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/quantsmithapp/datastation-backend/config"
	"google.golang.org/api/option"
)

var FirebaseClient *auth.Client

func InitFirebaseClient() {
	ctx := context.Background()
	cfg := config.GetConfig().Firebase
	credential := []byte(cfg.Credential)

	clientOption := option.WithCredentialsJSON(credential)
	app, err := firebase.NewApp(ctx, nil, clientOption)
	if err != nil {
		panic(err)
	}

	FirebaseClient, err = app.Auth(ctx)
	if err != nil {
		panic(err)
	}
}
