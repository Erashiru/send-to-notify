package dto

type ResponseStopList struct {
	Status       int           `json:"status"`
	CourseCounts []CourseCount `json:"course_counts"`
	ErrorResponse
}

type CourseCount struct {
	Id               string `json:"id"`
	CourseCategoryId string `json:"course_category_id"`
	DeviceId         string `json:"device_id"`
	Title            string `json:"title"`
	Count            string `json:"count"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
