package domain

type UserExportBundle struct {
	Wallet      string          `json:"wallet"`
	Artifacts   []Artifact      `json:"artifacts"`
	Memberships []DAOMembership `json:"memberships"`
	Permissions []Permission    `json:"permissions"`
	Transfers   []Transfer      `json:"transfers"`
}
