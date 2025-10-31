package model

type TopRankResponse struct {
	TopList []string `json:"top_list"`
}

type TopRank struct {
	Rank      int    `db:"rank"`
	BaseAsset string `db:"base_asset"`
}
