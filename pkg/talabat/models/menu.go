package models

type GetRequestStatusResponse struct {
	TimeStamp string `json:"timeStamp"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Status    int32  `json:"status"`
}

type CreateNewMenuRequest struct {
	RestaurantID string
	RequestID    string
	Menu         Menu
}

type CreateNewMenuErrorResponse struct {
	TimeStamp string `json:"timeStamp"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Status    int32  `json:"status"`
}
type BranchAvailability struct {
	BranchID string  `json:"branchId"`
	Status   bool    `json:"status"`
	Price    float64 `json:"price"`
}

type Branch struct {
	ID          int    `json:"id"`
	NameEn      string `json:"nameEn"`
	IsAvailable bool   `json:"isAvailable"`
}

type Schedule struct {
	DayOfWeek int    `json:"dayOfWeek"`
	DayName   string `json:"dayName"`
	TimeFrom  string `json:"timeFrom"`
	TimeTo    string `json:"timeTo"`
	AllDay    bool   `json:"allDay"`
}

type ItemDiscount struct {
	Name                  string     `json:"name"`
	CouponStatus          int        `json:"couponStatus"`
	MarketingTitleEn      string     `json:"marketingTitleEn"`
	MarketingTitleAr      string     `json:"marketingTitleAr"`
	DescriptionEn         string     `json:"descriptionEn"`
	DescriptionAr         string     `json:"descriptionAr"`
	TermsConditionsEn     string     `json:"termsConditionsEn"`
	TermsConditionsAr     string     `json:"termsConditionsAr"`
	BenefitDiscountAmount float64    `json:"benefitDiscountAmount"`
	NoExpiryDate          bool       `json:"noExpiryDate"`
	Schedules             []Schedule `json:"schedules"`
	IsBranchSpecific      bool       `json:"isBranchSpecific"`
	Branches              []Branch   `json:"branches"`
	Brands                []int      `json:"brands"`
	ItemIds               []int      `json:"itemIds"`
}

type ChoiceCategory struct {
	ID                string   `json:"id"`
	EnglishName       string   `json:"englishName"`
	ArabicName        string   `json:"arabicName"`
	MinimumSelections int      `json:"minimumSelections"`
	MaximumSelections int      `json:"maximumSelections"`
	SortOrder         int      `json:"sortOrder,omitempty"`
	Choices           []Choice `json:"choices"`
	ProductIDs        []string `json:"-"`
}

type Choice struct {
	ID                   string               `json:"id"`
	EnglishName          string               `json:"englishName"`
	ArabicName           string               `json:"arabicName"`
	Price                float64              `json:"price"`
	IsAvailable          bool                 `json:"isAvailable"`
	SortOrder            int                  `json:"sortOrder,omitempty"`
	BranchesAvailability []BranchAvailability `json:"branchesAvailability,omitempty"`
	ImageURL             string               `json:"imageURL,omitempty"`
	Thumbnail            string               `json:"thumbnail,omitempty"`
	ChoiceCategories     []ChoiceCategory     `json:"choiceCategories,omitempty"`
	AttributeGroupIDs    []string             `json:"-"`
}

type Item struct {
	ID                   string               `json:"id"`
	EnglishName          string               `json:"englishName"`
	ArabicName           string               `json:"arabicName"`
	EnglishDescription   string               `json:"englishDescription"`
	ArabicDescription    string               `json:"arabicDescription"`
	Price                float64              `json:"price"`
	IsAvailable          bool                 `json:"isAvailable"`
	SortOrder            int                  `json:"sortOrder,omitempty"`
	AvailableFrom        string               `json:"availableFrom"`
	AvailableTo          string               `json:"availableTo"`
	AvailableDays        string               `json:"availableDays"`
	BranchesAvailability []BranchAvailability `json:"branchesAvailability,omitempty"`
	ImageURL             string               `json:"ImageURL"`
	//ItemDiscount         ItemDiscount         `json:"itemDiscount,omitempty"`
	Thumbnail        string           `json:"Thumbnail"`
	ChoiceCategories []ChoiceCategory `json:"choiceCategories"`
	CategoryIDs      []string         `json:"-"`
}

type Category struct {
	ID          string `json:"id"`
	EnglishName string `json:"englishName"`
	ArabicName  string `json:"arabicName"`
	SortOrder   int    `json:"sortOrder,omitempty"`
	Items       []Item `json:"items"`
}

type Menu struct {
	PaperboyURL string     `json:"paperboyUrl"`
	CallbackURL string     `json:"callbackUrl"`
	ScheduledOn string     `json:"scheduledOn"`
	Categories  []Category `json:"categories"`
}
type UpdateItemsAvailabilityRequest struct {
	PaperboyUrl  string         `json:"paperboyUrl,omitempty"`
	CallbackUrl  string         `json:"callbackUrl,omitempty"`
	ScheduledOn  string         `json:"scheduledOn"`
	Availability []Availability `json:"availability"`
	RestaurantID string         `json:"-"`
	RequestID    string         `json:"-"`
}

type Availability struct {
	BranchId string         `json:"branchId"`
	Items    []ItemStoplist `json:"items"`
}

type ItemStoplist struct {
	CategoryId string `json:"categoryId,omitempty"`
	ItemId     string `json:"itemId"`
	Status     int    `json:"status"`
}
