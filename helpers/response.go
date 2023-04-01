package helpers

type BaseResponse struct {
	Meta struct {
		Message string   `json:"message"`
		Errors  []string `json:"error,omitempty"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

func NewSuccessResponse(param interface{}) BaseResponse {
	response := BaseResponse{}
	response.Meta.Message = "Success"
	response.Data = param

	return response
}

func NewSuccessInsertResponse(param interface{}) BaseResponse {
	response := BaseResponse{}
	response.Meta.Message = "Success Insert"
	response.Data = param

	return response
}

func NewErrorResponse(err error) BaseResponse {
	response := BaseResponse{}
	response.Meta.Message = "Something wrong"
	response.Meta.Errors = []string{err.Error()}

	return response
}
