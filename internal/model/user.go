package model

type User struct {
	Id      uint     `json:"id"`
	Name    string   `json:"user_name"`
	Age     int      `json:"age"`
	Friends []uint64 `json:"friends"`
}
