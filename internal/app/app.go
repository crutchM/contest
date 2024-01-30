package app

import (
	. "Contest/internal/domain"
	"Contest/internal/server/handlers"
	"Contest/internal/services"
	"Contest/internal/storage"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	port           int
	router         *mux.Router
	store          *storage.Storage
	compileService services.ICompileService
	testService    services.ITestService
}

func New(cfg *Config) *App {
	router := mux.NewRouter()

	store, err := storage.NewStorage(cfg.ConnStr)
	if err != nil {
		panic(err)
	}

	compileService := services.NewCompileSevice()
	testService := services.NewTestService(compileService, store.TestRepository)

	app := &App{
		port:           cfg.Port,
		router:         router,
		store:          store,
		compileService: compileService,
		testService:    testService,
	}
	app.setupRouter()

	return app
}

func (a *App) setupRouter() {
	compileSubrouter := a.router.PathPrefix("/compile").Subrouter()
	compileSubrouter.HandleFunc("/cpp", handlers.CompileCPP(a.compileService)).Methods("POST")

	a.router.HandleFunc("/test", handlers.RunTest(a.testService)).Methods("GET")

	crudSubrouter := a.router.PathPrefix("/crud").Subrouter()
	crudSubrouter.HandleFunc("/test", handlers.AddTest(a.testService)).Methods("PUT")
	crudSubrouter.HandleFunc("/test/{id}", handlers.DeleteTest(a.testService)).Methods("DELETE")
	crudSubrouter.HandleFunc("/test/{id}", handlers.UpdateTest(a.testService)).Methods("PATCH")
	crudSubrouter.HandleFunc("/test/{id}", handlers.GetTest(a.testService)).Methods("GET")
	crudSubrouter.HandleFunc("/tests", handlers.GetTests(a.testService)).Methods("GET")
	crudSubrouter.HandleFunc("/tests/{task_id}", handlers.GetTestsByTaskID(a.testService)).Methods("GET")
}

func (a *App) MustRun() {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), a.router); err != nil {
		panic(err.Error())
	}
}
