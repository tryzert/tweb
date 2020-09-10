package tool


var UserPassKeySessions map[interface{}]interface{}


func init() {
	UserPassKeySessions = make(map[interface{}]interface{})
}


func UserLoginValidate(username, password string) bool {
	if username == "maple" && password == "maple"{
		return true
	}
	return false
}


func UserPassKeyValidate(username, passkey interface{}) bool {
	if username == nil || passkey == nil {
		return false
	}
	if UserPassKeySessions[username] == passkey {
		return true
	}
	return false
}