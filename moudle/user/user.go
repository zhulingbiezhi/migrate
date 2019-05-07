package user

import (
	"github.com/jinzhu/gorm"
	"migrate/db"
	"time"
)

type User struct {
	ID         int64     `json:"id" gorm:"column:id"`
	Name       string    `json:"name" gorm:"column:name"`
	FirstName  string    `json:"first_name" gorm:"column:first_name"`
	LastName   string    `json:"last_name" gorm:"column:last_name"`
	Email      string    `json:"email" gorm:"column:email"`
	Sex        string    `json:"sex" gorm:"column:sex"`
	Flag       int       `json:"flag" gorm:"column:flag"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func (*User) TableName() string {
	return "user"
}

func ScopeIDRange(minID, maxID int64) func(scope *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		scope = scope.Where("id > ? and id <= ?", minID, maxID)
		return scope
	}
}

func ScopeID(id int64) func(scope *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		scope = scope.Where("id = ?", id)
		return scope
	}
}

func ScopeStartTime(t *time.Time) func(scope *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		scope = scope.Where("create_time > ?", t)
		return scope
	}
}

func ScopeEndTime(t *time.Time) func(scope *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		scope = scope.Where("create_time < ?", t)
		return scope
	}
}

func (u *User) Insert() (int64, error) {
	row := new(User)
	d := db.DB.Create(u).Scan(&row)
	if d.Error != nil {
		return 0, d.Error
	}
	return row.ID, nil
}

func Select(scopes ...func(scope *gorm.DB) *gorm.DB) ([]*User, error) {
	var users []*User
	d := db.DB.Model(&User{}).Scopes(scopes...).Find(&users).Order("id asc")
	if d.Error != nil {
		if d.Error == gorm.ErrRecordNotFound {
			return users, nil
		}
		return nil, d.Error
	}
	return users, nil
}

func MaxID() (int64, error) {
	u := User{}
	d := db.DB.Model(&User{}).Last(&u)
	if d.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return u.ID, d.Error
}
