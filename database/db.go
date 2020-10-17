package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	
)
const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = ""
	dbname = ""
)
func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	stmt := `INSERT INTO users (age, email, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err = db.Exec(stmt, 24, "jhoydich3@gmail.com", "Jeremiah", "Hoydich")
	if err != nil {
		panic(err)
	}
}
