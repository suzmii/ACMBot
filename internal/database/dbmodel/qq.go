package dbmodel

import (
	"gorm.io/gorm"
)

type QQBind struct {
	gorm.Model
	QID              uint64 `gorm:"uniqueIndex"`
	CodeforcesUserID uint   `gorm:"uniqueIndex"`
	QName            string
}

type QQGroup struct {
	gorm.Model
	GroupID uint64 `gorm:"index"`
	QID     uint64
}

type QQUser struct {
	QID              uint64
	CodeforcesRating uint
	QName            string
}

type QQGroupRank struct {
	GroupID uint64
	QQUsers []*QQUser
}

//
//func MigrateQQ() error {
//	return database.db.AutoMigrate(&QQBind{}, &QQGroup{})
//}
//
//var (
//	bindCfID datasync.Map
//	bindQID  datasync.Map
//)
//
//func BindQQToCodeforces(qqBind QQBind, group QQGroup) error {
//	v, _ := bindCfID.LoadOrStore(qqBind.CodeforcesUserID, &datasync.Mutex{})
//	u, _ := bindQID.LoadOrStore(qqBind.QID, &datasync.Mutex{})
//	cfIDLock := v.(*datasync.Mutex)
//	qIDLock := u.(*datasync.Mutex)
//	cfIDLock.Lock()
//	qIDLock.Lock()
//	defer cfIDLock.Unlock()
//	defer qIDLock.Unlock()
//	var IDUsed QQBind
//	if errs := database.db.Where(`codeforces_user_id = ?`, qqBind.CodeforcesUserID).First(&IDUsed).Error; errs == nil {
//		if IDUsed.QID != qqBind.QID {
//			return errs.ErrHandleHasBindByOthers
//		}
//	} else if !errors.Is(errs, gorm.ErrRecordNotFound) {
//		return errs
//	} else {
//		if errs = database.db.Where(`q_id = ?`, qqBind.QID).First(&IDUsed).Error; errs == nil {
//			if errs = database.db.Model(&QQBind{}).Where("q_id = ?", qqBind.QID).Updates(qqBind).Error; errs != nil {
//				return errs
//			}
//		} else if !errors.Is(errs, gorm.ErrRecordNotFound) {
//			return errs
//		}
//		if errs = database.db.Create(&qqBind).Error; errs != nil {
//			return errs
//		}
//	}
//	var existGroup QQGroup
//	if errs := database.db.Where(`group_id = ? and q_id = ?`, group.GroupID, group.QID).First(&existGroup).Error; errs == nil {
//		return nil
//	}
//	if errs := database.db.Create(&group).Error; errs != nil {
//		return errs
//	}
//	return nil
//}
//
//func CheckBind(QID uint) (bool, error) {
//	if errs := database.db.Where("q_id = ?", QID).First(&QQBind{}).Error; errs != nil {
//		if errors.Is(errs, gorm.ErrRecordNotFound) {
//			return false, nil
//		} else {
//			return false, errs
//		}
//	}
//	return true, nil
//}
//
//func GetQQGroupRank(QQGroupNumber uint64) (*QQGroupRank, error) {
//	var QIDs []uint64
//	if errs := database.db.Model(&QQGroup{}).Where("group_id = ?", QQGroupNumber).Pluck("q_id", &QIDs).Error; errs != nil {
//		return nil, errs
//	}
//	var qqGroupRank = QQGroupRank{
//		GroupID: QQGroupNumber,
//		Users: make([]*QQUser, 0),
//	}
//	for _, QID := range QIDs {
//		var qqBind QQBind
//		if errs := database.db.Where("q_id = ?", QID).First(&qqBind).Error; errs != nil {
//			return nil, errs
//		}
//		var cfRating uint
//		if errs := database.db.Model(&CodeforcesUser{}).Select("rating").Where("id = ?", qqBind.CodeforcesUserID).First(&cfRating).Error; errs != nil {
//			return nil, errs
//		}
//		var qqUser = QQUser{
//			QID:              QID,
//			Username:            qqBind.Username,
//			CodeforcesRating: cfRating,
//		}
//		qqGroupRank.Users = append(qqGroupRank.Users, &qqUser)
//	}
//	return &qqGroupRank, nil
//}
