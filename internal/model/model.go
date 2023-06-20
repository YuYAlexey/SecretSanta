package model

type Users struct {
	Id         uint64 `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Sex        string `json:"sex"`
	Age        uint8  `json:"age"`
	Gift       uint64 `json:"gift"`
	Is_player  bool   `json:"is_player"`
}

type Gift struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Is_selected bool   `json:"is_selected"`
}

type Config struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type User_user struct {
	Giver     uint64 `json:"giver"`
	Recipient uint64 `json:"recipient"`
}
