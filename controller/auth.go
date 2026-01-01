package controller

import (
	"cometosee/common"
	"cometosee/intailizer"
	"cometosee/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Api struct {
	Addr string
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body model.Auth

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid request", http.StatusBadRequest)
		return
	}

	//hash
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		common.WriteJSONError(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	query := `insert into cometoseeauth (username,email,password) values ($1,$2,$3)`

	_, err = intailizer.DB.Exec(query, body.Username, strings.ToLower(body.Email), string(hash))

	if err != nil {
		common.WriteJSONError(w, "error in inserting value in database or email already exists", http.StatusInternalServerError)
		return
	}

	common.WriteJSONMessage(w, "user created sucessfully")

}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body model.Auth

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(body.Email)

	var hashpassword, username string

	err := intailizer.DB.QueryRow(`select username, password from cometoseeauth where email=$1`, email).Scan(&username, &hashpassword)
	if err != nil {
		if err == sql.ErrNoRows {
			common.WriteJSONError(w, "Email not found", http.StatusBadRequest)
		} else {
			common.WriteJSONError(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(body.Password)); err != nil {
		common.WriteJSONError(w, "invalid password", http.StatusBadRequest)
		return
	}

	token, err := geneareJwt(email, username)
	if err != nil {
		common.WriteJSONError(w, "error in generating jwt token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "login sucessfully",
		"token":    token,
		"username": username,
	})

}

func geneareJwt(email string, username string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
