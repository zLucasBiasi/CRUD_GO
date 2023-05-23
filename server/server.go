package server

import (
	"crud/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// create user in database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}
	var user user

	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		w.Write([]byte("erro ao converter o usuario para struct"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao converter e conectar no banco de dados"))
		return
	}

	defer db.Close()
	//prepare statement

	statement, err := db.Prepare(("insert into users (nome, email) values (?,?)"))
	if err != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer statement.Close()

	insertion, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("Erro ao executar o statement"))
		return
	}

	idInsert, err := insertion.LastInsertId()
	if err != nil {
		w.Write([]byte("Erro ao obter o id inserido"))
		return
	}

	//status codes
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso! Id: %d", idInsert)))
}

// search all users in table
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Falha ao conectar com o banco"))
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		w.Write([]byte("erro ao buscar os usuarios"))
		return
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var user user

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("erro ao buscar o usuario"))
			return
		}
		users = append(users, user)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("erro ao converter os usuarios"))
		return
	}
}

// search a specific user in table
func SearchUser(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao converter o parametro para inteiro"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar ao banco"))
		return
	}

	defer db.Close()
	row, err := db.Query("select * from users where id = ?", ID)
	if err != nil {
		w.Write([]byte("Erro ao buscar o usuario"))
		return
	}

	var user user

	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao escanear o usuario"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("erro ao converter o usuario para json"))
		return
	}
}

// edit a specific user in table
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao converter o parametro para inteiro"))
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Erro ao buscar o corpo da request"))
		return
	}

	var user user
	if err := json.Unmarshal(bodyRequest, &user); err != nil {
		w.Write([]byte("Erro ao converter o usuario para struct"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar ao banco"))
		return
	}

	defer db.Close()
	//sempre usar o statement para coisas que nao seja consulta, post, delete, put

	statement, err := db.Prepare("update users set nome =?, email = ? where id = ?")

	if err != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer statement.Close()

	if _, err := statement.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	w.Write([]byte("usuario atualizado"))

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao converter o parametro para inteiro"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar ao banco"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from users where id = ?")
	if err != nil {
		w.Write([]byte("Erro ao criar o stateament"))
		return
	}

	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar o usuario"))
		return
	}

	w.Write([]byte("usuario deletado"))

	w.WriteHeader(http.StatusNoContent)
}
