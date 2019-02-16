package gateway

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusCode(err error) int {
	st := status.Convert(err)
	if st.Code() == codes.NotFound {
		return http.StatusNotFound
	} else {
		return http.StatusInternalServerError
	}
}
