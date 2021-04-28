package modelsv1

type Token struct {
	ID     int64  `json:"id" gorm:"primarykey"`
	Token  string `json:"token"`
	Status int    `json:"status"`
	IsJedi bool   `json:"is_jedi"`
}

func (t *Token) TableName() string {
	return "tokens"
}
