package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<title>Keep Up</title>
			</head>
			<body>
			<h1>Cookie Example</h1>
			</body>
			</html>
			`))
	})

	r.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   "not_encoded_cookie",
			Value:  "this is the value",
			Path:   "/",
			MaxAge: 10000,
		}

		http.SetCookie(w, cookie)
		w.Write([]byte(`<h1>Cookie Example</h1>`))
	})

	r.HandleFunc("/encoded_cookie", func(w http.ResponseWriter, r *http.Request) {
		secret := "987ef5cec47690b20e08dc77ae792079312e833221922c876a716c3e3472fc69"
		hashKey := []byte(secret)
		blockKey := []byte(nil)

		sc := securecookie.New(hashKey, blockKey)

		value := map[string]string{
			"value": "1234",
		}
		if encoded, err := sc.Encode("encoded_cookie", value); err == nil {
			encodedCookie := &http.Cookie{
				Name:   "encoded_cookie",
				Value:  encoded,
				MaxAge: 10000,
			}
			if err != nil {
				http.Error(w, "error making cookie", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, encodedCookie)
		}
	})

	r.HandleFunc("/decoded_cookie", func(w http.ResponseWriter, r *http.Request) {
		secret := "987ef5cec47690b20e08dc77ae792079312e833221922c876a716c3e3472fc69"

		hashKey := []byte(secret)
		blockKey := []byte(nil)

		sc := securecookie.New(hashKey, blockKey)

		value := make(map[string]string)

		if decodedCookie, err := r.Cookie("encoded_cookie"); err == nil {
			err = sc.Decode("encoded_cookie", decodedCookie.Value, &value)
			fmt.Fprintf(w, "This is the value: %s", value["value"])
			if err != nil {
				fmt.Println(err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		}
	})

	log.Infoln("Listening at http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}
