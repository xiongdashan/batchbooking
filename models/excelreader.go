package models

import (
	"fmt"
	"strings"

	ibelib "github.com/otwdev/ibepluslib/models"

	"github.com/otwdev/galaxylib"
	"github.com/tealeg/xlsx"
)

const (
	id           = 0
	pnr          = 1
	name         = 2
	gender       = 3
	idcard       = 4
	mobile       = 5
	obDate       = 6
	obdepart     = 7
	obarriv      = 8
	obflyNo      = 9
	obdepartTime = 10
	obarrivTime  = 11
	ibpnr        = 12
	ibDate       = 13
	ibdepart     = 14
	ibarriv      = 15
	ibflyNo      = 16
	ibdepartTime = 17
	ibarrivTime  = 18
)

type ExcelReader struct {
	FileName   string
	currentRow *xlsx.Row
	Office     string
}

func NewExcelReader(office, fileName string) *ExcelReader {
	return &ExcelReader{
		FileName: fileName,
		Office:   office,
	}
}

func (e *ExcelReader) Read(fn func(order *ibelib.OrderInfo)) {
	xls, err := xlsx.OpenFile(e.FileName)
	if err != nil {
		galaxylib.GalaxyLogger.Errorln(err)
		return
	}
	sheet := xls.Sheets[0]

	fmt.Println(len(sheet.Rows))

	for i, row := range sheet.Rows {

		if i < 1 {
			continue
		}

		e.currentRow = row

		cellsCount := len(row.Cells)

		if cellsCount < 12 || e.trimCell(pnr) != "" {
			continue
		}

		oOrder := e.generateOrder(name, gender, idcard, mobile, obdepart, obDate, obdepartTime, obarriv, obDate, obarrivTime, obflyNo)

		if oOrder != nil {
			fn(oOrder)
		}

		if cellsCount >= 18 {
			iOrder := e.generateOrder(name, gender, idcard, mobile, ibdepart, ibDate, ibdepartTime, ibarriv, ibDate, ibarrivTime, ibflyNo)
			if iOrder != nil {
				fn(iOrder)
			}
		}

	}

}

func (e *ExcelReader) generateOrder(nameOps, genderOps, idCardOps, mobileOps, depOps, depDateOps, depTimeOps,
	arrOps, arrDateOps, arrTimOps, flyNoOps int) *ibelib.OrderInfo {

	traverller := &ibelib.Traveler{
		PersonName: e.trimCell(nameOps),
		Gender:     e.transGender(genderOps), //e.trimCell(gender),
		IDCardType: "NI",
		IDCardNo:   e.trimCell(idCardOps),
		Mobile:     e.trimCell(mobileOps),
		Type:       string(ibelib.Adult),
	}

	if traverller.PersonName == "" || traverller.IDCardNo == "" {
		return nil
	}

	code, num := e.transFlyNo(flyNoOps)

	fltSegm := &ibelib.FlightSegment{
		DepartCityCode:   e.trimCell(depOps),
		DepTime:          e.trimCell(depTimeOps),
		FlyDate:          e.transDate(depDateOps), //e.trimCell(obDate),
		FlyNo:            num,                     //e.trimCell(obflyNo),
		MarketingAirLine: code,
		ArrTime:          e.trimCell(arrTimOps),
		//ArrDate:          e.transDate(arrDateOps), //e.trimCell(obDate),
		ArriveCityCode: e.trimCell(arrOps),
		TripSeq:        1,
		TripType:       "OB",
	}

	fltSegm.ArrDate = fltSegm.FlyDate

	if fltSegm.DepartCityCode == "" || fltSegm.ArriveCityCode == "" {
		return nil
	}

	p := &ibelib.PnrInfo{}
	p.OfficeNumber = e.Office
	p.AgencyCity = "BJS"

	p.TravelerInfos = append(p.TravelerInfos, traverller)
	p.FlightSegments = append(p.FlightSegments, fltSegm)

	order := &ibelib.OrderInfo{}
	order.PnrInofs = append(order.PnrInofs, p)
	order.ContactInfo = &ibelib.ContactInfo{
		PersonName:  "XXX",
		MobilePhone: "132333333",
	}
	return order
}

func (e *ExcelReader) trimCell(ops int) string {

	return strings.TrimSpace(e.currentRow.Cells[ops].String())
}

func (e *ExcelReader) transGender(ops int) string {
	val := e.trimCell(ops)
	if val == "男" {
		return "F"
	}
	if val == "女" {
		return "M"
	}
	return val
}

func (e *ExcelReader) transFlyNo(ops int) (string, string) {
	val := e.trimCell(ops)
	code := val[:2]
	number := val[2:]
	return code, number
}

func (e *ExcelReader) transDate(ops int) string {
	t, _ := e.currentRow.Cells[ops].GetTime(false)
	val := t.Format("2006-01-02")
	// val = strings.Replace(val, "月", "-", -1)
	// val = strings.Replace(val, "日", "", -1)
	return val
}
