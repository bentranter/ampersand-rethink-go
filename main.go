package main

import (
	"encoding/json"
	"log"
	"net/http"

	rdb "github.com/dancannon/gorethink"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Person is the model for our user
type Person struct {
	ID             string `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	CoolnessFactor int    `json:"coolnessFactor"`
}

// DB holds our connection to the database
type DB struct {
	Session *rdb.Session
}

func newDBConn() *DB {
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatal("Error: %s\n", err)
	}
	return &DB{
		Session: session,
	}
}

func main() {
	db := newDBConn()
	db.Init("arg", "People")
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
	})
	router := httprouter.New()

	router.GET("/api/people", db.List)
	router.GET("/api/people/:id", db.Get)
	router.PUT("/api/people/:id", db.Update)
	router.DELETE("/api/people/:id", db.Delete)
	router.POST("/api/people", db.Add)

	http.ListenAndServe(":8000", c.Handler(router))
}

// Init creates a new DB and a new table
func (db *DB) Init(dbName string, tableName string) error {
	resp, err := rdb.DBCreate(dbName).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Println("DB created: ", resp.DBsCreated)

	resp, err = rdb.TableCreate(tableName).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Println("Table created: ", resp.TablesCreated)

	return nil
}

// List all users
func (db *DB) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var people []Person
	rows, err := rdb.Table("People").Run(db.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = rows.All(&people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(people)
}

// Get all users
func (db *DB) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p Person
	row, err := rdb.Table("People").Get(ps.ByName("id")).Run(db.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	row.One(&p)
	json.NewEncoder(w).Encode(p)
}

// Update a specific user
func (db *DB) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p Person
	json.NewDecoder(r.Body).Decode(&p)

	resp, err := rdb.Table("People").Get(ps.ByName("id")).Update(p).RunWrite(db.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp.GeneratedKeys)
}

// Delete a specific user
func (db *DB) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp, err := rdb.Table("People").Get(ps.ByName("id")).Delete().RunWrite(db.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp)
}

// Add a new user
func (db *DB) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p Person
	json.NewDecoder(r.Body).Decode(&p)

	row, err := rdb.Table("People").Insert(p).Run(db.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = row.One(p)
	json.NewEncoder(w).Encode(p)
}
