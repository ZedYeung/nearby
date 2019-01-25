package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("service start")

	// with Auth
	// reference: https://auth0.com/blog/authentication-in-golang/
	// Create a new router on top of the existing http router as we need to check auth.
	r := mux.NewRouter()
	// Create a new JWT middleware with a Option that uses the key ‘mySigningKey’ such that we know this token is
	// from our server. The signing method is the default HS256 algorithm such that data is encrypted.
	// var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	// 		return signingKey, nil
	// 	},
	// 	SigningMethod: jwt.SigningMethodHS256,
	// })

	// It means we use jwt middleware to manage these endpoints and if they don’t have valid token, we will reject them.
	// r.Handle("/post", jwtMiddleware.Handler(http.HandlerFunc(Post))).Methods("POST")
	// r.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(Search))).Methods("GET")
	// login and signup don't need middleware, since it is the first time of request, we don't have token now.
	r.Handle("/login", http.HandlerFunc(Login)).Methods("POST")
	r.Handle("/signup", http.HandlerFunc(Signup)).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
