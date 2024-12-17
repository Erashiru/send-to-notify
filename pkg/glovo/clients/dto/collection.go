package dto

type SuperCollection struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Position    int      `json:"position"`
	ImageUrl    string   `json:"image_url"`
	Collections []string `json:"collections"`
}

type Collection struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name"`
	Position int       `json:"position"`
	ImageUrl string    `json:"image_url,omitempty"`
	Sections []Section `json:"sections"`
	Schedule *Schedule `json:"schedule,omitempty"`
}

type Schedule struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Availabilities []Availability `json:"availabilities"`
}

type Availability struct {
	Day       string     `json:"day,omitempty"`
	TimeSlots []TimeSlot `json:"time_slots,omitempty"`
}

type TimeSlot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Section struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Products []string `json:"products"`
}
