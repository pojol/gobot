package openapi

type response struct {
	Code int
	Msg  string
	Body interface{}
}

type behaviorInfo struct {
	Name   string
	Update int64
	Status string
	Tags   []string
	Desc   string
}

type prefabInfo struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
