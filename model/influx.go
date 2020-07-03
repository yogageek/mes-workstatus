package model

//定義要放入influxdb的格式

type WoFields struct {
	WoFields []WoField
}

//如果要append自己的struct必須寫一個Add
func (ws *WoFields) Add(w WoField) []WoField {
	ws.WoFields = append(ws.WoFields, w)
	// ws.OrderIds = append(ws.OrderIds, w.OId)
	return ws.WoFields
}

//工單(每個工站一個)
type WoField struct {
	Seq                    int    // wo id 1
	WorkorderId            string //
	Process                string // 前置工站1
	Status                 int    // 0,1,2
	GoodQty                int
	NgQty                  int
	Qty                    int    // total qty
	ProcessTimePerPcs      int    // 此工項每pcs標準工時
	CurrentStatus          string // 完成數量字串 100/300
	StardardCompletionTime int    // qty*每pcs標準工時
}

type WoTag struct {
	Type           string
	WorkorderIdTag string //要加上tag不然時間相同會被蓋掉
	ManorderIdTag  string
	OrderIdTag     string
}

//定單
type OrderField struct {
	ProductId      string
	ProductName    string
	PlanProduction int //order qty
	OrderOEE       int //order complete percentage
}

type OrderTag struct {
	Type       string
	OrderIdTag string
}

//制定單
type ManorderField struct {
	PlanProduction int    //order qty
	OrderOEE       int    //order complete percentage
	DueDate        string //預計完成日
}

type ManorderTag struct {
	Type          string
	ManorderIdTag string
	OrderIdTag    string
}
