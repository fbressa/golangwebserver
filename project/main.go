package main

/* Este projeto tem o intuito de praticar a criação de um web server em golang com abertura para aprimoramento */
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// definindo uma "table" no nosso "banco"
type User struct {
	Name string `json:"name"`
}

// o objeto sql.DB já é thread-safe por natureza
var db *sql.DB

func main() {
	var err error

	// String de conexão: "usuario:senha@tcp(host:porta)/nomedobanco"
	// Lembre que sua porta é a 3307!
	connStr := "root:123456@tcp(127.0.0.1:3307)/users_db"

	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Erro checando driver de conexão: ", err)
	}
	// O Open() não conecta na hora, ele só "abre" as configurações.
	// O Ping() de fato força a comunicação real com o Docker!
	err = db.Ping()
	if err != nil {
		log.Fatal("Não foi possível alcançar o MySQL: ", err)
	}

	fmt.Println("Conexão com MySQL estabelecida com sucesso!")

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

	// colocamos o ? para previnir SQL Injection
	query := "INSERT INTO users (name) VALUES (?)"

	result, err := db.Exec(query, user.Name)

	if err != nil {
		http.Error(
			w,
			"Erro ao salvar no banco"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	//Vamos pegar o ID gerado pelo auto increment
	id, _ := result.LastInsertId()
	fmt.Printf("Usuário adicionado no banco com o id: %d\n", id)

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

	var user User
	query := "SELECT name FROM users WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&user.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(
				w,
				"user not found",
				http.StatusNotFound,
			)
			return
		}
		http.Error(
			w,
			"Erro ao buscar do banco: "+err.Error(),
			http.StatusInternalServerError,
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

	// colocamos o ? para previnir SQL Injection
	query := "DELETE FROM users WHERE id = ?"

	result, err := db.Exec(query, id)

	if err != nil {
		http.Error(
			w,
			"Erro ao salvar do banco"+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// O método RowsAffected() nos diz quantas linhas o banco apagou de fato
	linhasDeletadas, _ := result.RowsAffected()

	if linhasDeletadas == 0 {
		http.Error(w, "Usuário não encontrado no banco", http.StatusNotFound)
		return
	}

	fmt.Printf("Usuário ID %d deletado com sucesso.\n", id)

	w.WriteHeader(http.StatusNoContent)
}
