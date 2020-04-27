package servers

import (
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	appserver "github.com/codemodify/systemkit-appserver"
	reflection "github.com/codemodify/systemkit-helpers-reflection"
	logging "github.com/codemodify/systemkit-logging"
)

// HTTPRequestHandler -
type HTTPRequestHandler func(rw http.ResponseWriter, r *http.Request)

// HTTPHandler -
type HTTPHandler struct {
	Route   string
	Verb    string
	Handler HTTPRequestHandler
}

// HTTPServer -
type HTTPServer struct {
	handlers []HTTPHandler
}

// NewHTTPServer -
func NewHTTPServer(handlers []HTTPHandler) appserver.IServer {
	return &HTTPServer{
		handlers: handlers,
	}
}

// Run - Implement `IServer`
func (thisRef *HTTPServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *HTTPServer) PrepareRoutes(router *mux.Router) {
	for _, handler := range thisRef.handlers {
		logging.Debugf("%s - for %s, from %s", handler.Route, handler.Verb, reflection.GetThisFuncName())

		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb, "OPTIONS").Name(handler.Route)
	}
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *HTTPServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	if enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		err := http.Serve(listener, corsSetterHandler)
		if err != nil {
			logging.Fatalf("%s, from %s", err.Error(), reflection.GetThisFuncName())

			os.Exit(-1)
		}
	} else {
		err := http.Serve(listener, router)
		if err != nil {
			logging.Fatalf("%s, from %s", err.Error(), reflection.GetThisFuncName())

			os.Exit(-1)
		}
	}
}
