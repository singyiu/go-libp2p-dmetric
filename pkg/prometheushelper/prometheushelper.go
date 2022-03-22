package prometheushelper

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RunServer(ctx context.Context, serverAddressStr string) {
	srv := &http.Server{Addr: serverAddressStr}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(context.TODO())
}
