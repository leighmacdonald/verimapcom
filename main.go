package main

import (
	"context"
	"github.com/leighmacdonald/verimapcom/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	w := web.New(ctx)
	if err := w.Setup(); err != nil {
		log.Fatalf("Could not run setup: %v", err)
	}
	defer w.Close()

	opts := web.DefaultHTTPOpts()
	opts.Handler = w.Handler

	srv := web.NewHTTPServer(opts)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Errorf("Shutdown unclean: %v", err)
	}
}
