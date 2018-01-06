package oauth2ns

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	rndm "github.com/nmrshll/rndm-go"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

type contextKey int

const (
	// PORT is the port that the temporary oauth server will listen on
	PORT                                  = 14565
	oauthStateStringContextKey contextKey = iota
)

// Authorize starts the login process
func Authorize(conf *oauth2.Config) *http.Client {
	// add transport for self-signed certificate to context
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, sslcli)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	conf.RedirectURL = fmt.Sprintf("http://127.0.0.1:%s/oauth/callback", strconv.Itoa(PORT))

	// Some random string, random for each request
	oauthStateString := rndm.String(8)
	ctx = context.WithValue(ctx, oauthStateStringContextKey, oauthStateString)
	url := conf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

	quitSignal := make(chan *http.Client)
	srv := startHTTPServer(ctx, conf, quitSignal)
	log.Println(color.CyanString("You will now be taken to your browser for authentication"))
	time.Sleep(600 * time.Millisecond)
	open.Run(url)
	time.Sleep(600 * time.Millisecond)
	// log.Printf("Authentication URL: %s\n", url)

	// When the callbackHandler returns a client, it's time to shutdown the server gracefully
	// timeout could be given instead of nil as a https://golang.org/pkg/context/
	client := <-quitSignal
	log.Printf("stopping HTTP server")
	shutdownContext, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(shutdownContext); err != nil {
		log.Fatal(err) // failure/timeout shutting down the server gracefully
	}
	fmt.Println("Server gracefully stopped")
	return client
}

func startHTTPServer(ctx context.Context, conf *oauth2.Config, quitSignal chan *http.Client) *http.Server {
	http.HandleFunc("/oauth/callback", callbackHandler(ctx, conf, quitSignal))
	srv := &http.Server{Addr: ":" + strconv.Itoa(PORT)}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver error: %s", err)
		}
	}()

	return srv
}

func callbackHandler(ctx context.Context, conf *oauth2.Config, quitSignal chan *http.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestStateString := ctx.Value(oauthStateStringContextKey).(string)
		responseStateString := r.FormValue("state")
		if responseStateString != requestStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", requestStateString, responseStateString)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := conf.Exchange(ctx, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// The HTTP Client returned by conf.Client will refresh the token as necessary
		client := conf.Client(ctx, token)

		// show success page
		successPage := `
		<div style="height:100px; width:100%!; display:flex; flex-direction: column; justify-content: center; align-items:center; background-color:#2ecc71; color:white; font-size:22"><div>Success!</div></div>
		<p style="margin-top:20px; font-size:18; text-align:center">You are authenticated, you can now return to the program. This will auto-close</p>
		<script>window.onload=function(){setTimeout(this.close, 4000)}</script>
		`
		fmt.Fprintf(w, successPage)
		quitSignal <- client
	}
}
