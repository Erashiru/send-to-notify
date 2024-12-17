package models

type ErrorSolution struct {
	ID                  string   `json:"id"  bson:"_id,omitempty"`
	Type                string   `json:"type" bson:"type"`
	Code                string   `json:"code" bson:"code"`
	BusinessName        string   `json:"business_name" bson:"business_name"`
	Reason              string   `json:"reason" bson:"reason"`
	Solution            string   `json:"solution" bson:"solution"`
	ContainsText        string   `json:"contains_text" bson:"contains_text"`
	PutProductOnStop    bool     `json:"put_product_on_stop" bson:"put_product_on_stop"`       //если true, во всех заказах где мы получаем эту ошибку, ставим продукты на стоп в агрегаторах
	StopListPosTypes    []string `json:"stop_list_pos_types" bson:"stop_list_pos_types"`       //когда ставим продукт на стоп, сначала проверяем, пос система ресторана входит в этот список
	RegexpToFindProduct string   `json:"regexp_to_find_product" bson:"regexp_to_find_product"` //чтобы найти продукт которую мы будем ставить на стоп, через этот regexp будем искать продукт айди в сообщение
	SendToTelegram      bool     `json:"send_to_telegram" bson:"send_to_telegram"`             //если true, отправляем в тг чат, где и какой продукт мы поставили на стоп
	Avoidable           bool     `json:"avoidable" bson:"avoidable"`                           //если true, тогда это ошибки на которые мы влияем, если false ошибки со стороны партнера
	IsTimeout           bool     `json:"is_timeout" bson:"is_timeout"`
}
