package save_order_request

type RK7Query struct {
	RK7CMD RK7CMD `xml:"RK7CMD"`
}

type LicenseInfo struct {
	Anchor          string           `xml:"anchor,attr,omitempty"`
	LicenseToken    string           `xml:"licenseToken,attr,omitempty"`
	LicenseInstance *LicenseInstance `xml:"LicenseInstance,omitempty"`
}

type LicenseInstance struct {
	Guid      string `xml:"guid,attr,omitempty"`
	SeqNumber string `xml:"seqNumber,attr,omitempty"`
}

type Order struct {
	Visit      string `xml:"visit,attr"`
	OrderIdent string `xml:"orderIdent,attr"`
}

type Session struct {
	Station Station `xml:"Station"`
	Dish    []Dish  `xml:"Dish"`
	Prepay  *Prepay `xml:"Prepay,omitempty"`
}

type Prepay struct {
	ID       string  `xml:"id,attr,omitempty"`
	Amount   string  `xml:"amount,attr,omitempty"`
	Promised string  `xml:"Promised,attr,omitempty"`
	Reason   *Reason `xml:"Reason,omitempty"`
}

type Reason struct {
	ID string `xml:"id,attr,omitempty"`
}

type Station struct {
	ID string `xml:"id,attr"`
}

type Dish struct {
	ID       string `xml:"id,attr"`
	Quantity string `xml:"quantity,attr"`
	Price    string `xml:"price,attr"`
	Modi     []Modi `xml:"Modi"`
}

type Modi struct {
	ID    string `xml:"id,attr"`
	Count string `xml:"count,attr"`
	Price string `xml:"price,attr"`
}

type Dishes []Dish

func (d Dishes) ConvertToCorrectFormat() Dishes {
	var (
		dishes = make([]Dish, 0, len(d))
	)

	for _, dish := range d {
		modificators := make([]Modi, len(dish.Modi))

		for i, modificator := range dish.Modi {
			modificator.Price += "00"
			modificators[i] = modificator
		}

		dishes = append(dishes, Dish{
			ID:       dish.ID,
			Price:    dish.Price + "00",
			Quantity: dish.Quantity + "000",
			Modi:     modificators,
		})
	}

	return dishes
}

type RK7CMD struct {
	CMD         string       `xml:"CMD,attr"`
	LicenseInfo *LicenseInfo `xml:"LicenseInfo,omitempty"`
	Order       Order        `xml:"Order"`
	Session     Session      `xml:"Session"`
}
