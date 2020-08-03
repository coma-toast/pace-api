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

	log.Fatal(http.ListenAndServe(":8001", app.getHandlers()))
}

func (a App) getHandlers() http.Handler {
	r := mux.NewRouter()
	// r.Use(authMiddle)
	// r.Handle("/api", authMiddle(blaHandler)).Methods(http.)
	// r.Methods("GET", "POST")
	r.HandleFunc("/api/ping", PingHandler)
	r.HandleFunc("/api/user", a.GetUserHandler).Methods("GET")
	r.HandleFunc("/api/user", a.UpdateUserHandler).Methods("POST")
	r.HandleFunc("/api/user", a.CreateUserHandler).Methods("PUT")
	r.HandleFunc("/api/user", a.DeleteUserHandler).Methods("DELETE")
	// TODO:  r.HandleFunc("/api/password", a.PasswordHandler).Methods("POST")
	r.HandleFunc("/api/contact", a.GetContactHandler).Methods("GET")
	r.HandleFunc("/api/contact", a.UpdateContactHandler).Methods("POST")
	r.HandleFunc("/api/contact", a.CreateContactHandler).Methods("PUT")
	r.HandleFunc("/api/contact", a.DeleteContactHandler).Methods("DELETE")
	r.HandleFunc("/api/company", a.GetCompanyHandler).Methods("GET")
	r.HandleFunc("/api/company", a.UpdateCompanyHandler).Methods("POST")
	r.HandleFunc("/api/company", a.CreateCompanyHandler).Methods("PUT")
	r.HandleFunc("/api/company", a.DeleteCompanyHandler).Methods("DELETE")
	r.HandleFunc("/api/project", a.GetProjectHandler).Methods("GET")
	r.HandleFunc("/api/project", a.UpdateProjectHandler).Methods("POST")
	r.HandleFunc("/api/project", a.CreateProjectHandler).Methods("PUT")
	r.HandleFunc("/api/project", a.DeleteProjectHandler).Methods("DELETE")
	r.HandleFunc("/api/inventory", a.GetInventoryHandler).Methods("GET")
	r.HandleFunc("/api/inventory", a.UpdateInventoryHandler).Methods("POST")
	r.HandleFunc("/api/inventory", a.CreateInventoryHandler).Methods("PUT")
	r.HandleFunc("/api/inventory", a.DeleteInventoryHandler).Methods("DELETE")
	r.HandleFunc("/api/inspection", a.GetInspectionHandler).Methods("GET")
	r.HandleFunc("/api/inspection", a.UpdateInspectionHandler).Methods("POST")
	r.HandleFunc("/api/inspection", a.CreateInspectionHandler).Methods("PUT")
	r.HandleFunc("/api/inspection", a.DeleteInspectionHandler).Methods("DELETE")

	// r.Use(loggingMiddleware)
	// Gorilla Mux's logging handler.
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	return loggedRouter
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
		jsonResponse(http.StatusBadRequest, err, w)
		return
	}

	provider, err := a.Container.UserProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
		return
	}

	updatedUser, err := provider.Update(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err, w)
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

	updatedUser, err := provider.Add(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error setting UserProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedUser, w)
}

// DeleteUserHandler deletes an existing user
func (a App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err = provider.Delete(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting User: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("User %s Deleted", user.Username), w)
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

	updatedUser, err := provider.Update(contact)
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

	updatedContact, err := provider.Add(contact)
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

	err = provider.Delete(contact)
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

	updatedUser, err := provider.Update(company)
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

	updatedUser, err := provider.Add(company)
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

	err = provider.Delete(company)
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

	updatedProject, err := provider.Update(user)
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

	updatedProject, err := provider.Add(user)
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

	err = provider.Delete(user)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting Project: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Project %s Deleted", user.Projectname), w)
}

// GetInventoryHandler Gets Inventory
func (a App) GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var inventory entity.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when inventory: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.tInventoryProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InventoryProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Get(inventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error getting Inventory %s: %s", inventory.ID, ))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("tInventory  : %s", inventory), w)
}
// UpdateInventoryHandler Updates Inventory
func (a App) UpdateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var eInventory entity.eInventory
	err := json.NewDecoder(r.Body).Decode(&eInventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when eInventory: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.eInventoryProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting eInventoryProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Upd(eInventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error eInventory"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("eInventory  : %s", eInventory), w)
}
// CreateInventoryHandler Creates Inventory
func (a App) CreateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var eInventory entity.eInventory
	err := json.NewDecoder(r.Body).Decode(&eInventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when eInventory: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.eInventoryProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting eInventoryProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Cre(eInventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error eInventory"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("eInventory  : %s", eInventory), w)
}
// DeleteInventoryHandler Deletes Inventory
func (a App) DeleteInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var eInventory entity.eInventory
	err := json.NewDecoder(r.Body).Decode(&eInventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when eInventory: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.eInventoryProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting eInventoryProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Del(eInventory)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error eInventory"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("eInventory  : %s", eInventory), w)
}
// GetInspectionHandler GetIs nspection
func (a App) GetInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var Inspection entity.Inspection
	err := json.NewDecoder(r.Body).Decode(&Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when Inspection: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Get(Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Inspection  : %s", Inspection), w)
}
// UpdateInspectionHandler UpdateIs nspection
func (a App) UpdateInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var Inspection entity.Inspection
	err := json.NewDecoder(r.Body).Decode(&Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when Inspection: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Upd(Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Inspection  : %s", Inspection), w)
}
// CreateInspectionHandler CreateIs nspection
func (a App) CreateInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var Inspection entity.Inspection
	err := json.NewDecoder(r.Body).Decode(&Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when Inspection: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Cre(Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Inspection  : %s", Inspection), w)
}
// DeleteInspectionHandler DeleteIs nspection
func (a App) DeleteInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var Inspection entity.Inspection
	err := json.NewDecoder(r.Body).Decode(&Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when Inspection: %w", err)r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
				rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %w", err)r)
				jsonResponse(http.StatusInternalServerError, err.Error(), w)
				return
	}
	
	err = provider.Del(Inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection"))
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Inspection  : %s", Inspection), w)
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
