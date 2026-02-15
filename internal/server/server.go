package server

import (
	"fmt"
	"net/http"

	"connectionInfo/internal/handler"
)

// Run starts the HTTP server on the specified port.
func Run(port string) error {
	h := handler.New()

	addr := fmt.Sprintf(":%s", port)
	return http.ListenAndServe(addr, h)
}
