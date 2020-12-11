package tapbag

type RequestContent struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
}

type ResponseContent struct {
	Code int `json:"code"`
	Tip string `json:"tip"`
	Data interface{} `json:"data"`
}


type File struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Openable bool `json:"openable"`
	Name string `json:"name"`
	Relpath string `json:"relpath"`
}
