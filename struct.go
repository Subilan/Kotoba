package kotoba

type reqCommonLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type reqCommonRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
}