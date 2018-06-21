package main

import (
	"batch-booking/models"
	"encoding/json"
	"fmt"

	"github.com/otwdev/galaxylib"
	"github.com/otwdev/ibepluslib/ibeplus"
	ibelib "github.com/otwdev/ibepluslib/models"
)

func main() {

	galaxylib.DefaultGalaxyConfig.InitConfig()

	galaxylib.DefaultGalaxyLog.ConfigLogger()

	officeNumber := galaxylib.GalaxyCfgFile.MustValue("ibe", "officeNumber")

	xlsReader := models.NewExcelReader(officeNumber, "./dbfile/1710内部人员信息-反馈.xlsx")
	xlsReader.Read(func(order *ibelib.OrderInfo) {
		buf, _ := json.Marshal(order)
		fmt.Println(string(buf))

		order.OfficeNumber = officeNumber
		//order.
		rq := ibeplus.NewPolicyDRQ(order) //ibeplus.NewFarePriceRequest(order) //ibeplus.NewAirAvail(order)
		root := rq.PolicyD()              //rq.AirAvailRQ()

		pnr := order.PnrInofs[0]

		err := root.GetPNRPrice(pnr)

		if err != nil {
			galaxylib.GalaxyLogger.Errorln(err)
			return
		}

		pnrbuf, _ := json.Marshal(pnr)
		galaxylib.GalaxyLogger.Infoln(string(pnrbuf))

	})

}
