package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/coma-toast/pace-api/pkg/firebase"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/hcl/strconv"
	"github.com/rollbar/rollbar-go"
)

func main() {
	conf = getConf()

	// Rollbar logging setup
	rollbar.SetToken(conf.RollbarToken)
	rollbar.SetEnvironment("production")                    // defaults to "development"
	rollbar.SetCodeVersion("v0.0.1")                        // optional Git hash/branch/tag (required for GitHub integration)
	rollbar.SetServerHost("web.1")                          // optional override; defaults to hostname
	rollbar.SetServerRoot("github.com/coma-toast/pace-api") // path of project (required for GitHub integration and non-project stacktrace collapsing)
	rollbar.Info("PACE-API starting up...")
	rollbar.Wait()

	r := mux.NewRouter()
	// r.Use(authMiddle)
	// r.Handle("/api", authMiddle(blaHandler)).Methods(http.)
	// r.Methods("GET", "POST")
	r.HandleFunc("/api/ping", PingHandler)
	r.HandleFunc("/api/user/{userName}", GetUserHandler).Methods("GET")
	r.HandleFunc("/api/user", UpdateUserHandler).Methods("POST")
	r.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8001", r))
}

// TODO: auth middleware

// PingHandler is just a quick test to ensure api calls are working.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Dev code alert
	rollbar.Info(
		fmt.Sprintf("Ping test sent from %s", r.Header.Get("X-FORWARDED-FOR")))
	data := "Pong"
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&data); err != nil {
		log.Println("Error encoding JSON: ", err)
	}
	log.Println(data)
	// json.NewEncoder(w).Encode(data)
}

// add user example:
// 	_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
//         "first": "Ada",
//         "last":  "Lovelace",
//         "born":  1815,
// })
// if err != nil {
//         log.Fatalf("Failed adding alovelace: %v", err)
// }

// GetUserHandler handles api calls for User
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make(map[string]string, 0)
	// Get all the URL vars .../{userName}/{whatever}
	vars := mux.Vars(r)
	userName := vars["userName"]
	ctx := context.Background()
	db := firebase.Connect(conf.FirebaseConfig)
	users := db.Collection("users").Where("username", "==", userName).Documents(ctx)
	allMatchingUsers, err := users.GetAll()
	if err != nil {
		rollbar.Warning(
			fmt.Sprintf("Error getting user %s from Firebase: %e", userName, err))
	}
	for _, fbUser := range allMatchingUsers {
		spew.Dump(fbUser.data())
		data = append(data, fbUser.data())
	}

	encoder := json.NewEncoder(w)
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
