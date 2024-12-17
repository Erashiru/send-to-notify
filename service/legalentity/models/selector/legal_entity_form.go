package selector

type LegalEntityForm struct {
	Name             string    `json:"name"`
	BIN              string    `json:"bin"`
	KNP              string    `json:"knp"`
	PaymentType      string    `json:"payment_type"`
	LinkedAccManager string    `json:"linked_acc_manager"`
	SalesID          string    `json:"sales_id"`
	Contacts         []Contact `json:"contacts"`
	SalesComment     string    `json:"sales_comment"`
	StoreIds         []string  `json:"store_ids"`
	PaymentCycle     int       `json:"payment_cycle"`
}

type Contact struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Comment  string `json:"comment,omitempty"`
}

func (le LegalEntityForm) HasName() bool {
	return le.Name != ""
}

func (le LegalEntityForm) HasBIN() bool {
	return le.BIN != ""
}

func (le LegalEntityForm) HasKNP() bool {
	return le.KNP != ""
}

func (le LegalEntityForm) HasPaymentType() bool {
	return le.PaymentType != ""
}

func (le LegalEntityForm) HasLinkedAccManager() bool {
	return le.LinkedAccManager != ""
}

func (le LegalEntityForm) HasSalesID() bool {
	return le.SalesID != ""
}

func (le LegalEntityForm) HasContacts() bool {
	var ok bool
	for _, c := range le.Contacts {
		if c.Email == "" || c.FullName == "" || c.Phone == "" || c.Position == "" {
			ok = false
		} else {
			ok = true
		}
	}
	return ok
}

func (le LegalEntityForm) HasSalesComment() bool {
	return le.SalesComment != ""
}

func (le LegalEntityForm) HasStoreIds() bool {
	return le.StoreIds != nil
}

func (le LegalEntityForm) HasPaymentCycle() bool {
	return le.PaymentCycle != 0
}
