package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/fittipaldi/go-r3-d3/app/apiv1/controllersv1"
	"github.com/fittipaldi/go-r3-d3/app/apiv1/middlewaresv1"
	"github.com/fittipaldi/go-r3-d3/app/model"
	"github.com/fittipaldi/go-r3-d3/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *sql.DB
	GormDB *gorm.DB
	Config *config.Config
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	a.Config = config

	//DATABASE CONNECTION
	a.DB = model.DatabaseConnection(config.DB)
	a.GormDB = model.InitGorm(config.DB.Type, config, a.DB)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "R3-D3 API")
	})

	// CUSTOM 404 NOT FOUND ROUTE
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{"message": "Not Found"}
		response, err := json.Marshal(payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(response))
	})

	a.Router = r
	//ROUTES TO API v1
	a.setApiV1Routes()
}

// setRouters sets the all required routers
func (a *App) setApiV1Routes() {
	commonMiddleware := middlewaresv1.CommonMiddlewareV3{
		GormDB: a.GormDB,
		Config: a.Config,
	}

	// Routing for handling the API-V1
	prefixApi := "/api/v1"
	r_v1 := a.Router.PathPrefix(prefixApi).Subrouter()
	r_v1.Use(commonMiddleware.Middleware)
	r_v1.HandleFunc("/", a.handleRequestV3(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API v1")
	})).Methods(http.MethodGet)

	//ROUTES USING API V1
	r_v1.HandleFunc("/spacecrafts", a.handleRequestV3(controllersv1.GetSpacecrafts)).Methods(http.MethodGet)

	l := 0
	a.Router.Walk(func(r *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tmpl, _ := r.GetPathTemplate()
		met, _ := r.GetMethods()

		matched, _ := regexp.MatchString(prefixApi, tmpl)
		if matched {
			ApiBase := strings.Replace(tmpl, prefixApi, "", 1)
			if len(ApiBase) > 1 {
				if l == 0 {
					fmt.Println("----------------------------- ROUTES V1 ---------------------------------")
				}
				mets := strings.Join(met, ",")
				fmt.Println(fmt.Sprintf("Endpoint[%s] Method[%s]", tmpl, mets))
				l++
			}
		}
		return nil
	})
}

// Run the app on it's router
func (a *App) Run(httpPorts string) {
	headersOk := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions})

	ports := strings.Split(httpPorts, ",")
	qtdPorts := len(ports)
	lastPort := ""
	for i, _ := range ports {
		idx := i + 1
		port := fmt.Sprintf(":%s", ports[i])
		if idx >= qtdPorts {
			lastPort = port
			break
		} else {
			go func() {
				fmt.Println(fmt.Sprintf("Server running in PORT [%s]", port))
				log.Fatal(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
			}()
		}
	}
	fmt.Println(fmt.Sprintf("Server running in PORT [%s]", lastPort))
	log.Fatal(http.ListenAndServe(lastPort, handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
}

type RequestHandlerFunctionV1 func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestV1(handler RequestHandlerFunctionV1) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, _ := a.checkAuthorizationHeader(r)
		if !auth {
			response, err := json.Marshal(map[string]interface{}{"message": "Unauthorized"})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(response))
			return
		}

		handler(a.DB, w, r)
	}
}

type RequestHandlerFunctionV3 func(w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestV3(handler RequestHandlerFunctionV3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip check if is a private port check, only for internal request.
		auth, httpCode := a.checkAuthorizationHeader(r)
		if !auth {
			message := "Unauthorized"
			switch httpCode {
			case 401:
				message = "Unauthorized"
				break
			case 404:
				message = "Not Found"
				break
			}
			payloadMsg := map[string]interface{}{"message": message}
			response, err := json.Marshal(payloadMsg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpCode)
			w.Write([]byte(response))
			return
		}
		handler(w, r)
	}
}

func (a *App) checkAuthorizationHeader(r *http.Request) (bool, int) {

	isApiv1 := strings.Split(r.URL.Path, "api/v1")
	if len(isApiv1) > 1 {
		return middlewaresv1.CheckTheAuthorization(r, a.GormDB)
	}

	return false, http.StatusUnauthorized
}
