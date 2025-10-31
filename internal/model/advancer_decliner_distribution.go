package model

type AdvancerDeclinerDistribution struct {
	Down15Percent     int `gorm:"column:down_15_pct" json:"down_15_pct"`
	Down10To15Percent int `gorm:"column:down_10_15_pct" json:"down_10_15_pct"`
	Down6To10Percent  int `gorm:"column:down_6_10_pct" json:"down_6_10_pct"`
	Down4To6Percent   int `gorm:"column:down_4_6_pct" json:"down_4_6_pct"`
	Down2To4Percent   int `gorm:"column:down_2_4_pct" json:"down_2_4_pct"`
	Down0To2Percent   int `gorm:"column:down_0_2_pct" json:"down_0_2_pct"`
	Even0Percent      int `gorm:"column:even_0_pct" json:"even_0_pct"`
	Up0To2Percent     int `gorm:"column:up_0_2_pct" json:"up_0_2_pct"`
	Up2To4Percent     int `gorm:"column:up_2_4_pct" json:"up_2_4_pct"`
	Up4To6Percent     int `gorm:"column:up_4_6_pct" json:"up_4_6_pct"`
	Up6To10Percent    int `gorm:"column:up_6_10_pct" json:"up_6_10_pct"`
	Up10To15Percent   int `gorm:"column:up_10_15_pct" json:"up_10_15_pct"`
	Up15Percent       int `gorm:"column:up_15_pct" json:"up_15_pct"`
}
