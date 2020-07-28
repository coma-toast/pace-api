package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/container"
	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/paceconfig"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	r.HandleFunc("/api/user", app.CreateUserHandler).Methods("PUT")
	r.HandleFunc("/api/user", app.DeleteUserHandler).Methods("DELETE")
	// TODO:  r.HandleFunc("/api/password", app.PasswordHandler).Methods("POST")
	r.HandleFunc("/api/contact", app.GetContactHandler).Methods("GET")
	r.HandleFunc("/api/contact", app.UpdateContactHandler).Methods("POST")
	r.HandleFunc("/api/contact", app.CreateContactHandler).Methods("PUT")
	r.HandleFunc("/api/contact", app.DeleteContactHandler).Methods("DELETE")
	r.HandleFunc("/api/company", app.GetCompanyHandler).Methods("GET")
	r.HandleFunc("/api/company", app.UpdateCompanyHandler).Methods("POST")
	r.HandleFunc("/api/company", app.CreateCompanyHandler).Methods("PUT")
	r.HandleFunc("/api/company", app.DeleteCompanyHandler).Methods("DELETE")
	r.HandleFunc("/api/project", app.GetProjectHandler).Methods("GET")
	r.HandleFunc("/api/project", app.UpdateProjectHandler).Methods("POST")
	r.HandleFunc("/api/project", app.CreateProjectHandler).Methods("PUT")
	r.HandleFunc("/api/project", app.DeleteProjectHandler).Methods("DELETE")

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

// GetContactHandler handles api calls for contacts
func (a App) GetContactHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := a.Container.ContactProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ContactProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
	}
	allContacts, err := provider.GetAll()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting All Contacts: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}
	jsonResponse(http.StatusOK, allContacts, w)
}

// UpdateContactHandler handles api calls for contacts
func (a App) UpdateContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact entity.Contact
	provider, err := a.Container.ContactProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ContactProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
	}

	err = json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a User: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	updatedUser, err := provider.UpdateContact(contact)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting ContactProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedUser, w)
}

// CreateContactHandler handles api calls for contacts
func (a App) CreateContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact entity.Contact
	provider, err := a.Container.ContactProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ContactProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
	}

	err = json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when creating a Contact: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	updatedContact, err := provider.AddContact(contact)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting ContactProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedContact, w)
}

// DeleteContactHandler handles api calls for contacts
func (a App) DeleteContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact entity.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a contact: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.ContactProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting contactProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	err = provider.DeleteContact(contact)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting contact: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("contact %s %s  Deleted", contact.FirstName, contact.LastName), w)
}

//GetCompanyHandler handles api calls for Company
func (a App) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := a.Container.CompanyProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting CompanyProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}
	allCompanies, err := provider.GetAll()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting All Companies: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}
	jsonResponse(http.StatusOK, allCompanies, w)
}

//UpdateCompanyHandler handles api calls for Company
func (a App) UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var company entity.Company
	provider, err := a.Container.CompanyProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting CompanyProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a Company: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	updatedUser, err := provider.UpdateCompany(company)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting CompanyProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedUser, w)
}

//CreateCompanyHandler handles api calls for Company
func (a App) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var company entity.Company
	provider, err := a.Container.CompanyProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting CompanyProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when creating a Company: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	updatedUser, err := provider.AddCompany(company)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting CompanyProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedUser, w)
}

//DeleteCompanyHandler handles api calls for Company
func (a App) DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var company entity.Company
	provider, err := a.Container.CompanyProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting CompanyProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when deleting a Company: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = provider.DeleteCompany(company)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting company: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("company %s Deleted", company.Name), w)
}

// GetProjectHandler handles api calls for User
func (a App) GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	projectName := r.URL.Query().Get("name")
	provider, err := a.Container.ProjectProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ProjectProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}
	if projectName == "" {
		allProjects, err := provider.GetAll()
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting All Projects: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, allProjects, w)
	} else {
		user, err := provider.GetByProjectname(projectName)
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting Project: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, user, w)
	}
}

// UpdateProjectHandler handles api calls for Project
func (a App) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.UpdateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a Project: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.ProjectProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ProjectProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	updatedProject, err := provider.UpdateProject(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting ProjectProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedProject, w)
}

// CreateProjectHandler adds a new user
func (a App) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.Project
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a Project: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.ProjectProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ProjectProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	updatedProject, err := provider.AddProject(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting ProjectProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedProject, w)
}

// DeleteProjectHandler deletes an existing user
func (a App) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.UpdateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating a Project: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.ProjectProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting ProjectProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	err = provider.DeleteProject(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting Project: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Project %s Deleted", user.Projectname), w)
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

// CreateUserHandler adds a new user
func (a App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
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

	updatedUser, err := provider.AddUser(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedUser, w)
}

// DeleteUserHandler deletes an existing user
func (a App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err = provider.DeleteUser(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting User: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("User %s Deleted", user.Username), w)
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

// TODO: remove or refactor. Switched to the Gorilla Mux logging middleware.
// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()
// 		buf, bodyErr := ioutil.ReadAll(r.Body)
// 		if bodyErr != nil {
// 			log.Print("bodyErr ", bodyErr.Error())
// 			http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		unquoteJSONString, err := strconv.Unquote(string(buf))
// 		if err != nil {
// 			rollbar.Warning(fmt.Sprintf("Error sanitizing JSON: %s", err), r)
// 		}

// 		rdr1 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
// 		rdr2 := ioutil.NopCloser(strings.NewReader(unquoteJSONString))
// 		r.Body = rdr2
// 		log.Println(r.Method + ": " + r.RequestURI)
// 		log.Printf("BODY: %q", rdr1)
// 		next.ServeHTTP(w, r)
// 	})
// }
