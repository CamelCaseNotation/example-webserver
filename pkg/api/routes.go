package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var l *logrus.Logger

// RootRouter returns a gorilla Router to be used as the Handler for the http.Server struct.
// TODO(josh): add additional routes + authentication
func RootRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	addV1Routes(r)
	return r
}

// V1Handler returns a Router struct which handles all requests prefixed with "/v1/"
func addV1Routes(r *mux.Router) {
	routerV1 := r.PathPrefix("/v1").Subrouter()
	routerV1.HandleFunc("/", randomNameJoke)
}

// Log returns a formatted and configured logrus.Logger struct such that
// the file name and line number is logged when used. Should be called by each handler
// which means they are responsible for setting the fields prior to logging.
// Log responsibly and usefully!
// TODO(josh): Ensure this does not add too much performance overhead at scale
func Log() *logrus.Logger {
	l = logrus.New()
	l.SetReportCaller(true)
	// TODO(josh): change to file/other if we're not cloud-native... but we should be :)
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			// https://github.com/sirupsen/logrus/issues/63#issuecomment-548792922
			return funcname, filename + ":" + strconv.Itoa(f.Line)
		},
	}
	return l
}