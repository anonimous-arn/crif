package main

import (
	"cibt/cibt"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type InputData struct {
	Inn string `json:"inn"`
}

var ErrorNotFound = map[string]string{
	"message": "client not found",
}

var ErrInternal = map[string]string{
	"message": "something went wrong or cibt wasn't answered",
}

var ErrBadRequest = map[string]string{
	"message": "bad request",
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ping", pong).Methods("GET")
	router.HandleFunc("/getCibtData", getCibtData).Methods("POST")
	http.Handle("/", router)

	fmt.Println("server starts")
	http.ListenAndServe(":8181", nil)
}

func pong(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Pong"))
}

func getCibtData(rw http.ResponseWriter, r *http.Request) {

	taxPayerId := r.URL.Query().Get("inn")
	if taxPayerId == "" {
		setHeader(rw)
		rw.WriteHeader(400)
		rw.Write([]byte("Inn param is empty"))
		rw.Header().Set("Content-Type", "application/xml")
		return
	}

	cibt_id := cibt.GetCibtId(taxPayerId)
	if cibt_id == "" {
		item, _ := json.Marshal(ErrorNotFound)
		setHeader(rw)
		rw.WriteHeader(404)
		rw.Write(item)
		rw.Header().Set("Content-Type", "application/xml")
		return
	}

	responseBody := cibt.GetCibtInfo(cibt_id)
	responseBody = strings.ReplaceAll(responseBody, "&lt;", "<")
	responseBody = strings.ReplaceAll(responseBody, "&gt;", ">")

	fmt.Println("bodyString -", responseBody)
	setHeader(rw)
	rw.Header().Set("Content-Type", "application/xml")
	rw.Write([]byte(responseBody))
}

func setHeader(rw http.ResponseWriter) http.ResponseWriter {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
	return rw
}
