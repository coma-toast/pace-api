package cmd

import (
	"encoding/json"
	"flag"
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

	configPath := flag.String("conf", ".", "config path")
	flag.Parse()

	conf, err := paceconfig.GetConf(*configPath)
	if err != nil {
		log.Fatalf("Error getting config: %s", err)
	}

	// Rollbar logging setup
	rollbar.SetToken(conf.RollbarToken)
	rollbar.SetEnvironment("development")                   // defaults to "development"
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
		fmt.Println("No name provided")
		allProjects, err := provider.GetAll()
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting All Projects: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, allProjects, w)
	} else {
		user, err := provider.GetByName(projectName)
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting Project: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, user, w)
	}
}

// UpdateProjectHandler handles api calls for Project
// TODO: search by project name or by ID. Otherwise you can't update the project name
func (a App) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var project entity.UpdateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&project)
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

	updatedProject, err := provider.Update(project)
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
	var project entity.Project
	err := json.NewDecoder(r.Body).Decode(&project)
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

	err = provider.Delete(project)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting Project: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Project %s Deleted", project.Name), w)
}

// GetInventoryHandler Gets Inventory
func (a App) GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var inventory entity.Inventory
	inventoryID := r.URL.Query().Get("id")

	provider, err := a.Container.InventoryProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InventoryProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if inventoryID == "" {
		allInventory, err := provider.GetAll()
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting All Inventory: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, allInventory, w)
	} else {
		inventory, err = provider.GetByID(inventory.ID)
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting Inventory %s: %s", inventory.ID, err), r)
			jsonResponse(http.StatusInternalServerError, err.Error(), w)
			return

		}
		jsonResponse(http.StatusOK, inventory, w)
	}
}

// UpdateInventoryHandler Updates Inventory
func (a App) UpdateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var inventoryRequest entity.UpdateInventoryRequest
	err := json.NewDecoder(r.Body).Decode(&inventoryRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting Inventory: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InventoryProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InventoryProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	inventoryData, err := provider.Update(inventoryRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error updating Inventory: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, inventoryData, w)
}

// CreateInventoryHandler Creates Inventory
func (a App) CreateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var inventoryRequest entity.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventoryRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting Inventory: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InventoryProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InventoryProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	inventoryData, err := provider.Add(inventoryRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error adding Inventory item: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, inventoryData, w)
}

// DeleteInventoryHandler Deletes Inventory
func (a App) DeleteInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var inventoryRequest entity.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventoryRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when Inventory: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InventoryProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InventoryProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	err = provider.Delete(inventoryRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error deleting Inventory: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Inventory item %s deleted", inventoryRequest.ID), w)
}

// GetInspectionHandler Gets Inspections
func (a App) GetInspectionHandler(w http.ResponseWriter, r *http.Request) {
	inspectionID := r.URL.Query().Get("id")

	provider, err := a.Container.InspectionProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if inspectionID == "" {
		allInspections, err := provider.GetAll()
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error Getting All Inspections: %s", err), r)
			jsonResponse(http.StatusInternalServerError, err, w)
			return
		}
		jsonResponse(http.StatusOK, allInspections, w)
	} else {
		inspection, err := provider.GetByID(inspectionID)
		if err != nil {
			rollbar.Warning(fmt.Sprintf("Error getting Inspection %s: %s", inspection.ID, err), r)
			jsonResponse(http.StatusInternalServerError, err.Error(), w)
			return

		}
		jsonResponse(http.StatusOK, inspection, w)
	}
}

// UpdateInspectionHandler Updates Inspection
func (a App) UpdateInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var inspectionRequest entity.UpdateInspectionRequest
	err := json.NewDecoder(r.Body).Decode(&inspectionRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when Inspection: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	updatedInspection, err := provider.Update(inspectionRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, updatedInspection, w)
}

// CreateInspectionHandler Creates Inspection
func (a App) CreateInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var inspection entity.UpdateInspectionRequest
	err := json.NewDecoder(r.Body).Decode(&inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when updating Inspection: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	newInspection, err := provider.Add(inspection)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, newInspection, w)
}

// DeleteInspectionHandler Deletes Inspection
func (a App) DeleteInspectionHandler(w http.ResponseWriter, r *http.Request) {
	var inspectionRequest entity.Inspection
	err := json.NewDecoder(r.Body).Decode(&inspectionRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when deleting Inspection: %s", err), r)
		jsonResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	provider, err := a.Container.InspectionProvider()
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error decoding JSON when getting InspectionProvider: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	err = provider.Delete(inspectionRequest)
	if err != nil {
		rollbar.Warning(fmt.Sprintf("Error Inspection: %s", err), r)
		jsonResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	jsonResponse(http.StatusOK, fmt.Sprintf("Inspection  : %s", inspectionRequest), w)
}

func jsonResponse(statusCode int, v interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, _ := json.Marshal(v)
	if string(data) == "null" {
		data = []byte("[]")
	}
	w.Write(data)
}
