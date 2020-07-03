package logic

import (
	"mes-workstatus/model"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

//process time == time per pcs of this wo
func takeProcessTimeFromWo(m model.WorkOrders, stepid int) int {
	stepidstr := strconv.Itoa(stepid)
	// if stepid == "" {
	// 	return 0
	// }
	// trimmedStepid := strings.Split(stepid, "-")[1]
	for _, w := range m.Product.Route {
		// fmt.Println(" id: ", w.Id)
		if w.Id == stepidstr {
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

func takeTextFromWo(m model.WorkOrders, stepid int) string {
	stepidstr := strconv.Itoa(stepid)
	// if stepid == "" {
	// 	return ""
	// }
	// trimmedStepid := strings.Split(stepid, "-")[1]
	for _, w := range m.Product.Route {
		// fmt.Println(" id: ", w.Id)
		if w.Id == stepidstr {
			return w.Text
		}
	}
	glog.Error("text not found ")
	return ""
}

//depre
func takeOIDFromWo(stepid string) string {
	if stepid == "" {
		return ""
	}
	oid := strings.Split(stepid, "-")[0]
	return oid
}

//depre
func takeWoIDFromWo(stepid string) string {
	if stepid == "" {
		return ""
	}
	woid := strings.Split(stepid, "-")[1]
	return woid
}
