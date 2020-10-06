package todolist

//数据库任务数据模型
type TaskModel struct {
	Id         int    `json:"id"`
	Task       string `json:"task"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Editstatus string `json:"editstatus"`
	Done       int    `json:"done"`
	Tag        string `json:"tag"`
	Deleted    int    `json:"deleted"`
	Deletetime string `json:"deletetime"`
	Lefttime   int    `json:"lefttime"`
}

//后端返回消息响应的模型
type ResponseContent struct {
	Code int         `json:"code"`
	Tip  string      `json:"tip"`
	Data interface{} `json:"data"`
}

//前端请求消息的模型
type RequestContent struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
