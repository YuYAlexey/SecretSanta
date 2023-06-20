package model

type Users struct {
	Id         uint64 `json:"id"` // REVIEW: Лучше писать ID
	Login      string `json:"login"`
	Password   string `json:"password"`
	First_name string `json:"first_name"` // REVIEW: не правильное название с точки зрения Go. Навзови FirstName
	Last_name  string `json:"last_name"`  // REVIEW: не правильное название с точки зрения Go. Навзови LastName
	Sex        string `json:"sex"`
	Age        uint8  `json:"age"`
	Gift       uint64 `json:"gift"`
	Is_player  bool   `json:"is_player"` // REVIEW: не правильное название с точки зрения Go. Навзови IsPlayer
}

type Gift struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Is_selected bool   `json:"is_selected"` // REVIEW: не правильное название с точки зрения Go. Навзови IsSelected
}

type Config struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type User_user struct { // REVIEW: не правильное название с точки зрения Go. Навзови UserUser, но навернео лучше названить GiverRecipient (в БД тоже)
	Giver     uint64 `json:"giver"`
	Recipient uint64 `json:"recipient"`
}
