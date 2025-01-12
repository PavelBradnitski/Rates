package models

type Rate struct {
	ID              uint    `gorm:"primaryKey"`
	CurID           int     `gorm:"not null" json:"Cur_ID"`
	Date            string  `json:"Date"`
	CurAbbreviation string  `json:"Cur_Abbreviation"`
	CurScale        int     `json:"Cur_Scale"`
	CurOfficialRate float64 `json:"Cur_OfficialRate"`
}
