package main

import (
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
)

var (
	client    *http.Client
	cookieJar *cookiejar.Jar
)

func init() {

	var err error

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

	log.Println("Getting members page..")
	req, err := http.NewRequest("GET", "https://www.puregym.com/members/", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.AddCookie(&http.Cookie{Name: ".AspNet.ApplicationCookie", Value: "jOyOArttKQTMnoAg7Eh7NnkuUwVOEb5cdNMoMqEjv-Vf4lk2V3hO9xwtGc8HWp6BFHihC8kxRGIoJ_VK-yq12z-hXYx9sj5oRqDQolKVQe2TwCvi1YAb-dsJcisqej-d14RtOBZ6myZwYxpc1xtBmzmI88sgKGvGGP3OA0lZGQs6X17YRKxeNZs4cLLuo9i9UvJnKz6rVqyFHqhdglmON7E1xbz1nC_7tX9xvWgnmbx2COmsq_Yjic2ZOOr1Uc7ftF2awz0762569FQQPdwh9UDPJlbnTHOFACABz6JbMHxYLR01Q2U0d3mnamZsImp4aM-uqCPuLUhcC9CzC8MnTvIk4O9BLTCUO3u10vbepPUR9RyB40Oto07ktXiInaEAMofgZD5zHg2q0j0pV00BLd0qNL8"})

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Request code: " + strconv.Itoa(resp.StatusCode))
	log.Println("Body len: " + strconv.Itoa(len(body)))
}
