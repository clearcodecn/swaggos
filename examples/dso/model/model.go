package model

type Group struct {
	Id             int64  `xorm:"id autoincr pk" json:"id"`
	Title          string `xorm:"title" json:"title"`
	Description    string `xorm:"description" json:"description"`
	OrganisationId int64  `xorm:"organisation_id" json:"organisationId"`
	KpiGroupId     int64  `xorm:"kpi_group_id" json:"kpiGroupId"`
	Subject        string `xorm:"subject" json:"subject"`

	CreatedAt int64 `xorm:"created" json:"createdAt"`
	UpdatedAt int64 `xorm:"updated" json:"updatedAt"`

	Items []*Item `json:"items" required:"true"`
}

func (s *Group) TableName() string {
	return "sequence_group"
}

type Item struct {
	Id             int64  `xorm:"id autoincr pk" json:"id"`
	Title          string `xorm:"title" json:"title"`
	Description    string `xorm:"description" json:"description"`
	OrganisationId int64  `xorm:"organisation_id" json:"organisationId"`
	KpiGroupId     int64  `xorm:"kpi_group_id" json:"kpiGroupId"`
	Subject        string `xorm:"subject" json:"subject"`

	CreatedAt int64 `xorm:"created" json:"createdAt"`
	UpdatedAt int64 `xorm:"updated" json:"updatedAt"`
}

type Group2 struct {
	Id             int64  `xorm:"id autoincr pk" json:"id"`
	Title          string `xorm:"title" json:"title"`
	Description    string `xorm:"description" json:"description"`
	OrganisationId int64  `xorm:"organisation_id" json:"organisationId"`
	KpiGroupId     int64  `xorm:"kpi_group_id" json:"kpiGroupId"`
	Subject        string `xorm:"subject" json:"subject"`

	CreatedAt int64   `xorm:"created" json:"createdAt"`
	UpdatedAt int64   `xorm:"updated" json:"updatedAt"`
	Items     []*Item `json:"items" required:"true"`
}
