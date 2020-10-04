package todolist

//数据库任务数据模型
type ThingModel struct {
	Id         int    `json:"id"`
	Thing      string `json:"thing"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Editstatus string `json:"editstatus"`
	Done       int    `json:"done"`
	Tag        string `json:"tag"`
	Isdeleted  int    `json:"isdeleted"`
	Deletetime string `json:"deletetime"`
	Lefttime   int    `json:"lefttime"`
}

//后端返回消息响应的模型
type ResponseContent struct {
	ResponseCode int         `json:"response_code"`
	Tip          string      `json:"tip"`
	Data         interface{} `json:"data"`
}
