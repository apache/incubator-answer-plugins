package wecom

type Department struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	ParentID         int      `json:"parentid"`
	Order            int      `json:"order"`
	DepartmentLeader []string `json:"department_leader"`
}

type Employee struct {
	Name       string `json:"name"`
	Department []int  `json:"department"`
	Userid     string `json:"userid"`
}

type AuthUserInfoResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Userid  string `json:"userid"`
	Mobile  string `json:"mobile"`
	Gender  string `json:"gender"`
	Email   string `json:"email"`
	Avatar  string `json:"avatar"`
	QrCode  string `json:"qr_code"`
	Address string `json:"address"`
}

type UserInfo struct {
	Userid        string `json:"userid"`
	Mobile        string `json:"mobile"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	QrCode        string `json:"qr_code"`
	Address       string `json:"address"`
	Name          string `json:"name"`
	DepartmentIDs []int  `json:"department"`
	Position      string `json:"position"`
	IsAvailable   bool   `json:"is_available"`
}

type UserDetailInfo struct {
	Errcode        int    `json:"errcode"`
	Errmsg         string `json:"errmsg"`
	Userid         string `json:"userid"`
	Name           string `json:"name"`
	Department     []int  `json:"department"`
	Position       string `json:"position"`
	Status         int    `json:"status"`
	Isleader       int    `json:"isleader"`
	EnglishName    string `json:"english_name"`
	Telephone      string `json:"telephone"`
	Enable         int    `json:"enable"`
	HideMobile     int    `json:"hide_mobile"`
	Order          []int  `json:"order"`
	MainDepartment int    `json:"main_department"`
	Alias          string `json:"alias"`
}
