package resty

type requestType string

const (
	GET    requestType = "GET"
	PUT    requestType = "PUT"
	POST   requestType = "POST"
	PATCH  requestType = "PATCH"
	DELETE requestType = "DELETE"
)
