package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// Endpoint is a node in circular linked list
type EndPoint struct {
	URL            string
	Next           *EndPoint
	Last3Responses []Response // to check if Endpoint is healthy or not
}

type Response struct {
	StatusCode  int
	RequestTime time.Time
}

// Constructor
func NewResponse(code int, t time.Time) Response {
	return Response{code, t}
}

func NewEndPoint(url string) *EndPoint {
	return &EndPoint{URL: url}
}

// start is the starting of the circular linked list
// end is the end of the circular linked list
// nextServer will always point to the Endpoint which will be called in Round Robin approach
// mp is used to store distinct URLs in linked list
var (
	start      *EndPoint
	end        *EndPoint
	nextServer *EndPoint
	mp         = make(map[string]bool)
)

//Append new Endpoint at the end of Circular Linked List
func Append(urls []URL) {

	for _, url := range urls {

		if !mp[url.Url] {
			endpoint := NewEndPoint(url.Url)
			if start == nil {
				start = endpoint
			} else {
				end.Next = endpoint
			}
			end = endpoint
			end.Next = start
			mp[url.Url] = true
		}
	}

}

// this function will find next healthy Endpoint in Round Robin approach
// if there is no healthy Endpoints, then it will return nil response
func FindHealthyServerInRRFashion(nextServer *EndPoint, r *http.Request) (*http.Response, error, *EndPoint) {

	var (
		res        *http.Response = nil
		statusCode int
		err        error
		client     = &http.Client{}
	)

	for res == nil || statusCode != http.StatusOK {

		if nextServer == nil {
			return nil, nil, nextServer
		}

		fmt.Println("--------- trying to make request to -----", nextServer.URL)
		request, _ := http.NewRequest(r.Method, nextServer.URL, r.Body)

		res, err = client.Do(request)
		if res != nil {
			statusCode = res.StatusCode
		}

		// if there is no response or statuscode is not 200(OK)
		if err != nil || (res != nil && res.StatusCode != http.StatusOK) {

			// mark it as failure
			if len(nextServer.Last3Responses) >= 3 {
				// pop first response
				nextServer.Last3Responses = nextServer.Last3Responses[1:]
			}

			nextServer.Last3Responses = append(nextServer.Last3Responses, NewResponse(statusCode, time.Now()))

			// check if current endpoint becomes unhealthy
			// if yes, then remove it from registered Endpoints

			if len(nextServer.Last3Responses) == 3 {

				for _, response := range nextServer.Last3Responses {

					if response.StatusCode != http.StatusOK && IsTimeDifferenceLessThan15Seconds(response.RequestTime) {

						// remove Endpoint
						fmt.Println("-----------", nextServer.URL, "------ becomes unhealthy. removing it-----")

						nextServer = RemoveCurrentEndPoint(nextServer)
						break
					}
				}
			}

			if nextServer != nil {
				nextServer = nextServer.Next
			}

		} else {
			return res, nil, nextServer.Next // if we get valid response from current Endpoint
		}

	}
	return nil, err, nil
}

// calculate time difference between time.Now() and last failed request time
func IsTimeDifferenceLessThan15Seconds(t time.Time) bool {

	t2 := time.Now()
	d := t2.Sub(t).Seconds()
	if d <= float64(15) {
		return true
	}
	return false
}

// delete nextServer Endpoint from circular linked list
func RemoveCurrentEndPoint(nextServer *EndPoint) *EndPoint {

	mp[nextServer.URL] = false

	// if there is only one Endpoint
	if start == end {
		start = nil
		end = nil
		return nil
	}
	var (
		prev *EndPoint
	)
	curr := start
	for curr != nextServer {
		prev = curr
		curr = curr.Next
	}
	prev.Next = nextServer.Next
	if end == nextServer {
		end = prev
	}
	nextServer.Next = nil

	return prev
}

type URL struct {
	Url string `json:"url"`
}

func NewURL(url string) URL {
	return URL{url}
}
