package tapbag

type RequestContent struct {
	Code int `json:"code"`
	Data string `json:"data"`
}

type ResponseContent struct {
	Code int `json:"code"`
	Tip string `json:"tip"`
	Data []*File `json:"data"`
}


type File struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Openable bool `json:"openable"`
	Name string `json:"name"`
	Relpath string `json:"relpath"`
}
