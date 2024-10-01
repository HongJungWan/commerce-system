package response

type OrderStatsResponse struct {
	Month         string `json:"month"`
	TotalSales    int64  `json:"total_sales"`
	TotalCanceled int64  `json:"total_canceled"`
}
