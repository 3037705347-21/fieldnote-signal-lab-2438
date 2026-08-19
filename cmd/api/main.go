package main

import (
	"example.com/fieldnote-signal-lab/internal/lab"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("FIELDNOTE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:18081"
	}
	service := lab.NewService(lab.NewMemoryStore(), lab.DefaultCatalog())
	lab.SeedDemo(service)
	log.Printf("fieldnote signal lab listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, lab.NewServer(service)))
}
