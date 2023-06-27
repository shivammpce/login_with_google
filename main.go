package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Replace the following client ID and client secret with your own
const (
	clientID     = ""
	clientSecret = ""
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
		<head>
			<style>
				body {
					background-color: #f2f2f2;
					font-family: Arial, sans-serif;
				}

				h1 {
					color: #333;
					text-align: center;
					margin-top: 100px;
				}

				.container {
					text-align: center;
					margin-top: 30px;
				}

				.btn-login {
					display: inline-block;
					padding: 10px 20px;
					background-color: #4285F4;
					color: #fff;
					text-decoration: none;
					border-radius: 4px;
					font-size: 16px;
					transition: background-color 0.3s ease;
				}

				.btn-login:hover {
					background-color: #3367D6;
				}
			</style>
		</head>
		<body>
			<h1>Welcome to the Home Page</h1>
			<div class="container">
				<p>Login with Google:</p>
				<a class="btn-login" href="/login">Login</a>
			</div>
		</body>
	</html>`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal(err)
	}

	client := oauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	// Parse the response and extract the name
	// Replace the JSON parsing with your own logic based on the response structure
	// In this example, we assume the response contains a 'name' field
	var userInfo struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, `<html>
		<body>
			<h1>Welcome, %s!</h1>
		</body>
	</html>`, userInfo.Name)
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
