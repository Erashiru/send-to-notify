package selector

type Menu struct {
	ID      string
	StoreID string
	Name    string
	Token   string

	ProductExtID string // productID is a
	SectionID    string
	GroupID      string

	IsProductAvailable *bool

	Sorting
	Pagination
}

func EmptyMenuSearch() Menu {
	return Menu{}
}

func MenuSearch() Menu {
	return Menu{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (m Menu) MenuID() string {
	return m.ID
}

func (m Menu) HasMenuID() bool {
	return m.ID != ""
}

func (m Menu) HasStoreID() bool {
	return m.StoreID != ""
}

func (m Menu) HasProductExtID() bool {
	return m.ProductExtID != ""
}

func (m Menu) HasMenuName() bool {
	return m.Name != ""
}

func (m Menu) HasSectionID() bool {
	return m.SectionID != ""
}

func (m Menu) HasGroupID() bool {
	return m.GroupID != ""
}

func (m Menu) HasToken() bool {
	return m.Token != ""
}

func (m Menu) ProductAvailable() bool {
	if m.IsProductAvailable != nil && *m.IsProductAvailable {
		return true
	}
	return false
}

func (m Menu) HasProductIsAvailable() bool {
	return m.IsProductAvailable != nil
}

func (m Menu) SetMenuID(id string) Menu {
	m.ID = id
	return m
}

func (m Menu) SetStoreID(id string) Menu {
	m.StoreID = id
	return m
}

func (m Menu) SetToken(id string) Menu {
	m.Token = id
	return m
}

func (m Menu) SetMenuName(name string) Menu {
	m.Name = name
	return m
}

func (m Menu) SetSectionID(sectionID string) Menu {
	m.SectionID = sectionID
	return m
}

func (m Menu) SetGroupID(groupID string) Menu {
	m.GroupID = groupID
	return m
}

func (m Menu) SetProductExtID(productID string) Menu {
	m.ProductExtID = productID
	return m
}

func (m Menu) SetProductIsAvailable(isAvailable *bool) Menu {
	m.IsProductAvailable = isAvailable
	return m
}

func (m Menu) SetPage(page int64) Menu {
	if page > 0 {
		m.Pagination.Page = page - 1
	}
	return m
}

func (m Menu) SetLimit(limit int64) Menu {
	if limit > 0 {
		m.Pagination.Limit = limit
	}
	return m
}

func (m Menu) SetSorting(key string, dir int8) Menu {
	m.Sorting.Param = key
	m.Sorting.Direction = dir
	return m
}
