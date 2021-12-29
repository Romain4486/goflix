package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
)

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Goflix")
	}
}

func (s *server) handleTokenCreate() http.HandlerFunc {
	type request struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token string `json:"token"`
	}
	type responseError struct {
		Error string `json:"error"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse the login body. Err=%v", err)
			log.Panicln(msg)
			s.respond(w, r, responseError{
				Error: msg,
			}, http.StatusBadRequest)
			return
		}

		if req.UserName != "Golang" || req.Password != "rocks" {
			s.respond(w, r, responseError{
				Error: "Invalid credentials",
			}, http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.UserName,
			"exp":      time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat":      time.Now().Unix(),
		})

		ts, err := token.SignedString([]byte(JWT_APP_KEY))
		if err != nil {
			msg := fmt.Sprintf("Cannot generate JWT err=%v", err)
			s.respond(w, r, msg, http.StatusInternalServerError)
		}

		s.respond(w, r, response{
			Token: ts,
		}, http.StatusOK)

	}
}
