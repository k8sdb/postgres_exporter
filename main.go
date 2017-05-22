package main

import (
	"flag"
	"net/http"
	"os"

	pge "github.com/k8sdb/postgres_exporter/exporter"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var Version string = "0.0.1"

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9187", "Address to listen on for web interface and telemetry.")
		metricPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		queriesPath   = flag.String("extend.query-path", "", "Path to custom queries to run.")
		onlyDumpMaps  = flag.Bool("dumpmaps", false, "Do not run, simply dump the maps.")
	)
	flag.Parse()

	if *onlyDumpMaps {
		pge.DumpMaps()
		return
	}

	dsn := os.Getenv("DATA_SOURCE_NAME")
	if len(dsn) == 0 {
		log.Fatal("couldn't find environment variable DATA_SOURCE_NAME")
	}

	exporter := pge.NewExporter(dsn, *queriesPath)
	prometheus.MustRegister(exporter)

	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		landingPage := []byte(`<html>
<head><title>Postgres exporter</title></head>
<body>
<h1>Postgres exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)
		w.Write(landingPage)
	})

	log.Infof("Starting Server: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
