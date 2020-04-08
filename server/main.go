package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"hermes"
	"net/http"
)

var (
	router *mux.Router
)

func main() {
	var (
		listen string
		// cert    string
		// key     string
		// keypass string
		dbname string
		dbhost string
		dbuser string
		dbpass string
	)
	flag.StringVar(&listen, "listen", "localhost:5587", "Address to listen to for Hermes server")
	flag.StringVar(&dbname, "dbname", "hermes", "Database name where to store user information")
	flag.StringVar(&dbuser, "dbuser", "hermes", "Database user for database")
	flag.StringVar(&dbpass, "dbpass", "hermes", "Database password for database (pass it via env var)")
	flag.StringVar(&dbhost, "dbhost", "localhost:5432", "Database host and port")
	flag.Parse()

	hermes.InitStorage(dbname, dbhost, dbuser, dbpass)
	hermes.InitQueue()
	router = mux.NewRouter()
	router.HandleFunc("/send", send).Name("send")
	router.HandleFunc("/status/{uuid}", status).Name("status")
	server := http.Server{
		Addr:    listen,
		Handler: router,
	}
	server.ListenAndServe()
	hermes.CloseQueue()
	hermes.CloseStorage()
}

func send(w http.ResponseWriter, r *http.Request) {
	customer, err := hermes.GetCustomer(r.Header.Get("Hermes-Token"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	message, err := hermes.ParseMessage(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	handle, err := hermes.Send(customer, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	statusURL, _ := router.Get("status").URL("uuid", handle)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", statusURL.String())
	w.Write([]byte(fmt.Sprintf("{uuid: \"%s\"}", handle)))
}

func status(w http.ResponseWriter, r *http.Request) {
	// var (
	// 	vars = mux.Vars(r)
	// )
	// customer, err := hermes.GetCustomer(r.Header.Get("Hermes-Token"))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// status, err := GetStatus(customer, vars["uuid"])
	w.Write([]byte("OK"))
}
