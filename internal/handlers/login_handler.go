package handlers

import (
	"encoding/json"
	"net/http"
	"the-gym-app/internal/middleware"
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
	}
}
