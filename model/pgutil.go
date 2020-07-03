package model

import (
	"strconv"

	"github.com/golang/glog"
)

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
