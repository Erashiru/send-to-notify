package selector

type Store struct {
	ID                      string
	Name                    string
	Currency                string
	LanguageCode            string
	StoreGroupId            string
	Usernames               []string
	Street                  string
	City                    string
	Timezone                string
	UtcOffset               float64
	ClientSecret            string
	DeliveryService         string
	ExternalStoreID         string
	ExternalDeliveryService string
	PosType                 string
	Hash                    string
	PosOrganizationID       string
	Token                   string
	AggregatorMenuIDs       []string
	IsActiveMenu            *bool
	IDs                     []string
	Express24StoreId        []string
	PosterAccountNumber     string
	IsVirtualStore          *bool
	TalabatRemoteBranchId   string
	ScheduledStatusChange   bool
	YarosStoreId            string
	IsChildStore            *bool
	DeferSubmission         *bool
	OrderAutoClose          *bool
	Status                  string
}

func NewEmptyStoreSearch() Store {
	return Store{}
}

func (s Store) SetHasVirtualStore(has *bool) Store {
	s.IsVirtualStore = has
	return s
}

func (s Store) HasVirtualStore() bool {
	return s.IsVirtualStore != nil
}

func (s Store) HasExpress24StoreId() bool {
	return len(s.Express24StoreId) > 0
}

func (s Store) SetExpress24StoreId(express24StoreId []string) Store {
	s.Express24StoreId = express24StoreId
	return s
}

func (s Store) SetTalabatRemoteBranchId(remoteId string) Store {
	s.TalabatRemoteBranchId = remoteId
	return s
}

func (s Store) SetYarosStoreId(storeId string) Store {
	s.YarosStoreId = storeId
	return s
}

func (s Store) HasYarosStoreId() bool {
	return s.YarosStoreId != ""
}
func (s Store) SetStreet(street string) Store {
	s.Street = street
	return s
}

func (s Store) HasStreet() bool {
	return s.Street != ""
}

func (s Store) SetCity(city string) Store {
	s.City = city
	return s
}

func (s Store) HasCity() bool {
	return s.City != ""
}

func (s Store) SetTimezone(timezone string) Store {
	s.Timezone = timezone
	return s
}

func (s Store) HasTimezone() bool {
	return s.Timezone != ""
}

func (s Store) SetUtcOffset(utcOffset float64) Store {
	s.UtcOffset = utcOffset
	return s
}

func (s Store) HasUtcOffset() bool {
	return s.UtcOffset != 0
}

func (s Store) SetUsernames(usernames []string) Store {
	s.Usernames = usernames
	return s
}

func (s Store) HasUsernames() bool {
	return len(s.Usernames) > 0
}

func (s Store) SetName(name string) Store {
	s.Name = name
	return s
}

func (s Store) HasName() bool {
	return s.Name != ""
}

func (s Store) HasTalabatRemoreBranchId() bool {
	return s.TalabatRemoteBranchId != ""
}

func (s Store) SetCurrency(currency string) Store {
	s.Currency = currency
	return s
}

func (s Store) HasCurrency() bool {
	return s.Currency != ""
}

func (s Store) SetLanguageCode(languageCode string) Store {
	s.LanguageCode = languageCode
	return s
}

func (s Store) HasLanguageCode() bool {
	return s.LanguageCode != ""
}

func (s Store) SetStoreGroupId(storeGroupId string) Store {
	s.StoreGroupId = storeGroupId
	return s
}

func (s Store) HasStoreGroupId() bool {
	return s.StoreGroupId != ""
}

func (s Store) ActiveMenu() bool {
	if s.IsActiveMenu != nil && *s.IsActiveMenu {
		return true
	}
	return false
}

func (s Store) SetIsActiveMenu(isActive *bool) Store {
	s.IsActiveMenu = isActive
	return s
}

func (s Store) HasIsActiveMenu() bool {
	return s.IsActiveMenu != nil
}

func (s Store) AggregatorMenuID() string {
	if len(s.AggregatorMenuIDs) != 1 {
		return ""
	}
	return s.AggregatorMenuIDs[0]
}

func (s Store) HasAggregatorMenuID() bool {
	if len(s.AggregatorMenuID()) > 1 || len(s.AggregatorMenuIDs) == 0 {
		return false
	}
	return s.AggregatorMenuIDs[0] != ""
}

func (s Store) HasAggregatorMenuIDs() bool {
	return len(s.AggregatorMenuIDs) > 0
}

func (s Store) SetAggregatorMenuID(menuID string) Store {
	s.AggregatorMenuIDs = []string{menuID}
	return s
}

func (s Store) SetAggregatorMenuIDs(menuIDs []string) Store {
	s.AggregatorMenuIDs = menuIDs
	return s
}

func (s Store) SetExternalDeliveryService(service string) Store {
	s.ExternalDeliveryService = service
	return s
}

func (s Store) HasExternalDeliveryService() bool {
	return s.ExternalDeliveryService != ""
}

func (s Store) SetToken(token string) Store {
	s.Token = token
	return s
}

func (s Store) HasToken() bool {
	return s.Token != ""
}

func (s Store) SetID(id string) Store {
	s.ID = id
	return s
}

func (s Store) HasID() bool {
	return s.ID != ""
}

func (s Store) SetClientSecret(secret string) Store {
	s.ClientSecret = secret
	return s
}

func (s Store) HasClientSecret() bool {
	return s.ClientSecret != ""
}

func (s Store) SetHash(value string) Store {
	s.Hash = value
	return s
}
func (s Store) HasHash() bool {
	return s.Hash != ""
}

func (s Store) SetDeliveryService(deliveryService string) Store {
	s.DeliveryService = deliveryService
	return s
}

func (s Store) HasDeliveryService() bool {
	return s.DeliveryService != ""
}

func (s Store) SetExternalStoreID(id string) Store {
	s.ExternalStoreID = id
	return s
}

func (s Store) HasExternalStoreID() bool {
	return s.ExternalStoreID != ""
}

func (s Store) SetPosType(posType string) Store {
	s.PosType = posType
	return s
}

func (s Store) HasPosType() bool {
	return s.PosType != ""
}

func (s Store) HasPosOrganizationID() bool {
	return s.PosOrganizationID != ""
}

func (s Store) SetPosOrganizationID(orgID string) Store {
	s.PosOrganizationID = orgID
	return s
}

func (s Store) SetStoreIDs(storeIDs []string) Store {
	s.IDs = storeIDs
	return s
}

func (s Store) HasStoreIDs() bool {
	return s.IDs != nil && len(s.IDs) != 0
}

func (s Store) HasAccountNumber() bool {
	return s.PosterAccountNumber != ""
}

func (s Store) SetPosterAccountNumber(accountNumber string) Store {
	s.PosterAccountNumber = accountNumber
	return s
}

func (s Store) SetScheduledStatusChange(value bool) Store {
	s.ScheduledStatusChange = value
	return s
}

func (s Store) HasScheduledStatusChange() bool {
	return s.ScheduledStatusChange
}

func (s Store) SetIsChildStore(isChildStore *bool) Store {
	s.IsChildStore = isChildStore
	return s
}

func (s Store) HasIsChildStore() bool {
	return s.IsChildStore != nil
}

func (s Store) SetDeferSubmission(isDeferSubmission *bool) Store {
	s.DeferSubmission = isDeferSubmission
	return s
}

func (s Store) HasDeferSubmission() bool {
	return s.DeferSubmission != nil
}

func (s Store) SetOrderAutoClose(autoClose *bool) Store {
	s.OrderAutoClose = autoClose
	return s
}

func (s Store) HasOrderAutoClose() bool {
	return s.OrderAutoClose != nil
}
