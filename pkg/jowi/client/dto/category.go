package dto

type ResponseCourseCategory struct {
	Status           int              `json:"status"`
	CourseCategories []CourseCategory `json:"course_categories"`
	Pagination
	ErrorResponse
}

type CourseCategory struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	DepartmentId string `json:"department_id"`
	ParentId     string `json:"parent_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
