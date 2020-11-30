package http

type RequestEntity struct {
	InjectId string
	AppName  string
	Topic    string
	Opt      string
	RuleJson string
	Operator string
}

type ResponseEntity struct {
	Success  bool   `json:"success"`
	Msg      string `json:"msg"`
	ErrorMsg string `json:"errorMsg"`
}

type AwatchHandler interface {
	Handle(req *RequestEntity, resp *ResponseEntity)
}

const (
	TopicInjectRule   = "preset"
	TopicInjectManage = "awatch"
	TopicInjectSwitch = "switch"
)

var HandlerRepo = make(map[string]AwatchHandler)

func init() {
	HandlerRepo[TopicInjectRule] = new(InjectRuleHandler)
	HandlerRepo[TopicInjectManage] = new(InjectManageHandler)
	HandlerRepo[TopicInjectSwitch] = new(InjectSwitchHandler)
}
