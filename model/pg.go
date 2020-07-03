package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/golang/glog"
)

//工項單(每個工站一個)
type WorkOrders struct {
	WorkorderId string    `json:"workorder_id,omitempty"`
	State       State     `json:"state,omitempty"` //良品 不良品
	Product     Product   `json:"product,omitempty"`
	OrderInfo   OrderInfo `json:"order_info,omitempty"` //取代manorder
	Time        time.Time
	WoId        int `json:"wo_id,omitempty"` //工單號 //deprecated
	// Manorder    Manorder `json:"manorder,omitempty"` //訂單數量 //deprecated
}

func (w *WorkOrders) CalGoodSpendedTime() int {
	goodqty := w.State.AccGood
	processtime := w.TakeProcessTimeFromWo()
	return goodqty * processtime
}

//process time == time per pcs of this workorder
func (w *WorkOrders) TakeProcessTimeFromWo() int {
	stepid := strconv.Itoa(w.OrderInfo.StepId)

	for _, w := range w.Product.Route {
		// fmt.Println(" id: ", w.Id)
		if w.Id == stepid {
			for _, t := range w.Lines {
				// fmt.Println("id: ", w.Id, " processtime: ", t.ProcessTime)
				t, err := strconv.Atoi(t.ProcessTime)
				if err != nil {
					glog.Error(err)
				}
				return t
			}
		}
	}
	glog.Error("time per pcs of this workorder not found ")
	return 0
}
func (w *WorkOrders) TakeTextFromWo() string {
	stepid := strconv.Itoa(w.OrderInfo.StepId)

	for _, w := range w.Product.Route {
		// fmt.Println(" id: ", w.Id)
		if w.Id == stepid {
			return w.Text
		}
	}
	glog.Error("text not found ")
	return ""
}
func (w *WorkOrders) TakeTextSeq() string {
	return w.TakeTextFromWo() + strconv.Itoa(w.OrderInfo.StepId)
}

type State struct {
	Status  string `json:"status,omitempty"`
	AccNg   int    `json:"acc_ng,omitempty"`
	AccGood int    `json:"acc_good,omitempty"`
}

// depre
// The Attrs struct represents the data in the JSON/JSONB column. We can use
// struct tags to control how each field is encoded.
// type Manorder struct {
// 	Qty     int `json:"qty,omitempty"`
// 	Product struct {
// 		Workflow []struct {
// 			Id       string `json:"id,omitempty"`
// 			Text     string `json:"text,omitempty"`
// 			Machines []struct {
// 				ProcessTime string `json:"process_time,omitempty"`
// 			} `json:"machines,omitempty"`
// 		} `json:"workflow,omitempty"`
// 	} `json:"product,omitempty"`
// }

//部分為舊版manorder
type Product struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Spec    string `json:"spec,omitempty"`
	RouteId string `json:"route_id,omitempty"`
	Route   []struct {
		Id    string `json:"id,omitempty"`
		Text  string `json:"text,omitempty"`
		Lines []struct {
			ProcessTime string `json:"process_time,omitempty"`
		}
	} `json:"route,omitempty"`
}

//取代manorder
type OrderInfo struct {
	OrderId     string `json:"order_id,omitempty"`
	ManorderId  string `json:"manorder_id,omitempty"`
	WorkOrderId string `json:"workOrder_id,omitempty"`
	StepId      int    `json:"step_id,omitempty"`
	ProductId   string `json:"product_id,omitempty"`
	Qty         int    `json:"qty,omitempty"`
}

// special calcuate for select order qty
type OrderJoinWorkorder struct {
	OrderId     string
	TotalQty    int //加總
	Qty         int
	AccNg       int `json:"acc_ng,omitempty"`   //加總
	AccGood     int `json:"acc_good,omitempty"` //加總
	Time        time.Time
	ProductId   string
	ProductName string
}

type OrderJoinManorder struct {
	OrderId      string
	ManorderId   string
	TotalQty     int //加總
	Qty          int
	AccGood      int       //加總
	AccNg        int       //加總
	WorkedTime   int       //GoodQty已花費時間
	RequiredTime int       //Qty共需花費時間
	DueDate      time.Time //制令單預計完成日
	Time         time.Time
}

//如果要讓pg json格式資料轉go struct, 需要增加func(指定struct)的 Scan & Value方法

func (a State) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *State) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (a OrderInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *OrderInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (a Product) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Product) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
