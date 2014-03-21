package main

// Use code.google.com/p/goauth2/oauth client to test
// Open url in browser:
// http://localhost:4000/app

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"net/http"
)

func main() {
	client := &oauth.Config{
		ClientId:     "36a4ad6bd2ae7de9",
		ClientSecret: "bc0f1d4094acc95b0229cb6916113b37",
		RedirectURL:  "http://localhost:4000/appauth/code",
		AuthURL:      "http://localhost:3000/authorize",
		TokenURL:     "http://localhost:3000/token",
	}
	ctransport := &oauth.Transport{Config: client}

	// Application home endpoint
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>"))
		//w.Write([]byte(fmt.Sprintf("<a href=\"/authorize?response_type=code&client_id=1234&state=xyz&scope=everything&redirect_uri=%s\">Login</a><br/>", url.QueryEscape("http://localhost:4000/appauth/code"))))
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Login</a><br/>", client.AuthCodeURL(""))))
		w.Write([]byte("</body></html>"))
	})

	// Application destination - CODE
	http.HandleFunc("/appauth/code", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		code := r.Form.Get("code")

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CODE<br/>"))

		if code != "" {

			var jr *oauth.Token
			var err error

			// if parse, download and parse json
			if r.Form.Get("doparse") == "1" {
				jr, err = ctransport.Exchange(code)
				if err != nil {
					jr = nil
					w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", err)))
				}
			}

			// show json access token
			if jr != nil {
				w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", jr.AccessToken)))
				if jr.RefreshToken != "" {
					w.Write([]byte(fmt.Sprintf("REFRESH TOKEN: %s<br/>\n", jr.RefreshToken)))
				}
			}

			w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

			cururl := *r.URL
			curq := cururl.Query()
			curq.Add("doparse", "1")
			cururl.RawQuery = curq.Encode()
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))
		} else {
			w.Write([]byte("Nothing to do"))
		}

		w.Write([]byte("</body></html>"))
	})

	http.ListenAndServe(":4000", nil)
}
