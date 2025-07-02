package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"the-gym-app/internal/middleware"
	"the-gym-app/internal/models"
	"the-gym-app/internal/services"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

func NewLoginHandler(db *services.DatabaseService) *LoginHandler {
	return &LoginHandler{db}
}

type LoginHandler struct {
	db *services.DatabaseService
}

func (l *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//decode the parameters in the request to a LoginUser object and then make db calls to validate
	//if the password for this user matches
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Request is not valid", http.StatusBadRequest)
		return
	}
	//validate the user credentials
	exists, err := l.db.CheckIfUserExists(credentials.Username)
	if err != nil {
		http.Error(w, "Server is facing problems at the moment", http.StatusInternalServerError)
		return
	}
	if exists {
		//verify password
		verified, err := l.db.CheckIfUserPasswordCorrect(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "Server is facing problems at the moment", http.StatusInternalServerError)
			return
		}
		if verified {
			//generate jwt token for the user
			jwtToken, err := middleware.GenerateToken(credentials.Username)
			if err != nil {
				http.Error(w, "Server is facing problems at the moment", http.StatusInternalServerError)
				return
			}
			response := map[string]string{
				"token": jwtToken,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	} else {
		http.Error(w, "User does not exist. Please sign up", http.StatusUnauthorized)
	}
}

func (l *LoginHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	var newUser Credentials
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Problem creating new user, please check in again later", http.StatusInternalServerError)
		return
	}
	var newUserToDb models.User
	newUserToDb.Username = newUser.Username
	newUserToDb.PasswordHash = newUser.Password
	newUserToDb.Email = newUser.Email
	if err := l.db.SaveUser(&newUserToDb); err != nil {
		http.Error(w, "Problem saving new user, please check in again later", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Welcome to The Gym App! Please login using your credentials")
}
