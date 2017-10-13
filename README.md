# oauth2ns: easy oauth2 for purely client side apps (desktop GUI, cli, mobile)
While oauth works fine for apps the have a server side, I personally find it a pain to use when developing simple desktop or cli applications.  
That's why needed something to turn a complete oauth flow into a one-liner (well, that's clearly exaggerated, but it's barely more).  

So here's how it works !  

# Usage

On the web service that you want your app to authenticate into, register you app (aka client) to get a `client id` and a `client secret`. 

**IMPORTANT**: you must set the redirection URL to `http://127.0.0.1:14565/oauth/callback` for this lib to function properly.  

Here's an example of creating an app on bitbucket. The UI is usually similar on other web services.  

![alt text](./.readme/creating-oauth-apps.png "app creation parameters")

This will give you the `client id` and `client secret` you need to authenticate your apps' users.



And then, from your Go program, authenticate your user like this :  

[embedmd]:# (./example/example.go)
```go
package main

import (
	oauth2ns "github.com/n-marshall/oauth2-noserver"
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
	/*client := ...*/ _ = oauth2ns.Authorize(conf)
	// use client.Get / client.Post for further requests, the token will automatically be there
}
```

The `AuthURL` and `TokenURL` can be found in the service's oauth documentation.
