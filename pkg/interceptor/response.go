package interceptor

const (
	DEFAULT_ERROR_MSG         = "Something went wrong"
	DEFAULT_HTTP_ERROR_CODE   = 500
	DEFAULT_HTTP_SUCCESS_CODE = 200
)

// Response declares the methods to access  the member-variables of response.
type Response interface {
	Error() interface{}
	ErrorMessage() string
	Success() bool
	Data() interface{}
}

// response implements Response and defines the structure of the response to send.
type response struct {
	ResSuccess      bool        `json:"success"`
	ResData         interface{} `json:"data"`
	ResError        interface{} `json:"error"`
	ResErrorMessage string      `json:"message"`
}

// CreateResponse creates a new istence of response using the input params and returns.
// If errMsg is not provided then "default error message" will be used.
func CreateResponse(success bool, data interface{}, err interface{}, errMsg string) Response {
	if !success && errMsg == "" {
		errMsg = DEFAULT_ERROR_MSG
	}
	return &response{
		ResSuccess:      success,
		ResData:         data,
		ResError:        err,
		ResErrorMessage: errMsg,
	}
}

// Error returns the cause of the error.
func (r *response) Error() interface{} {
	return r.ResError
}

// ErrorMessage returns an user friendly error message.
func (r *response) ErrorMessage() string {
	return r.ResErrorMessage
}

// Success returns boolean value representing the response status(eg: true from successfull response).
func (r *response) Success() bool {
	return r.ResSuccess
}

// Data returns the response data.
func (r *response) Data() interface{} {
	return r.ResData
}
