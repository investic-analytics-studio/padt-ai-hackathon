package model

type FeaturesEntities struct {
	FeatureName   string `json:"feature_name" db:"feature_name"`
	FeatureEnable bool   `json:"feature_enable" db:"feature_enable"`
}
