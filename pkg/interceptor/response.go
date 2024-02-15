package interceptor

const (
	DEFAULT_ERROR_MSG         = "Something went wrong"
	DEFAULT_HTTP_ERROR_CODE   = 500
	DEFAULT_HTTP_SUCCESS_CODE = 200
)

// Response declares the methods to access  the member-variables of Response.
type ResponseI interface {
	Error() interface{}
	ErrorMessage() string
	Success() bool
	Data() interface{}
}

// Response implements Response and defines the structure of the Response to send.
type Response struct {
	ResSuccess      bool        `json:"success"`
	ResData         interface{} `json:"data"`
	ResError        interface{} `json:"error"`
	ResErrorMessage string      `json:"message"`
}

// CreateResponse creates a new istence of Response using the input params and returns.
// If errMsg is not provided then "default error message" will be used.
func CreateResponse(success bool, data interface{}, err interface{}, errMsg string) ResponseI {
	if !success && errMsg == "" {
		errMsg = DEFAULT_ERROR_MSG
	}
	return &Response{
		ResSuccess:      success,
		ResData:         data,
		ResError:        err,
		ResErrorMessage: errMsg,
	}
}

// Error returns the cause of the error.
func (r *Response) Error() interface{} {
	return r.ResError
}

// ErrorMessage returns an user friendly error message.
func (r *Response) ErrorMessage() string {
	return r.ResErrorMessage
}

// Success returns boolean value representing the Response status(eg: true from successfull Response).
func (r *Response) Success() bool {
	return r.ResSuccess
}

// Data returns the Response data.
func (r *Response) Data() interface{} {
	return r.ResData
}
