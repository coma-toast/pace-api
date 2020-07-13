package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	helper "github.com/coma-toast/pace-api/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/hcl/strconv"
)

// TODO: save to DB instead of local json so it looks better for resume

func main() {
	r := mux.NewRouter()
	// r.Use(authMiddle)
	// r.Handle("/api", authMiddle(blaHandler)).Methods(http.)
	// r.Methods("GET", "POST")
	r.HandleFunc("/ping", PingHandler)
	r.HandleFunc("/api/paceData", GetPaceDataHandler).Methods("GET")
	r.HandleFunc("/api/paceData", UpdatePaceDataHandler).Methods("POST")
	r.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8000", r))
}

// TODO: auth middleware

// PingHandler is just a quick test to ensure api calls are working.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong\n"))
}

// GetPaceDataHandler handles api calls for paceData
func GetPaceDataHandler(w http.ResponseWriter, r *http.Request) {
	data := "test data - GetPaceDataHandler()"
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

// UpdatePaceDataHandler handles api calls for paceData
func UpdatePaceDataHandler(w http.ResponseWriter, r *http.Request) {
	var data string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println("Error decoding JSON: ", err)
	} else {
		helper.UpdateData(data)
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
			log.Println("Error sanitizing JSON ", err)
		}

		rdr1 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
		rdr2 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
		r.Body = rdr2
		log.Println(r.Method + ": " + r.RequestURI)
		log.Printf("BODY: %q", rdr1)
		next.ServeHTTP(w, r)
	})
}
