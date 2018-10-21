package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	. "github.com/travierm/PeerChat/server/lib"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func GetSignal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	hash := ps.ByName("hash")

	data := signals.GetByHash(hash)

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

	signals.Store(hash, data)

	fmt.Println("Signal stored for", hash)
	fmt.Fprintf(w, "ok")
}

func GetAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hash := ps.ByName("hash")

	data := answers.GetByHash(hash)

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

	answers.Push(hash, data)

	fmt.Println("Answer stored for", hash)
	fmt.Fprintf(w, "ok")
}

func debugmsg(data interface{}) {
	fmt.Printf("%+v\n", data)
}

//setup signal and answer servers
var signals = SignalServer{Cache: make(map[string]interface{})}
var answers = AnswerServer{Cache: make(map[string][]interface{})}

func main() {

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
