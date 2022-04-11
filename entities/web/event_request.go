package web

type EventRequest struct {
	Title string `form:"title"`
	HostedBy string `form:"hosted_by"`
	Cover string `form:"cover"`
	DatetimeEvent string `form:"datetime_event"`
	CategoryID uint `form:"category_id"`
	Location string `form:"location"`
	Description string `form:"description"`
}