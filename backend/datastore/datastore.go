package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/gym-tracker/config"
	"github.com/Jleagle/gym-tracker/log"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var (
	client *datastore.Client
	ctx    context.Context
)

func init() {

	ctx = context.Background()

	var err error
	client, err = datastore.NewClient(ctx, config.DatastoreProject, option.WithAPIKey(config.DatastoreKey))
	if err != nil {
		log.Instance.Error("create datastore client", zap.Error(err))
	}
}

type Credential struct {
	Email string `datastore:"email,noindex"`
	PIN   string `datastore:"pin,noindex"`
	Gym   string `datastore:"gym"`
}

func SaveNewCredential(email, pin, gym string) error {

	key := datastore.IncompleteKey("Credential", nil)
	data := Credential{Email: email, PIN: pin, Gym: gym}

	_, err := client.Put(ctx, key, data)
	return err
}

func GetCredentials() (credsMap map[string][]Credential, err error) {

	query := datastore.NewQuery("Credential")

	var creds []Credential
	_, err = client.GetAll(ctx, query, &creds)
	if err != nil {
		return nil, err
	}

	credsMap = map[string][]Credential{}
	for _, v := range creds {
		credsMap[v.Gym] = append(credsMap[v.Gym], v)
	}

	return credsMap, nil
}
