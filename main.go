package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/hcl/strconv"
)

func main() {
	conf := getConf()
	_ = conf
	r := mux.NewRouter()
	// r.Use(authMiddle)
	// r.Handle("/api", authMiddle(blaHandler)).Methods(http.)
	// r.Methods("GET", "POST")
	r.HandleFunc("/ping", PingHandler)
	r.HandleFunc("/api/user", GetUserHandler).Methods("GET")
	r.HandleFunc("/api/user", UpdateUserHandler).Methods("POST")
	r.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8001", r))
}

// TODO: auth middleware

// PingHandler is just a quick test to ensure api calls are working.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong\n"))
}

// GetUserHandler handles api calls for User
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// db := firebase.Connect(conf.FirebaseConfig)
	data := "test data - GetUserHandler()"
	// data, err := helper.ReadSectorData()
	// if err != nil {
	// 	log.Panicln("Error decoding cached data", err)
	// }
	encoder := json.NewEncoder(w)
	// TODO: finish here  https://yourbasic.org/golang/json-example/#encode-marshal-struct-to-json
	if err := encoder.Encode(&data); err != nil {
		log.Println("Error encoding JSON: ", err)
	}
	log.Println(data)
}

// UpdateUserHandler handles api calls for User
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var data string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println("Error decoding JSON: ", err)
		// } else {
		// 	helper.UpdateData(data)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Settin'\n"))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			log.Print("bodyErr ", bodyErr.Error())
			http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
			return
		}

		unquoteJSONString, err := strconv.Unquote(string(buf))
		if err != nil {
			log.Println("Error sanitizing JSON: ", err)
		}

		rdr1 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
		rdr2 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
		r.Body = rdr2
		log.Println(r.Method + ": " + r.RequestURI)
		log.Printf("BODY: %q", rdr1)
		next.ServeHTTP(w, r)
	})
}
