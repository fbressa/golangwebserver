package main

/* Este projeto tem o intuito de praticar a criação de um web server em golang com abertura para aprimoramento */
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// definindo uma "table" no nosso "banco"
type User struct {
	Name string `json:"name"`
}

// Definindo um "banco"
var userCache = make(map[int]User)

// Mutex para fazer o request thread safe
var cacheMutex sync.RWMutex

func main() {
	//Definimos um router
	mux := http.NewServeMux()

	//Definimos um handler pra uma route root
	mux.HandleFunc("/", handleRoot)

	//Criamos um novo handler para uma route de criar usuário
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc(
		"GET /users/{id}",
		getUser,
	)
	mux.HandleFunc(
		"DELETE /users/{id}",
		deleteUser,
	)

	//Iniciamos o servidor
	fmt.Println("Server runing on :8080")
	http.ListenAndServe(":8080", mux)
}

// A route root aqui
func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "hello world")
}

// A route de criar usuário aqui
func createUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if user.Name == "" {
		http.Error(
			w,
			"Name is required",
			http.StatusBadRequest,
		)
		return
	}

	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)

}

func getUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(
			w,
			"user not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func deleteUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if _, ok := userCache[id]; !ok {
		http.Error(
			w,
			"user not found",
			http.StatusBadRequest,
		)
		return
	}

	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
