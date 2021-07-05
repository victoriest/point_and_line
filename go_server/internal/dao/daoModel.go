package dao

type User struct {
	Id       int64
	Name     string
	Round    int
	WinCount int
	Rank     int
	Pwd      string
	OpenId   string
}
