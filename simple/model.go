package simple

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    uint      `gorm:"primaryKey;" json:"id" form:"id"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime;comment:'UpdatedTime'" json:"utime"`
}

type BaseUidModel struct {
	Id    string    `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime time.Time `gorm:"autoCreateTime:milli;comment:'CreatedTime'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime:milli;comment:'UpdatedTime'" json:"utime"`
}

type BaseUidModelWithDel struct {
	BaseUidModel
	Dtime gorm.DeletedAt `gorm:"index;comment:'DeletedTime'" json:"-"`
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = xid.New().String()
	return
}

// func (m *BaseUidModelWithDel) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.Id = xid.New().String()
// 	return
// }

type User struct {
	BaseUidModelWithDel
	Username     string `gorm:"index;type:varchar(50)" json:"username"  form:"username"`
	PasswordHash string `gorm:"type:varchar(200)" json:"-"  form:"password_hash"`
}
