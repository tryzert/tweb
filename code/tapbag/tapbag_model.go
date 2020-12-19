package tapbag

import (
	"time"
)

type RequestContent struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ResponseContent struct {
	Code int         `json:"code"`
	Tip  string      `json:"tip"`
	Data interface{} `json:"data"`
}

type File struct {
	Id       int    `json:"id"`
	Type     string `json:"type"`
	Openable bool   `json:"openable"`
	Name     string `json:"name"`
	Relpath  string `json:"relpath"`
}

type Folder struct {
	Name        string    `json:"name"`
	Src         string    `json:"src"`
	HasChildren bool      `json:"hasChildren"`
	Children    []*Folder `json:"children"`
}

type FileItem struct {
	Paths         []string
	Size          int64
	VisitHot      int
	LastVisitTime time.Time
}
