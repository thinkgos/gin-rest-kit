package httputil

type Response struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	TraceId string `json:"traceId,omitempty"`
	Detail  string `json:"detail,omitempty"`
	Data    any    `json:"data"`
}
