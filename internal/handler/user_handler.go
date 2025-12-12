package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"go-api/cmd/server/auth"
	"go-api/internal/database"
	"go-api/internal/models"
)

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(
		"INSERT INTO users (name, email, password) VALUES ($1,$2,$3)",
		user.Name, user.Email, user.Password,
	)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	result, err := database.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "User not found", 404)
		return
	}

	w.Write([]byte("User deleted"))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		"UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		user.Name, user.Email, user.Password, id,
	)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("User updated"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&loginData)

	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	var user models.User
	var storedPassword string

	err := database.DB.QueryRow(
		"SELECT id, name, email, password FROM users WHERE email = $1",
		loginData.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &storedPassword)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Simple password check (in production, use bcrypt)
	if storedPassword != loginData.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
