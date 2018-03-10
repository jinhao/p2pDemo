package proto

const (
	CMD_LOGIN = iota + 1
	CMD_LOGOUT
	CMD_LIST
	CMD_CONE
	CMD_MSG
)

type Proto struct {
	Cmd      int      `json:cmd`
	ConeAddr string   `json:coneAddr`
	Clients  []string `json: clients`
	Data     string   `json: data`
}
