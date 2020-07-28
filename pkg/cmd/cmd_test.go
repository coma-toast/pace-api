package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRoute(t *testing.T) {
	a := App{}
	testingServer := httptest.NewServer(a.getHandlers())
	defer testingServer.Close()
	response, err := http.Get(fmt.Sprintf("%s/api/ping", testingServer.URL))
	if err != nil {
		t.Error("Error getting Ping response: ", err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	if !reflect.DeepEqual(body, []byte("\"Pong\"")) {
		t.Error("Body does not match expected value: ", body)
	}
}

func TestUser(t *testing.T) {
	a := App{}
	testingServer := httptest.NewServer(a.getHandlers())
	defer testingServer.Close()
	response, err := http.Get(fmt.Sprintf("%s/api/ping", testingServer.URL))
	if err != nil {
		t.Error("Error getting Ping response: ", err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	if !reflect.DeepEqual(body, []byte("\"Pong\"")) {
		t.Error("Body does not match expected value: ", body)
	}
}
