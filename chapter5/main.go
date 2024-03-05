package main

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := Put(key, string(value)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	logger.WritePut(key, string(value))
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if err := Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	logger.WriteDelete(key)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := Get(key)
	if err != nil {
		if errors.Is(err, ErrorNoSuchKey) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.Write([]byte(value))
}

func main() {
	if err := initializeTransactionLog(); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods(http.MethodPut)
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods(http.MethodDelete)
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
