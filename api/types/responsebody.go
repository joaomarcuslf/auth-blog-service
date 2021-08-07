package types

type ResponseBody struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
