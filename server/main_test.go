package main

import (
	"math/rand"
	"testing"
	"time"
)

type testSignal struct {
	Name string
}

type testAnswer struct {
	Name string
}

func init() {
	seed := time.Now().Unix()
	rand.Seed(seed)
}

func TestSignalServer(t *testing.T) {
	server := SignalServer{cache: make(map[string]interface{})}
	server.store("ABCD", testSignal{Name: "room1"})
	server.store("CHDD", testSignal{Name: "room2"})

	result := server.getByHash("CHDD")

	if result.(testSignal).Name != "room2" {
		t.Errorf("Signal server has bad state")
	}
}

func TestAnswerServer(t *testing.T) {
	server := AnswerServer{cache: make(map[string][]interface{})}
	server.push("ABCD", testAnswer{Name: "Mar"})
	server.push("ABCD", testAnswer{Name: "Jan"})
	server.push("DCBA", testAnswer{Name: "Oct"})

	results := server.getByHash("ABCD")
	signal := results[1]

	if signal.(testAnswer).Name != "Jan" {
		t.Errorf("Answer server has bad state")
	}

}
