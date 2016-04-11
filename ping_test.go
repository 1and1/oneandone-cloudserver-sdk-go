package oneandone

import (
	"fmt"
	"testing"
)

// ping tests

func TestPing(t *testing.T) {
	fmt.Println("PING...")
	// API client with no token
	client := New("", BaseUrl)
	pong, err := client.Ping()
	if err != nil {
		t.Errorf("Ping failed. Error: " + err.Error())
	}
	if len(pong) == 0 {
		t.Errorf("Empty PING response.")
		return
	}
	if pong[0] != "PONG" {
		t.Errorf("Invalid PING response.")
	}
}

func TestPingAuth(t *testing.T) {
	fmt.Println("PING with authentication...")
	pong, err := api.PingAuth()
	if err != nil {
		t.Errorf("Ping with authentication failed. Error: " + err.Error())
	}
	if len(pong) == 0 {
		t.Errorf("Empty PING authentication response.")
		return
	}
	if pong[0] != "PONG" {
		t.Errorf("Invalid PING authentication response.")
	}
}
