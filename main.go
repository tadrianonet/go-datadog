package main

import (
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"net/http"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {

	rules := []tracer.SamplingRule{tracer.RateRule(1)}
	tracer.Start(
		tracer.WithSamplingRules(rules),
		tracer.WithService("service"),
		tracer.WithEnv("env"),
	)
	defer tracer.Stop()

	err := profiler.Start(
		profiler.WithService("service"),
		profiler.WithEnv("env"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	// Create a traced mux router
	mux := httptrace.NewServeMux()
	// Continue using the router as you normally would.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log := logrus.New()

		// Incluindo data e hora no log usando campos estruturados
		log.WithFields(logrus.Fields{
			"component": "log-validation",
			"status":    "starting",
			"datetime":  time.Now().Format(time.RFC3339),
		}).Info("Iniciando a validação de logs")

		w.Write([]byte("Hello World!"))
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
