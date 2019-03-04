package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"github.com/ZedYeung/nearby/model"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func Signup(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Received one signup request")

	// Connect to the postgres db
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))

	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	// Parse and decode the request body into a new `User` instance
	user := &model.User{}
	err = json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var dbUsername string
	err = db.QueryRow("SELECT username FROM users WHERE username=$1", user.Username).Scan(&dbUsername)

	switch {
	case err == sql.ErrNoRows:
		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account. Cannot get hashedPassword.", 500)
			return
		}

		// Next, insert the username, along with the hashed password into the database
		_, err = db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", user.Username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account. Cannot insert into database.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		fmt.Println(err)
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		res.Write([]byte("User already exists"))
		http.Redirect(res, req, "/login/", 301)
	}
}

// If login is successful, a new token is created.
func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Received one login request")

	// Connect to the postgres db
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))

	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	// Parse and decode the request body into a new `User` instance
	user := &model.User{}
	err = json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var dbPassword string

	err = db.QueryRow("SELECT password FROM users WHERE username=$1", user.Username).Scan(&dbPassword)

	switch {
	case err == sql.ErrNoRows:
		http.Error(res, "User not exists", http.StatusUnauthorized)
		return
	case err == nil:
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password))
		if err != nil {
			http.Redirect(res, req, "/login", 301)
			return
		}

		// Create a new token object to store.
		token := jwt.New(jwt.SigningMethodHS256)
		// Convert it into a map for lookup
		claims := token.Claims.(jwt.MapClaims)
		/*
			Set token claims
			Store username and expiration into it.
		*/
		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		/* Sign (Encrypt) the token with our secret such that only server knows it. */
		tokenString, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))
		if err != nil {
			fmt.Println(err)
			http.Error(res, "cannot sign token", http.StatusInternalServerError)
		}

		/* Finally, write the token to the browser window */
		res.Write([]byte(tokenString))
		fmt.Println(tokenString)
		fmt.Println([]byte(tokenString))

		// res.Write([]byte("Welcome back " + user.Username))
	default:
		http.Redirect(res, req, "/login/", 301)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
