package models

type FeatureStatus string

const (
	CheckStatus        FeatureStatus = "checked"
	RejectStatus       FeatureStatus = "rejected"
	InProgressStatus   FeatureStatus = "in-progress"
	IntermediateStatus FeatureStatus = "intermediate"
)
