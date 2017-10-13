package main

import (
	"github.com/n-marshall/oauth2ns"
	"golang.org/x/oauth2"
)

func main() {
	conf := &oauth2.Config{
		ClientID:     "tRQ9V3cyLxHu3xXFWM",               // also known as slient key sometimes
		ClientSecret: "PZEU9ruk3eZxYvAeCRgS9YyYpAybLG4P", // also known as secret key
		Scopes:       []string{"account"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://bitbucket.org/site/oauth2/authorize",
			TokenURL: "https://bitbucket.org/site/oauth2/access_token",
		},
	}
	/*r := ...*/ _ = oauth2ns.Authorize(conf)
	// use r.client.Get / r.client.Post for further requests, the token will automatically be there
}
