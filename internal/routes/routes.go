package routes

import (
	"net/http"
	//"os"

	apihandlers "rainier/internal/handlers"
	mw "rainier/internal/middleware"
	"rainier/internal/util"

	"github.com/rs/zerolog/log"

	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//GetRoutes get the routes for the application
func GetRoutes(hctx *apihandlers.HTTPHandlerContext, mwctx *mw.MiddleWareContext) *mux.Router {

	staticAssets, err := util.SiteTemplate("/static")
	if err != nil {
		log.Fatal().Err(err).Str("service", "wpengine").Msg("Please ensure the web template directory exists and that you have permissions to access it")
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	// Index Page
	r.Handle(
		"/",
		mwctx.WPLog(
			http.HandlerFunc(hctx.IndexHandler))).Methods("GET", "POST")

	// Test
	r.Handle(
		"/test",
		mwctx.WPLog(
			http.HandlerFunc(hctx.TestHandler))).Methods("GET")

	// Health check page
	r.Handle(
		"/healthcheck957873",
		http.HandlerFunc(hctx.ErrorHandler)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(apihandlers.NotFound)

	// Static pages
	ServeStatic(r, "./"+staticAssets)

	// Default handler
	http.Handle("/", r)

	return r
}

// ServeStatic  serve static content from the appropriate location
func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"css":     staticDirectory + "/css/",
		"images":  staticDirectory + "/images/",
		"scripts": staticDirectory + "/scripts/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}
