package core
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


type Settings struct {
	RootDir string `json:"rootdir"`
	Port string `json:"port"`
	Services Services `json:"services"`
}

type Services struct {
	Run_todolist bool `json:"run_todolist"`
	Run_tim      bool `json:"run_tim"`
}


//初始化配置文件，settings.json
func initSettings() {
	services := Services{
		Run_todolist: true,
		Run_tim:      true,
	}
	st := &Settings{
		RootDir: "/media/maple/E盘",
		Port: ":9010",
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
	_, err := os.Stat("settings.json")
	//如果配置文件不存在或被删除了，则重新初始化一个
	if os.IsNotExist(err) {
		initSettings()
	}
}


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

