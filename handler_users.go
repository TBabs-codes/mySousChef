package main

import (
	"encoding/json"
	"net/http"

	"github.com/TBabs-codes/recipe_book/auth"
)

// Creates user
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
		ID            uuid.UUID `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Is_Chirpy_Red bool      `json:"is_chirpy_red"`
	}

	//SQL query to create user if no error els
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          userInfo.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Printf("Error adding user to database: %s", err)
		return
	}

	respBody := UserCreationResp{
		ID:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		Is_Chirpy_Red: user.IsChirpyRed,
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