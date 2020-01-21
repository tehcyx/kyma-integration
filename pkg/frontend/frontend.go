package frontend

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/tehcyx/kyma-integration/internal/tmpl"
)

var (
	templates = template.Must(template.New("").
		Funcs(template.FuncMap{}).Parse(""))
)

func init() {
	for _, tpl := range tmpl.TMPLMap {
		templates = template.Must(templates.Parse(tpl))
	}
}

// IndexHandler handler returning a tiny frontend to do the auto register process of the application
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("home --")

	var hostURL string
	if envIP := os.Getenv("INSTANCE_IP"); envIP != "" {
		hostURL = fmt.Sprintf("%s:8080", envIP)
	} else {
		hostURL = "http://localhost:8080"
	}

	isRedirect := false
	vals := r.URL.Query()
	if _, ok := vals["redirect"]; ok {
		isRedirect = true
	}

	hasError := false
	if _, ok := vals["error"]; ok {
		hasError = true
	}

	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
		"host":       hostURL,
		"isRedirect": isRedirect,
		"hasError":   hasError,
	}); err != nil {
		log.Error(err)
	}
}
