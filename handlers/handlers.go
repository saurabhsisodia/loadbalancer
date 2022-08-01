package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// Register array of Endpoints
func Register(rw http.ResponseWriter, r *http.Request) {
	var urls []URL
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	Append(urls)
}

func Proxy(rw http.ResponseWriter, r *http.Request) {

	// no healthy upstreams
	if start == nil {
		http.Error(rw, "no healthy upstream", http.StatusServiceUnavailable)
		return
	}

	if nextServer == nil {
		nextServer = start
	}

	res, _, server := FindHealthyServerInRRFashion(nextServer, r)
	nextServer = server

	// check if there is no healthy Endpoints, after removing unhealthy Endpoints
	if res == nil && start == nil {
		http.Error(rw, "no healthy upstream", http.StatusServiceUnavailable)
		return
	}

	// copy all the header from res to rw
	for k, _ := range res.Header {
		rw.Header().Set(k, res.Header.Get(k))
	}

	// copy res.Body to rw
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func Get(rw http.ResponseWriter, r *http.Request) {

	curr := start
	urls := []URL{}
	for curr != nil {

		urls = append(urls, NewURL(curr.URL))
		if curr == end {
			break
		}
		curr = curr.Next

	}

	err := json.NewEncoder(rw).Encode(&urls)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

}

func Delete(rw http.ResponseWriter, r *http.Request) {
	start = nil
	end = nil
	nextServer = nil

	for k, _ := range mp {
		delete(mp, k)
	}

}
