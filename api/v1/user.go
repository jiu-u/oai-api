package v1

type UserInfo struct {
	Id              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	LinuxDoId       string `json:"linuxDoId"`
	LinuxDoUsername string `json:"linuxDoUsername"`
}
