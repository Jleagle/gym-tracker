package datastore

import (
	"context"
	"math/rand"
	"path/filepath"

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

	abs, err := filepath.Abs("./")
	if err != nil {
		log.Instance.Error("abs", zap.Error(err))
	}

	ctx = context.Background()
	client, err = datastore.NewClient(ctx, config.GoogleProject, option.WithCredentialsFile(abs+"/gcp-auth.json"))
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

	_, err := client.Put(ctx, key, &data)
	return err
}

func GetCredentials() (ret []Credential, err error) {

	query := datastore.NewQuery("Credential") // Grab whole table

	var creds []Credential
	_, err = client.GetAll(ctx, query, &creds)
	if err != nil {
		return nil, err
	}

	// Group by gym and email
	credsMap := map[string]map[string]Credential{}
	for _, v := range creds {
		if credsMap[v.Gym] == nil {
			credsMap[v.Gym] = map[string]Credential{}
		}
		credsMap[v.Gym][v.Email] = v
	}

	creds = []Credential{}
	for _, gymCreds := range credsMap {
		var s []Credential
		for _, vv := range gymCreds {
			s = append(s, vv)
		}
		if len(s) > 0 {
			ret = append(ret, s[rand.Intn(len(s))])
		}
	}

	return ret, nil
}
