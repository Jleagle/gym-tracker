package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36"
)

var (
	loginURL *url.URL
	countURL *url.URL

	client    *http.Client
	cookieJar *cookiejar.Jar
)

func init() {

	var err error

	loginURL, err = url.Parse("https://www.puregym.com/api/members/login/")
	if err != nil {
		log.Fatal(err)
	}

	countURL, err = url.Parse("https://www.puregym.com/members/")
	if err != nil {
		log.Fatal(err)
	}

	cookieJar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client = &http.Client{
		Jar: cookieJar,

		// Don't follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func main() {

	log.Println("Trying to get count..")
	req, err := http.NewRequest("GET", countURL.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	log.Println("Need to login first..")
	if resp.StatusCode != 200 {
		err = login(client)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func login(client *http.Client) (err error) {

	loginPayload := LoginPayload{
		AssociateAccount: "false",
		Email:            os.Getenv("PUREGYM_EMAIL"),
		Pin:              os.Getenv("PUREGYM_PIN"),
	}

	loginPayloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", loginURL.String(), bytes.NewReader(loginPayloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		loginResponse := LoginResponse{}
		err = json.Unmarshal(b, &loginResponse)
		if err != nil {
			return err
		}

		return errors.New(loginResponse.Message)
	}

	return nil
}

type LoginPayload struct {
	AssociateAccount string `json:"associateAccount"`
	Email            string `json:"email"`
	Pin              string `json:"pin"`
}

type LoginResponse struct {
	Message string `json:"message"`
}
