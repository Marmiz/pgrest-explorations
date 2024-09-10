package main

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func (app *application) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.logger.Error("failed to read JSON", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
		return
	}

	// perform some basic validation
	v := validateEmail(input.Email)
	if !v {
		app.logger.Error("invalid email address", "email", input.Email)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	v = validatePasswordText(input.Password)
	if !v {
		app.logger.Error("invalid password", "attempted", input.Password)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	dbUser, err := app.queries.GetUser(r.Context(), input.Email)
	if err != nil {
		app.logger.Info("failed to get user", "user", input.Email, "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
		return
	}

	match, err := comparePasswords(dbUser.PasswordHash, []byte(input.Password))
	if err != nil {
		app.logger.Error("failed to compare passwords", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
		return
	}

	if !match {
		app.logger.Info("invalid credentials", "email", input.Email)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// generate a new jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": dbUser.Role,
		"exp":  jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"iat":  jwt.NewNumericDate(time.Now()),
		"nbf":  jwt.NewNumericDate(time.Now()),
		"aud":  "https://example.com",
		"iss":  "https://example.com",
	})

	tokenString, err := token.SignedString([]byte(app.config.jwt.secret))
	if err != nil {
		app.logger.Error("failed to sign token", "error", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string]string{"token": tokenString}, nil)
	if err != nil {
		app.logger.Error("failed to write JSON", err)
		http.Error(w, "Server encountered a problem", http.StatusInternalServerError)
	}
}

func validateEmail(email string) bool {
	if email == "" {
		return false
	}

	return EmailRX.MatchString(email)
}

func validatePasswordText(password string) bool {
	if password == "" {
		return false
	}

	return len(password) >= 8 && len(password) <= 72
}

func comparePasswords(p1 []byte, p2 []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p1, p2)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
