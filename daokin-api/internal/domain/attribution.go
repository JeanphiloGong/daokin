package domain

type AttributionTrail struct {
	ArtifactID   string       `json:"artifact_id"`
	Permissions  []Permission `json:"permissions"`
	Transfers    []Transfer   `json:"transfers"`
	TotalRecords int          `json:"total_records"`
}
