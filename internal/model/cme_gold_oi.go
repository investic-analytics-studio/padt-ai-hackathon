package model

type CMEGoldOI struct {
	MonthID     string `gorm:"column:monthid" json:"monthid"`
	Globex      int    `gorm:"column:globex" json:"globex"`
	OpenOutcry  int    `gorm:"column:openoutcry" json:"openoutcry"`
	TotalVolume int    `gorm:"column:totalvolume" json:"totalvolume"`
	BlockVolume int    `gorm:"column:blockvolume" json:"blockvolume"`
	EFPVol      int    `gorm:"column:efpvol" json:"efpvol"`
	EFRVol      int    `gorm:"column:efrvol" json:"efrvol"`
	EOOVol      int    `gorm:"column:eoovol" json:"eoovol"`
	EFSVol      int    `gorm:"column:efsvol" json:"efsvol"`
	SubVol      int    `gorm:"column:subvol" json:"subvol"`
	PNTVol      int    `gorm:"column:pntvol" json:"pntvol"`
	TASVol      int    `gorm:"column:tasvol" json:"tasvol"`
	Deliveries  int    `gorm:"column:deliveries" json:"deliveries"`
	OPNT        int    `gorm:"column:opnt" json:"opnt"`
	AON         int    `gorm:"column:aon" json:"aon"`
	AtClose     int    `gorm:"column:atclose" json:"atclose"`
	Change      int    `gorm:"column:change" json:"change"`
	Strike      string `gorm:"column:strike" json:"strike"`
	Exercises   int    `gorm:"column:exercises" json:"exercises"`
	Type        string `gorm:"column:type" json:"type"`
	Month       string `gorm:"column:month" json:"month"`
	Year        string `gorm:"column:year" json:"year"`
	LastUpdated string `gorm:"column:last_updated" json:"last_updated"`
}
