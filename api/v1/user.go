package v1

type UserInfo struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Status      int    `json:"status"`
	Nickname    string `json:"nickname"`
	LastLoginAt string `json:"lastLoginAt"`
	LastLoginIP string `json:"lastLoginIP"`
}

type UserListRequest struct {
	Page     int    `json:"page" form:"page" default:"1"`
	PageSize int    `json:"pageSize" form:"pageSize" default:"10"`
	Username string `json:"username" form:"username"`
	Status   int8   `json:"status" form:"status"`
}

type UserListResponse struct {
	Total    int        `json:"total"`
	List     []UserInfo `json:"list"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
}
