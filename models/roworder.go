package models

import "github.com/jinzhu/gorm"

type RowOrder struct {
	gorm.Model
	GroupID    string
	PersonName string
	Gender     string
	IDCard     string
	Mobile     string
	OBPnrCode  string
	OBDate     string
	OBDep      string //运程出发地
	OBArr      string //去程目的地
	OBFlyNo    string //去程航班号
	OBDeptm    string //去程出发时间
	OBArrtm    string //去程到达时间
	IBPnrCode  string //回PNR
	IBDate     string //回程出发时间
	IBDep      string //回程出发地
	IBArr      string //回程目的地
	IBFlyNo    string //回程航班号
	IBDeptm    string //回程出发时间
	IBArrtm    string //回程到达时间

}

func (r *RowOrder) FirstOrCreate(db *gorm.DB) int64 {

	var row *RowOrder
	result := db.Where(&RowOrder{GroupID: r.GroupID, PersonName: r.PersonName, IDCard: r.IDCard}).Assign(&r).FirstOrCreate(&row)
	return result.RowsAffected
}
