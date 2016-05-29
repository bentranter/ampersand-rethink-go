package main

import (
	"encoding/json"
	"log"
	"net/http"

	rdb "github.com/dancannon/gorethink"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var (
	session   = newDBConn()
	tableName = "People"
	c         = newCORSOptions()
)

// Person is the model for our user.
type Person struct {
	ID             string `json:"id" gorethink:"id,omitempty"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	CoolnessFactor int    `json:"coolnessFactor"`
}

func newDBConn() *rdb.Session {
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fataln("Error: %s", err)
	}

	resp, err := rdb.TableCreate(tableName).RunWrite(session)
	if err != nil {
		log.Printf("Note: %s\n", err)
	}
	log.Println("Tables created: ", resp.TablesCreated)

	return session
}

func newCORSOptions() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
	})
	return c
}

// List all users in the DB.
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var people []Person
	rows, err := rdb.Table(tableName).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = rows.All(&people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(people)
}

// Get all users in the DB.
func Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p Person
	row, err := rdb.Table(tableName).Get(ps.ByName("id")).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	row.One(&p)
	json.NewEncoder(w).Encode(p)
}

// Update a specific user
func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p Person
	json.NewDecoder(r.Body).Decode(&p)

	resp, err := rdb.Table(tableName).Get(ps.ByName("id")).Update(p).RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp.GeneratedKeys)
}

// Delete a specific user from the DB>
func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp, err := rdb.Table(tableName).Get(ps.ByName("id")).Delete().RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp)
}

// Add a new user
func Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p Person
	json.NewDecoder(r.Body).Decode(&p)

	row, err := rdb.Table(tableName).Insert(p).RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	id := row.GeneratedKeys
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": id[0],
	})
}

func main() {
	router := httprouter.New()

	router.GET("/api/people", List)
	router.GET("/api/people/:id", Get)

	router.PUT("/api/people/:id", Update)
	router.DELETE("/api/people/:id", Delete)
	router.POST("/api/people", Add)

	http.ListenAndServe(":8000", c.Handler(router))
}
