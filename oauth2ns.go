package oauth2ns

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

const (
	PORT = 14565
)

var (
	ctx context.Context
)

type Result struct {
	Client *http.Client
	Token  *oauth2.Token
}

func Authorize(conf *oauth2.Config) *Result {
	// add transport for self-signed certificate to context
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(context.Background(), oauth2.HTTPClient, sslcli)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	conf.RedirectURL = fmt.Sprintf("http://127.0.0.1:%s/oauth/callback", strconv.Itoa(PORT))
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	quitSignal := make(chan *Result)
	srv := startHttpServer(conf, quitSignal)
	log.Println(color.CyanString("You will now be taken to your browser for authentication"))
	time.Sleep(600 * time.Millisecond)
	open.Run(url)
	time.Sleep(600 * time.Millisecond)
	log.Printf("Authentication URL: %s\n", url)

	log.Printf("main: serving for 10 seconds")
	log.Printf("main: stopping HTTP server")

	// now close the server gracefully ("shutdown")
	// timeout could be given instead of nil as a https://golang.org/pkg/context/
	r := <-quitSignal
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	return r
}

func startHttpServer(conf *oauth2.Config, quitSignal chan *Result) *http.Server {
	http.HandleFunc("/oauth/callback", callbackHandler(conf, quitSignal))
	srv := &http.Server{Addr: ":" + strconv.Itoa(PORT)}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Printf("Httpserver error: %s", err)
			}
		}
	}()

	return srv
}

func callbackHandler(conf *oauth2.Config, quitSignal chan *Result) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParts, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("<p>%s</p>"), err.Error())
		}
		code := queryParts["code"][0]

		// Exchange will do the handshake to retrieve the initial access token.
		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		// The HTTP Client returned by conf.Client will refresh the token as necessary.
		client := conf.Client(ctx, tok)

		// show success page
		successPage := `
		<p style="height:100px; width:100%; display:flex; flex-direction: column; justify-content: center; text-align:center; background-color:#2ecc71; color:white; font-size:22">Success!</p>
		<p style="margin-top:20px; font-size:16">You are authenticated, you can now return to the CLI. This will auto-close</p>
		<script>window.onload=function(){setTimeout(this.close, 2000)}</script>
		`
		fmt.Fprintf(w, successPage)
		quitSignal <- &Result{client, tok}
	}
}
