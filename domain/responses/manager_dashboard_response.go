package responses

type Spark struct {
	Data []int64 `json:"data"`
}

type TopProduct struct {
	Name  string `json:"name"`
	Units int64  `json:"units"`
}

type ManagerDashboardResponse struct {
	TotalRevenueTHB int64       `json:"totalRevenueTHB"`
	NewMembers30d    int64       `json:"newMembers30d"`
	ActiveMembers    int64       `json:"activeMembers"`
	CheckinsToday    int64       `json:"checkinsToday"`
	CompletedPT30d   int64       `json:"completedPT30d"`

	RevenueSpark   Spark        `json:"revenueSpark"`
	NewMembersSpark Spark       `json:"newMembersSpark"`
	CheckinsSpark   Spark       `json:"checkinsSpark"`
	PTSpark         Spark       `json:"ptSpark"`

	TopProducts    []TopProduct `json:"topProducts"`
}