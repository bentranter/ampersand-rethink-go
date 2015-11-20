package main

import (
	"encoding/json"
	"log"
	"net/http"

	rdb "github.com/dancannon/gorethink"
	"github.com/julienschmidt/httprouter"
)

// Conn holds our connection to the database
type Conn struct {
	Session *rdb.Session
}

func newConn() Conn {
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatal("Error: %s\n", err)
	}
	return Conn{
		Session: session,
	}
}

func main() {
	conn := newConn()
	conn.Init("arg", "People")
	router := httprouter.New()

	router.GET("/api/people", conn.List)
	router.GET("/api/people/:id", conn.Get)
	router.PUT("/api/people/:id", conn.Update)
	router.DELETE("/api/people/:id", conn.Delete)
	router.POST("/api/people", conn.Add)

	http.ListenAndServe(":3000", router)
}

// Init creates a new DB and a new table
func (c *Conn) Init(db string, table string) error {
	resp, err := rdb.DBCreate(db).RunWrite(c.Session)
	if err != nil {
		return err
	}
	log.Println("DB created: ", resp.DBsCreated)

	resp, err = rdb.TableCreate(table).RunWrite(c.Session)
	if err != nil {
		return err
	}
	log.Println("Table created: ", resp.TablesCreated)

	return nil
}

// List all users
func (c *Conn) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp, err := rdb.Table("People").Run(c.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Will this work?
	json.NewEncoder(w).Encode(resp)
}

// Get all users
func (c *Conn) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps.ByName("id")
	resp, err := rdb.Table("People").Get(user).RunWrite(c.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp)
}

// Update a specific user
func (c *Conn) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

// Delete a specific user
func (c *Conn) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

// Add a new user
func (c *Conn) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp, err := rdb.Table("People").Insert(r.Body).Run(c.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp)
}
