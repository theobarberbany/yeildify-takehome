package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theobarberbany/yeildify-takehome/common"
)

// Server implements http.Handler
type Server struct {
}

// ServeHTTP checks the request body for a message, if one is found it is
// reversed and returned
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var m common.Message

		// Attempt to decode the request body into the Message struct. If there
		// is an error, respond to the client with the error message and a 400
		// status code
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&m)
		if err != nil {
			response := common.ErrorResponse{
				ErrorMessage: fmt.Sprintf("error unmarshalling json: %s ", err.Error()),
			}
			common.WriteJSONResponse(w, response, http.StatusBadRequest)

			return
		}

		// If there is no message passed to reverse, fail? Maybe I don't need this?
		if m.Message == "" {
			response := common.ErrorResponse{
				ErrorMessage: fmt.Sprintf("message passed is empty"),
			}
			common.WriteJSONResponse(w, response, http.StatusBadRequest)

			return
		}

		// Reverse the provided message
		reversed := reverse(m.Message)
		m.Message = reversed
		common.WriteJSONResponse(w, m, http.StatusOK)

		return

	default:
		response := common.ErrorResponse{
			ErrorMessage: fmt.Sprintf("method %s not implemented, please use POST only", r.Method),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
