package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var signals = SignalServer{cache: make(map[string]interface{})}
var answers = AnswerServer{cache: make(map[string][]interface{})}

type DataPayload struct {
	Hash string
	Data interface{}
}

type SignalServer struct {
	cache map[string]interface{}
}

type AnswerServer struct {
	cache map[string][]interface{}
}

func (s SignalServer) store(hash string, signal interface{}) {
	s.cache[hash] = signal
}

func (s SignalServer) getByHash(hash string) interface{} {
	return s.cache[hash]
}

func (s AnswerServer) push(hash string, answer interface{}) {
	if s.cache[hash] == nil {
		s.cache[hash] = make([]interface{}, 0)
	}

	s.cache[hash] = append(s.cache[hash], answer)
}

func (s AnswerServer) getByHash(hash string) []interface{} {
	return s.cache[hash]
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func GetSignal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	hash := ps.ByName("hash")

	data := signals.getByHash(hash)

	if data == nil {
		http.Error(w, "Invalid hash given", 400)
		return
	}

	json.NewEncoder(w).Encode(data)

	fmt.Println("Fetched signal for", hash)
}

func PostSignal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var payload DataPayload
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hash := payload.Hash
	data := payload.Data

	signals.store(hash, data)

	fmt.Println("Signal stored for", hash)
	fmt.Fprintf(w, "ok")
}

func GetAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hash := ps.ByName("hash")

	data := answers.getByHash(hash)

	if data == nil {
		http.Error(w, "Invalid hash given", 400)
		return
	}

	json.NewEncoder(w).Encode(data)

	fmt.Println("Fetched answer for", hash)
}

func PostAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var payload DataPayload
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hash := payload.Hash
	data := payload.Data

	answers.push(hash, data)

	fmt.Println("Answer stored for", hash)
	fmt.Fprintf(w, "ok")
}

func debugmsg(data interface{}) {
	fmt.Printf("%+v\n", data)
}

func main() {
	//setup signal and answer server

	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/signal/:hash", GetSignal)
	router.POST("/signal", PostSignal)

	router.GET("/answer/:hash", GetAnswer)
	router.POST("/answer", PostAnswer)

	fmt.Println("Signal server running")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":3000", handler))

}
