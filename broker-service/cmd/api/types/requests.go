package types

type JsonResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RouteRequestBody struct {
	Action  string `json:"action"`
	Payload any    `json:"payload,omitempty"`
}

type MethodCallInfo struct {
	Method   string
	Endpoint string
	Body     interface{}
}

type RabbitPayload struct {
	Endpoint string `endpoint:"endpoint"`
	Method   string `endpoint:"method"`
	Data     any    `json:"data,omitempty"`
}
