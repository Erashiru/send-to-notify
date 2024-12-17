package selector

type StoreGroup struct {
	ID       string
	Name     string
	StoreIDs []string
	Country  string
	Category string
	Status   string

	Countries  []string
	Categories []string
	Statuses   []string
	DomainName string

	Pagination
	Sorting
}

func NewEmptyStoreGroupSearch() StoreGroup {
	return StoreGroup{}
}

func (s StoreGroup) HasStoreIDs() bool {
	return len(s.StoreIDs) > 0
}

func (s StoreGroup) SetStoreIDs(storeIDs []string) StoreGroup {
	s.StoreIDs = append(s.StoreIDs, storeIDs...)
	return s
}

func (s StoreGroup) SetID(id string) StoreGroup {
	s.ID = id
	return s
}

func (s StoreGroup) HasID() bool {
	return s.ID != ""
}

func (s StoreGroup) SetName(name string) StoreGroup {
	s.Name = name
	return s
}

func (s StoreGroup) HasName() bool {
	return s.Name != ""
}

func (s StoreGroup) HasCountry() bool {
	return s.Country != ""
}

func (s StoreGroup) SetCountry(country string) StoreGroup {
	s.Country = country
	return s
}

func (s StoreGroup) HasCategory() bool {
	return s.Category != ""
}

func (s StoreGroup) SetCategory(category string) StoreGroup {
	s.Category = category
	return s
}

func (s StoreGroup) HasStatus() bool {
	return s.Status != ""
}

func (s StoreGroup) SetStatus(status string) StoreGroup {
	s.Status = status
	return s
}

func (s StoreGroup) HasDomainName() bool {
	return s.DomainName != ""
}

func (s StoreGroup) SetDomainName(domainName string) StoreGroup {
	s.DomainName = domainName
	return s
}

func (s StoreGroup) HasCountries() bool {
	return s.Countries != nil && len(s.Countries) != 0
}

func (s StoreGroup) SetCountries(countries []string) StoreGroup {
	s.Countries = append(s.Countries, countries...)
	return s
}

func (s StoreGroup) HasCategories() bool {
	return s.Categories != nil && len(s.Categories) != 0
}

func (s StoreGroup) SetCategories(categories []string) StoreGroup {
	s.Categories = append(s.Categories, categories...)
	return s
}

func (s StoreGroup) HasStatuses() bool {
	return s.Statuses != nil && len(s.Statuses) != 0
}

func (s StoreGroup) SetStatuses(statuses []string) StoreGroup {
	s.Statuses = append(s.Statuses, statuses...)
	return s
}

func (s StoreGroup) SetPage(page int64) StoreGroup {
	if page > 0 {
		s.Page = page - 1
	}
	return s
}

func (s StoreGroup) HasPage() bool {
	return s.Page != 0
}

func (s StoreGroup) SetLimit(limit int64) StoreGroup {
	if limit > 0 {
		s.Limit = limit
	}
	return s
}

func (s StoreGroup) SetSorting(param string, dir int8) StoreGroup {
	s.Sorting.Param = param
	s.Sorting.Direction = dir
	return s
}
