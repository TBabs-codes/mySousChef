package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TBabs-codes/mySousChef/auth"
	"github.com/TBabs-codes/mySousChef/internal/database"
	"github.com/google/uuid"
)

//Creates user
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type userCreationReq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userInfo := userCreationReq{}
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("Error decoding user creation request: %s", err)
		w.WriteHeader(500)
		return
	}

	hashed_password, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		log.Printf("Error creating password: %s", err)
		w.WriteHeader(500)
		return
	}

	// err = auth.CheckPasswordHash(userInfo.Password, hashed_password)
	// if err != nil {
	// 	log.Printf("Error creating password2: %s", err)
	// 	w.WriteHeader(500)
	// 	return
	// }

	type UserCreationResp struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Premium    bool      `json:"premium"`
	}

	//SQL query to create user if no error els
	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:             uuid.New().String(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          userInfo.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Printf("Error adding user to database: %s", err)
		return
	}


	//Create JWT
	jwt, err := auth.MakeJWT(uuid.MustParse(user.ID), cfg.JWTSecret, cfg.JWT_Timeout)
	if err != nil {
		log.Printf("Error creating JSON web token: %s", err)
		return
	}

	//Add JWT to DB
	refresh_token, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: jwt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),

	})
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		return
	}

	//Update response header to include JWT
	w.Header().Set("Authorization", "Bearer " + refresh_token.Token)


	respBody := UserCreationResp{
		ID:         uuid.MustParse(user.ID),
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
		Premium:    user.Premium,
	}
	
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)

}

//Logins User
func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type userLoginReq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userInfo := userLoginReq{}
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("Error decoding user login request: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.DB.ReturnUserFromEmail(r.Context(), userInfo.Email)
	if err != nil {
		log.Printf("Error finding user with email provided: %s", err)
		w.WriteHeader(500)
		return
	}

	err = auth.CheckPasswordHash(userInfo.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Password unhashing failed: %s", err)
		w.WriteHeader(500)
		return
	}

	//Create JWT
	jwt, err := auth.MakeJWT(uuid.MustParse(user.ID), cfg.JWTSecret, cfg.JWT_Timeout)
	if err != nil {
		log.Printf("Error creating JSON web token: %s", err)
		return
	}

	//Add JWT to DB
	refresh_token, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: jwt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),

	})
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		return
	}

	//Update response header to include JWT
	w.Header().Set("Authorization", "Bearer " + refresh_token.Token)


	//Preparing Response Body
	type UserLoginResp struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Premium    bool      `json:"premium"`
	}

	respBody := UserLoginResp{
		ID:         uuid.MustParse(user.ID),
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
		Premium:    user.Premium,
	}
	
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)

}
