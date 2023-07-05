package model

type Users struct {
	ID        uint64 `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       string `json:"sex"`
	Age       string `json:"age"`
	Gift      string `json:"gift"`
	IsPlayer  bool   `json:"is_player"`
}

type Gift struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Description string `json:"description"`
	IsSelected  bool   `json:"is_selected"`
}

type Config struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type GiverRecipient struct {
	Giver     uint64 `json:"giver"`
	Recipient uint64 `json:"recipient"`
}
