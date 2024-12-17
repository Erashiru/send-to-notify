package models

import "time"

type Config struct {
	BaseURL       string `json:"base_url"`
	UserID        string `json:"user_id"`
	Secret        string `json:"token"`
	GroupID       string `json:"group_id"`
	ResponsibleID string `json:"responsible_id"`

	Duties string `json:"duties"`
}

type Task struct {
	ID                 string `json:"ID"`
	ParentID           string `json:"PARENT_ID"`
	StageID            string `json:"STAGE_ID"`
	GroupID            string `json:"GROUP_ID"`
	ResponsibleID      string `json:"RESPONSIBLE_ID"`
	GUID               string `json:"GUID"`
	ForumID            string `json:"FORUM_ID"`
	ForumTopicID       string `json:"FORUM_TOPIC_ID"`
	ForkedByTemplateID int    `json:"FORKED_BY_TEMPLATE_ID,string"`
	SiteID             string `json:"SITE_ID"`

	Title               string `json:"TITLE"`
	Description         string `json:"DESCRIPTION"`
	Deadline            string `json:"DEADLINE"`
	StartDatePlan       string `json:"START_DATE_PLAN"`
	EndDatePlan         string `json:"END_DATE_PLAN"`
	AllowChangeDeadline string `json:"ALLOW_CHANGE_DEADLINE"`
	TaskControl         string `json:"TASK_CONTROL"`
	TimeEstimate        string `json:"TIME_ESTIMATE"`
	CreatedBy           string `json:"CREATED_BY"`
	DescriptionInBbcode string `json:"DESCRIPTION_IN_BBCODE"`
	DeclineReason       string `json:"DECLINE_REASON"`

	RealStatus string `json:"REAL_STATUS"`
	Status     string `json:"STATUS"`

	ResponsibleName       string `json:"RESPONSIBLE_NAME"`
	ResponsibleLastName   string `json:"RESPONSIBLE_LAST_NAME"`
	ResponsibleSecondName string `json:"RESPONSIBLE_SECOND_NAME"`

	DateStart    string      `json:"DATE_START"`
	DurationPlan interface{} `json:"DURATION_PLAN"`
	DurationFact string      `json:"DURATION_FACT"`
	DurationType string      `json:"DURATION_TYPE"`

	Priority string `json:"PRIORITY"`
	Mark     string `json:"MARK"`

	CreatedDate       string `json:"CREATED_DATE"`
	ChangedBy         string `json:"CHANGED_BY"`
	ChangedDate       string `json:"CHANGED_DATE"`
	StatusChangedBy   string `json:"STATUS_CHANGED_BY"`
	StatusChangedDate string `json:"STATUS_CHANGED_DATE"`
	ClosedBy          string `json:"CLOSED_BY"`
	ClosedDate        string `json:"CLOSED_DATE"`
	ActivityDate      string `json:"ACTIVITY_DATE"`
	ViewedDate        string `json:"VIEWED_DATE"`
	TimeSpentInLogs   string `json:"TIME_SPENT_IN_LOGS"`

	Favorite          string `json:"FAVORITE"`
	AllowTimeTracking string `json:"ALLOW_TIME_TRACKING"`
	MatchWorkTime     string `json:"MATCH_WORK_TIME"`
	AddInReport       string `json:"ADD_IN_REPORT"`

	CommentsCount  string         `json:"COMMENTS_COUNT"`
	Subordinate    string         `json:"SUBORDINATE"`
	Multitask      string         `json:"MULTITASK"`
	AllowedActions AllowedActions `json:"ALLOWED_ACTIONS"`

	NotViewed            string        `json:"notViewed"`
	Replicate            string        `json:"replicate"`
	XMLID                int           `json:"xmlId,string"`
	ServiceCommentsCount string        `json:"serviceCommentsCount"`
	ExchangeModified     interface{}   `json:"exchangeModified"`
	ExchangeID           int           `json:"exchangeId,string"`
	OutlookVersion       string        `json:"outlookVersion"`
	Sorting              interface{}   `json:"sorting"`
	IsMuted              string        `json:"isMuted"`
	IsPinned             string        `json:"isPinned"`
	IsPinnedInGroup      string        `json:"isPinnedInGroup"`
	Auditors             []interface{} `json:"auditors"`
	Accomplices          []interface{} `json:"accomplices"`
	NewCommentsCount     int           `json:"newCommentsCount"`
	Group                interface{}   `json:"group"` // in all tasks this field is like `GroupData struct`, but in one and simgle it is empty [], so -- interface{}
	Creator              UserData      `json:"creator"`
	Responsible          UserData      `json:"responsible"`
	SubStatus            int           `json:"subStatus,string"`
}

type AllowedActions struct {
	Accept                bool `json:"ACTION_ACCEPT"`
	Decline               bool `json:"ACTION_DECLINE"`
	Complete              bool `json:"ACTION_COMPLETE"`
	Approve               bool `json:"ACTION_APPROVE"`
	Disapprove            bool `json:"ACTION_DISAPPROVE"`
	Start                 bool `json:"ACTION_START"`
	Pause                 bool `json:"ACTION_PAUSE"`
	Delegate              bool `json:"ACTION_DELEGATE"`
	Remove                bool `json:"ACTION_REMOVE"`
	Edit                  bool `json:"ACTION_EDIT"`
	Defer                 bool `json:"ACTION_DEFER"`
	Renew                 bool `json:"ACTION_RENEW"`
	Create                bool `json:"ACTION_CREATE"`
	ChangeDeadline        bool `json:"ACTION_CHANGE_DEADLINE"`
	ChecklistAddItems     bool `json:"ACTION_CHECKLIST_ADD_ITEMS"`
	ChecklistReorderItems bool `json:"ACTION_CHECKLIST_REORDER_ITEMS"`
	ChangeDirector        bool `json:"ACTION_CHANGE_DIRECTOR"`
	ElapsedTimeAdd        bool `json:"ACTION_ELAPSED_TIME_ADD"`
	StartTimeTracking     bool `json:"ACTION_START_TIME_TRACKING"`
	AddFavourite          bool `json:"ACTION_ADD_FAVORITE"`
	DeleteFavourite       bool `json:"ACTION_DELETE_FAVORITE"`
	Rate                  bool `json:"ACTION_RATE"`
}

type GroupData struct {
	ID             int           `json:"id,string"`
	Name           string        `json:"name"`
	Opened         bool          `json:"opened"`
	MembersCount   int           `json:"membersCount"`
	Image          string        `json:"image"`
	AdditionalData []interface{} `json:"additionalData"`
}

type UserData struct {
	ID   int    `json:"id,string"`
	Name string `json:"name"`
	Link string `json:"link"`
	Icon string `json:"icon"`
}

type TasksListResponse struct {
	Result []Task `json:"result"`
	Next   int    `json:"next"`
	Total  int    `json:"total"`
	Time   Time   `json:"time"`
}

type Time struct {
	Start            float64   `json:"start"`
	FinishAt         float64   `json:"finish"`
	Duration         float64   `json:"duration"`
	Processing       float64   `json:"processing"`
	StartAt          time.Time `json:"date_start"`
	EndAt            time.Time `json:"date_finish"`
	OperatingResetAt int       `json:"operating_reset_at"`
	Operating        float64   `json:"operating"`
}

type AddTask struct {
	Title         string     `json:"TITLE"`
	Description   string     `json:"DESCRIPTION"`
	GroupID       string     `json:"GROUP_ID"`
	ResponsibleID int        `json:"RESPONSIBLE_ID"`
	Deadline      *time.Time `json:"DEADLINE,omitempty"`
}

type AddTaskResponse struct {
	Result int  `json:"result"`
	Time   Time `json:"time"`
}

type GetLeadResponse struct {
	Result Lead `json:"result"`
	Time   Time `json:"time"`
}

type Lead struct {
	ID                      string      `json:"ID"`
	Title                   string      `json:"TITLE"`
	Honorific               interface{} `json:"HONORIFIC"`
	Name                    string      `json:"NAME"`
	SecondName              interface{} `json:"SECOND_NAME"`
	LastName                interface{} `json:"LAST_NAME"`
	CompanyTitle            interface{} `json:"COMPANY_TITLE"`
	CompanyID               interface{} `json:"COMPANY_ID"`
	ContactID               interface{} `json:"CONTACT_ID"`
	IsReturnCustomer        string      `json:"IS_RETURN_CUSTOMER"`
	Birthdate               string      `json:"BIRTHDATE"`
	SourceID                string      `json:"SOURCE_ID"`
	SourceDescription       interface{} `json:"SOURCE_DESCRIPTION"`
	StatusID                string      `json:"STATUS_ID"`
	StatusDescription       interface{} `json:"STATUS_DESCRIPTION"`
	Post                    interface{} `json:"POST"`
	Comments                interface{} `json:"COMMENTS"`
	CurrencyID              string      `json:"CURRENCY_ID"`
	Opportunity             string      `json:"OPPORTUNITY"`
	IsManualOpportunity     string      `json:"IS_MANUAL_OPPORTUNITY"`
	HasPhone                string      `json:"HAS_PHONE"`
	HasEmail                string      `json:"HAS_EMAIL"`
	HadImol                 string      `json:"HAS_IMOL"`
	AssignedByID            string      `json:"ASSIGNED_BY_ID"`
	CreatedByID             string      `json:"CREATED_BY_ID"`
	ModifyByID              string      `json:"MODIFY_BY_ID"`
	DateCreate              time.Time   `json:"DATE_CREATE"`
	DateModify              time.Time   `json:"DATE_MODIFY"`
	DateClosed              string      `json:"DATE_CLOSED"`
	StatusSemanticID        string      `json:"STATUS_SEMANTIC_ID"`
	Opened                  string      `json:"OPENED"`
	OriginatorID            interface{} `json:"ORIGINATOR_ID"`
	OriginID                interface{} `json:"ORIGIN_ID"`
	MoveByID                string      `json:"MOVED_BY_ID"`
	MoveTime                time.Time   `json:"MOVED_TIME"`
	Address                 interface{} `json:"ADDRESS"`
	Address2                interface{} `json:"ADDRESS_2"`
	AddressCity             interface{} `json:"ADDRESS_CITY"`
	AddressPortalCode       interface{} `json:"ADDRESS_POSTAL_CODE"`
	AddressRegion           interface{} `json:"ADDRESS_REGION"`
	AddressProvince         interface{} `json:"ADDRESS_PROVINCE"`
	AddressCountry          interface{} `json:"ADDRESS_COUNTRY"`
	AddressCountryCode      interface{} `json:"ADDRESS_COUNTRY_CODE"`
	AddressLocAddrID        interface{} `json:"ADDRESS_LOC_ADDR_ID"`
	UtmSource               interface{} `json:"UTM_SOURCE"`
	UtcMedium               interface{} `json:"UTM_MEDIUM"`
	UtmCampaign             interface{} `json:"UTM_CAMPAIGN"`
	UtmContent              interface{} `json:"UTM_CONTENT"`
	UtmTerm                 interface{} `json:"UTM_TERM"`
	LastActivityBy          string      `json:"LAST_ACTIVITY_BY"`
	LastActivityTime        time.Time   `json:"LAST_ACTIVITY_TIME"`
	UfCrmLead1718344223283  string      `json:"UF_CRM_LEAD_1718344223283"`
	UfCrmLead1718696937920  string      `json:"UF_CRM_LEAD_1718696937920"`
	UfCrmInstagramWz        string      `json:"UF_CRM_INSTAGRAM_WZ"`
	UfCrmVkWz               string      `json:"UF_CRM_VK_WZ"`
	UfCrmAvitoWz            string      `json:"UF_CRM_AVITO_WZ"`
	UfCrmTelegramusernameWz string      `json:"UF_CRM_TELEGRAMUSERNAME_WZ"`
	UfCrmTelegramidWz       string      `json:"UF_CRM_TELEGRAMID_WZ"`
	UfCrmLead1719560953439  string      `json:"UF_CRM_LEAD_1719560953439"`
	Phone                   []LeadInfo  `json:"PHONE"`
	Im                      []LeadInfo  `json:"IM"`
	Link                    []LeadInfo  `json:"LINK"`
}

type LeadInfo struct {
	ID        string `json:"ID"`
	ValueType string `json:"VALUE_TYPE"`
	Value     string `json:"VALUE"`
	TypeID    string `json:"TYPE_ID"`
}
