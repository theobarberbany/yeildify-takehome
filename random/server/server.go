package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"log"

	"github.com/theobarberbany/yeildify-takehome/common"
)

// Server implements http.Handler. It procides a http.Client for making requests
type Server struct {
	// c is an http client
	c *http.Client
}

// Initialise sets up the server's http client and initialises a random seed
func (s *Server) Initialise() {
	s.c = &http.Client{}
	rand.Seed(time.Now().UnixNano())

	log.Print("server initialised")

	return
}

// these could be read in as command line args
const (
	reverseEndpoint = "reverse"
	reverseHost     = "reverse"
	reversePort     = "80"
	proxyScheme     = "http"
)

// ServeHTTP generates a random number then makes a request to the reversal
// service before returning the result
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// pass request to reverse service
		proxyReq, err := newProxyRequest(r)
		if err != nil {
			log.Printf("error building proxy request: %s", err.Error())
			response := common.ErrorResponse{
				ErrorMessage: fmt.Sprintf("error building proxy request: %s", err.Error()),
			}
			common.WriteJSONResponse(w, response, http.StatusInternalServerError)

			return
		}

		// possibly degrade gracefully here? wrap in a retry function
		// in case of intermittent network failure?
		log.Printf("proxying request: %s", proxyReq.URL.String())
		r, err := s.c.Do(proxyReq)
		if err != nil {
			log.Printf("error making request to reverse service: %s ", err.Error())
			response := common.ErrorResponse{
				ErrorMessage: fmt.Sprintf("error making request to reverse service: %s ", err.Error()),
			}
			common.WriteJSONResponse(w, response, http.StatusBadGateway)

			return
		}

		// If the reverse service fails for some reason, return an error along
		// with the body of the response
		if r.StatusCode != 200 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			log.Printf("error making request to reverse service: returned %d: %s", r.StatusCode, buf.String())
			response := common.ErrorResponse{
				ErrorMessage: fmt.Sprintf("error making request to reverse service: returned %d: %s", r.StatusCode, buf.String()),
			}
			common.WriteJSONResponse(w, response, http.StatusBadRequest)

			return
		}

		log.Printf("request %s completed ok", proxyReq.URL.String())

		// Attempt to decode the request body into the Message struct. If there
		// is an error, respond to the client with the error message and a 500
		// status code. If this is malformed, something has gone seriously
		// wrong.
		var m common.Message
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err = decoder.Decode(&m)
		if err != nil {
			log.Printf("error unmarshalling json: %s ", err.Error())
			response := common.ErrorResponse{
				ErrorMessage: fmt.Sprintf("error unmarshalling json: %s ", err.Error()),
			}
			common.WriteJSONResponse(w, response, http.StatusInternalServerError)

			return
		}

		// obtain random number
		randFloat := rand.Float64()
		m.RandomNumber = randFloat
		common.WriteJSONResponse(w, m, http.StatusOK)

		return

	default:
		response := common.ErrorResponse{
			ErrorMessage: fmt.Sprintf("method %s not implemented, please use POST only", r.Method),
		}

		common.WriteJSONResponse(w, response, http.StatusBadRequest)

		return
	}
}

// buildNewURL builds a new URL for proxying a reverse request. It takes an
// incoming http.Request and parses it, changing the port on the host and the
// endpoint requested
func buildNewURL(r *http.Request) string {
	splitURL := strings.Split(r.URL.String(), "/")
	splitURL = splitURL[:len(splitURL)-1]

	return fmt.Sprintf("%s://%s:%s%s/%s", proxyScheme, reverseHost, reversePort, strings.Join(splitURL, "/"), reverseEndpoint)
}

// newProxyRequest builds a new request to be proxied
func newProxyRequest(r *http.Request) (*http.Request, error) {
	newURL := buildNewURL(r)
	req, err := http.NewRequest(r.Method, newURL, r.Body)
	if err != nil {
		return nil, err
	}

	req.Header = r.Header
	req.Header.Set("Host", r.Host)
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)

	return req, nil
}
