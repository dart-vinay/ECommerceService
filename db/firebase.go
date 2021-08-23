package db

import (
	"golang.org/x/net/context"
	"log"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

func GetFirebaseClient() (*context.Context, *firebase.App) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./resources/locklly-5692a-firebase-adminsdk-q079a-be8c6cfda3.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error Creating Firebase Client %s", )
	}
	return &ctx, app
}
