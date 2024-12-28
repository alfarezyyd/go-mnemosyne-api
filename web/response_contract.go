package web

type ResponseContract struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewResponseContract(statusResp bool, messageResp string, dataResp interface{}, Error interface{}) *ResponseContract {
	return &ResponseContract{
		Status:  statusResp,
		Message: messageResp,
		Data:    dataResp,
		Error:   Error,
	}
}
