package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v3"
)

func main() {
	fmt.Println("service start")

	// Initialize Elasticsearch
	ctx := context.Background()
	initES(ctx)
	// with Auth
	// reference: https://auth0.com/blog/authentication-in-golang/
	// Create a new router on top of the existing http router as we need to check auth.
	r := mux.NewRouter()
	// Create a new JWT middleware with a Option that uses the key ‘mySigningKey’ such that we know this token is
	// from our server. The signing method is the default HS256 algorithm such that data is encrypted.
	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SIGN_KEY")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	// It means we use jwt middleware to manage these endpoints and if they don’t have valid token, we will reject them.
	r.Handle("/post/", jwtMiddleware.Handler(http.HandlerFunc(Post))).Methods("POST", "OPTIONS")
	r.Handle("/search/", jwtMiddleware.Handler(http.HandlerFunc(Search))).Methods("GET", "OPTIONS")
	// login and signup don't need middleware, since it is the first time of request, we don't have token now.
	r.Handle("/login/", http.HandlerFunc(Login)).Methods("POST", "OPTIONS")
	r.Handle("/signup/", http.HandlerFunc(Signup)).Methods("POST", "OPTIONS")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")), nil))
}

func initES(ctx context.Context) {
	// Create a client, which means we create a connection to ES. If there is err, return.
	client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_URL")), elastic.SetSniff(false))
	if err != nil {
		panic(err)
		return
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(os.Getenv("ES_INDEX")).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		// Create a new index.
		// If not, create a new mapping. For other fields (user, message, etc.)
		// no need to have mapping as they are default. For geo location (lat, lon),
		// we need to tell ES that they are geo points instead of two float points
		// such that ES will use Geo-indexing for them (K-D tree)

		mapping := `{
			"mappings":{
				"post":{
					"properties":{
						"location":{
							"type":"geo_point"
						}
					}
				}
			}
		}
		`
		_, err := client.CreateIndex(os.Getenv("ES_INDEX")).Body(mapping).Do(ctx) // Create this index
		if err != nil {
			// Handle error
			panic(err)
		}
	}
}
