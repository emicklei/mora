package response

import restful "github.com/emicklei/go-restful"

type Response struct {
	Success bool        `json:"success"`
	Error   *Error      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteSuccess(resp *restful.Response) {
	NewResponse(true).WriteStatus(200, resp)
}

func WriteResponse(data interface{}, resp *restful.Response) {
	WriteResponseStatus(200, data, resp)
}

func WriteResponseStatus(status int, data interface{}, resp *restful.Response) {
	success := NewResponse(true)
	success.Data = data
	success.WriteStatus(status, resp)
}

func NewResponse(success bool) *Response {
	return &Response{Success: success}
}

func NewErrorResponse(err error) *Response {
	res := &Response{Success: false, Error: &Error{}}
	res.SetError(err)
	return res
}

func (r *Response) SetError(err error) {
	if err != nil {
		if r.Error == nil {
			r.Error = &Error{}
		}
		r.Error.Name = err.Error()
	}
}

func (r *Response) WriteStatus(status int, resp *restful.Response) {
	if r.Error != nil && r.Error.Code == 0 {
		r.Error.Code = status
	}
	resp.WriteHeaderAndEntity(status, r)
}
