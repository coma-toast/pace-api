package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/container"
	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/paceconfig"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/hcl/hcl/strconv"
	"github.com/rollbar/rollbar-go"
)

// App is the app container
type App struct {
	Config    *paceconfig.Config
	Container container.Container
}

// TODO: look at Aaron's hub repo to see how to do the providers/connections.

// UserProvider provides a firestore Client for Users
type UserProvider struct {
	User *firestore.Client
}

// Run runs the app
func Run() {
	app := App{}
	var err error

	conf, err := paceconfig.GetConf()
	if err != nil {
		log.Fatalf("Error getting config: %s", err)
	}

	// Rollbar logging setup
	rollbar.SetToken(conf.RollbarToken)
	rollbar.SetEnvironment("production")                    // defaults to "development"
	rollbar.SetCodeVersion("v0.0.1")                        // optional Git hash/branch/tag (required for GitHub integration)
	rollbar.SetServerHost("web.1")                          // optional override; defaults to hostname
	rollbar.SetServerRoot("github.com/coma-toast/pace-api") // path of project (required for GitHub integration and non-project stacktrace collapsing)
	rollbar.Info("PACE-API starting up...")
	rollbar.Wait()

	app.Container = container.NewProduction(conf)

	r := mux.NewRouter()
	// r.Use(authMiddle)
	// r.Handle("/api", authMiddle(blaHandler)).Methods(http.)
	// r.Methods("GET", "POST")
	r.HandleFunc("/api/ping", PingHandler)
	r.HandleFunc("/api/user", app.GetUserHandler).Methods("GET")
	r.HandleFunc("/api/user", app.UpdateUserHandler).Methods("POST")
	// r.HandleFunc("/api/user", CreateUserHandler).Methods("PUT") // TODO: create user? or update auto creates?
	// r.Use(loggingMiddleware)
	// Gorilla Mux's logging handler.
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	log.Fatal(http.ListenAndServe(":8001", loggedRouter))
}

// TODO: auth middleware

// PingHandler is just a quick test to ensure api calls are working.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Dev code alert
	rollbar.Info(
		fmt.Sprintf("Ping test sent from %s", r.Header.Get("X-FORWARDED-FOR")), r)
	data := "Pong"
	jsonResponse(http.StatusOK, data, w)
}

// GetUserHandler handles api calls for User
func (a App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username")
	provider, err := a.Container.UserProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}
	if userName == "" {
		allUsers, err := provider.GetAll()
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting All Users: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, allUsers, w)
	} else {
		user, err := provider.GetByUsername(userName)
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting User: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, user, w)
	}
}

// UpdateUserHandler handles api calls for User
func (a App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a User: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.UserProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	updatedUser, err := provider.UpdateUser(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedUser, w)
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

func jsonResponse(statusCode int, v interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, _ := json.Marshal(v)
	w.Write(data)
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
			rollbar.Warning(fmt.Sprintf("Error sanitizing JSON: %s", err), r)
		}

		rdr1 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
		rdr2 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
		r.Body = rdr2
		log.Println(r.Method + ": " + r.RequestURI)
		log.Printf("BODY: %q", rdr1)
		next.ServeHTTP(w, r)
	})
}
