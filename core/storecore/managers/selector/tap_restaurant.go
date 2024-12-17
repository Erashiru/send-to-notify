package selector

type TapRestaurant struct {
	ID         string
	Name       string
	QRMenuLink string
	Tel        string
	Instagram  string
	Website    string
	Sorting
	Pagination
}

func EmptyTapRestaurant() TapRestaurant {
	return TapRestaurant{}
}
func TapRestaurantSearch() TapRestaurant {
	return TapRestaurant{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (tr TapRestaurant) SetID(id string) TapRestaurant {
	tr.ID = id
	return tr
}
func (tr TapRestaurant) HasID() bool {
	return tr.ID != ""
}

func (tr TapRestaurant) SetName(name string) TapRestaurant {
	tr.Name = name
	return tr
}
func (tr TapRestaurant) HasName() bool {
	return tr.Name != ""
}

func (tr TapRestaurant) SetQRMenuLink(qrMenuLink string) TapRestaurant {
	tr.QRMenuLink = qrMenuLink
	return tr
}
func (tr TapRestaurant) HasQRMenuLink() bool {
	return tr.QRMenuLink != ""
}

func (tr TapRestaurant) SetTel(tel string) TapRestaurant {
	tr.Tel = tel
	return tr
}
func (tr TapRestaurant) HasTel() bool {
	return tr.Tel != ""
}

func (tr TapRestaurant) SetInstagram(instagram string) TapRestaurant {
	tr.Instagram = instagram
	return tr
}
func (tr TapRestaurant) HasInstagram() bool {
	return tr.Instagram != ""
}

func (tr TapRestaurant) SetWebsite(website string) TapRestaurant {
	tr.Website = website
	return tr
}
func (tr TapRestaurant) HasWebsite() bool {
	return tr.Website != ""
}

func (tr TapRestaurant) SetPage(page int64) TapRestaurant {
	if page > 0 {
		tr.Pagination.Page = page - 1
	}
	return tr
}
func (tr TapRestaurant) SetLimit(limit int64) TapRestaurant {
	if limit > 0 {
		tr.Pagination.Limit = limit
	}
	return tr
}
func (tr TapRestaurant) SetSorting(key string, dir int8) TapRestaurant {
	tr.Sorting.Param = key
	tr.Sorting.Direction = dir
	return tr
}
