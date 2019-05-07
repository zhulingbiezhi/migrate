package user_v2

import (
	"github.com/jinzhu/gorm"
	"migrate/db"
	"time"
)

type UserV2 struct {
	ID         int64     `json:"id" gorm:"column:id"`
	Name       string    `json:"name_v_2" gorm:"column:name_v_2"`
	FirstName  string    `json:"first_name_v_2" gorm:"column:first_name_v_2"`
	LastName   string    `json:"last_name_v_2" gorm:"column:last_name_v_2"`
	Email      string    `json:"email_v_2" gorm:"column:email_v_2"`
	Sex        string    `json:"sex_v_2" gorm:"column:sex_v_2"`
	Flag       int       `json:"flag_v_2" gorm:"column:flag_v_2"`
	CreateTime time.Time `json:"create_time_v_2" gorm:"column:create_time_v_2"`
}

func (*UserV2) TableName() string {
	return "user_v2"
}

func (u *UserV2) Insert() (int64, error) {
	row := new(UserV2)
	d := db.DB.Create(u).Scan(&row)
	if d.Error != nil {
		return 0, d.Error
	}
	return row.ID, nil
}

func Select(scopes ...func(scope *gorm.DB) *gorm.DB) ([]*UserV2, error) {
	var users []*UserV2
	d := db.DB.Model(&UserV2{}).Scopes(scopes...).Find(&users)
	if d.Error != nil {
		if d.Error == gorm.ErrRecordNotFound {
			return users, nil
		}
		return nil, d.Error
	}
	return users, nil
}
func MaxID() (int64, error) {
	u := UserV2{}
	d := db.DB.Model(&UserV2{}).Last(&u).Debug()
	if d.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return u.ID, d.Error
}
