package dto

type ResponseCourse struct {
	Status  int      `json:"status"`
	Courses []Course `json:"courses"`

	Pagination
	ErrorResponse
}

type Course struct {
	Id                  string `json:"id"`
	CourseCategoryId    string `json:"course_category_id"`
	Title               string `json:"title"`
	Price               string `json:"price"`                  // WRONG JOWI DOCS TYPE
	PriceForOnlineOrder string `json:"price_for_online_order"` // WRONG JOWI DOCS TYPE
	IsException         bool   `json:"is_exception"`
	OnlineOrder         bool   `json:"online_order"`
	IsVisible           bool   `json:"is_visible"`
	IsPiece             bool   `json:"is_piece"`
	CostPrice           string `json:"cost_price"` // WRONG JOWI DOCS TYPE
	PrepareTime         string `json:"prepare_time"`
	UnitName            string `json:"unit_name"`
	Weight              string `json:"weight"`
	Description         string `json:"description"`
	ImageUrl            string `json:"image_url"`
	ImageUpdatedAt      string `json:"image_updated_at"`
	ParentId            string `json:"parent_id"`
	CompanySaleId       string `json:"company_sale_id"`
	TaxGroupId          string `json:"tax_group_id"`
	CookingTechnology   string `json:"cooking_technology"`
	FullOut             string `json:"full_out"`
	KKL                 string `json:"kkl"`
	KDJ                 string `json:"kdj"`
	IKPU                string `json:"ikpu"`
	PackageCode         string `json:"package_code"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
