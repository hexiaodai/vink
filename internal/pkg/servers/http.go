package servers

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server interface {
	Run() error
	Stop() error
}

func httpErrorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(err.Error()))
}
