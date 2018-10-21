package main

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/travierm/PeerChat/server/lib"
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
	server := SignalServer{Cache: make(map[string]interface{})}
	server.Store("ABCD", testSignal{Name: "room1"})
	server.Store("CHDD", testSignal{Name: "room2"})

	result := server.GetByHash("CHDD")

	if result.(testSignal).Name != "room2" {
		t.Errorf("Signal server has bad state")
	}
}

func TestAnswerServer(t *testing.T) {
	server := AnswerServer{Cache: make(map[string][]interface{})}
	server.Push("ABCD", testAnswer{Name: "Mar"})
	server.Push("ABCD", testAnswer{Name: "Jan"})
	server.Push("DCBA", testAnswer{Name: "Oct"})

	results := server.GetByHash("ABCD")
	signal := results[1]

	if signal.(testAnswer).Name != "Jan" {
		t.Errorf("Answer server has bad state")
	}

}
