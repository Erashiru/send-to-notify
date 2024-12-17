package selector

type Menu struct {
	ID        string
	Name      string
	SectionID string

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

func (m Menu) HasMenuName() bool {
	return m.Name != ""
}

func (m Menu) HasSectionID() bool {
	return m.SectionID != ""
}

func (m Menu) SetMenuID(id string) Menu {
	m.ID = id
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
