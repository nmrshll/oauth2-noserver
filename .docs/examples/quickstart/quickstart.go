package main

import (
	"github.com/nmrshll/oauth2-noserver"
	"golang.org/x/oauth2"
)

func main() {
	conf := &oauth2.Config{
		ClientID:     "________________",            // also known as client key sometimes
		ClientSecret: "___________________________", // also known as secret key
		Scopes:       []string{"account"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://bitbucket.org/site/oauth2/authorize",
			TokenURL: "https://bitbucket.org/site/oauth2/access_token",
		},
	}
	/*client := ...*/ _ = oauth2ns.Authorize(conf)
	// use client.Get / client.Post for further requests, the token will automatically be there
}
