package modelsv1

type Spacecraft struct {
	ID     int64  `json:"id" gorm:"primarykey"`
	Name   string `json:"name"`
	Class  string `json:"class"`
	Crew   string `json:"crew"`
	Image  string `json:"image"`
	Value  string `json:"value"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

func (s *Spacecraft) TableName() string {
	return "spacecraft"
}
