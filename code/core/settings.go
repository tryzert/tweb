package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"tweb/code/tool"
)

//定义配置文件结构体
type Settings struct {
	Port     string   `json:"port"`
	Services Services `json:"services"`
}

type Services struct {
	TodolistConf *TodolistConf `json:"todolist_conf"`
	TapbagConf   *TapbagConf   `json:"tapbag_conf"`
}

type TodolistConf struct {
	Run bool `json:"run"`
}

type TapbagConf struct {
	Run     bool   `json:"run"`
	SrcPath string `json:"src_path"`
}

//初始化配置文件，settings.json
func initSettings() {
	services := Services{
		&TodolistConf{true},
		&TapbagConf{true, "./"},
	}
	st := &Settings{
		Port:     ":9010",
		Services: services,
	}
	data, err := json.MarshalIndent(st, "", "	")
	if err != nil {
		fmt.Println("配置信息初始化出错!")
		return
	}
	err = ioutil.WriteFile("settings.json", data, 0644)
	if err != nil {
		fmt.Println("配置信息写入出错！")
	}
}

func init() {
	//如果配置文件不存在或被删除了，则重新初始化一个
	if exist, _ := tool.FileExist("settings.json"); !exist {
		initSettings()
	}
}

//读取配置文件
func getSettings() *Settings {
	v := &Settings{}
	content, err := ioutil.ReadFile("settings.json")
	if err != nil {
		fmt.Println("配置信息读取失败！")
		return nil
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		fmt.Println("配置信息格式错误！")
		return nil
	}
	return v
}

//设置配置文件
func setSettings(st *Settings) error {
	content, err := json.Marshal(st)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("settings.json", content, 0644)
	if err != nil {
		fmt.Println("写入配置文件出错！")
	}
	return err
}
