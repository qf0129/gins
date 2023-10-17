package simple

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    uint      `gorm:"primaryKey;" json:"id"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime;comment:'UpdatedTime'" json:"utime"`
}

type BaseUidModel struct {
	Id    string    `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime;comment:'UpdatedTime'" json:"utime"`
}

type BaseAssociatedModel struct {
	Ctime time.Time      `gorm:"autoCreateTime;comment:'CreatedTime'" json:"ctime"`
	Dtime gorm.DeletedAt `gorm:"index;comment:'DeletedTime'" json:"-" `
}

type BaseUidModelWithDel struct {
	BaseUidModel
	Dtime gorm.DeletedAt `gorm:"index;comment:'DeletedTime'" json:"-" `
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Id == "" {
		m.Id = xid.New().String()
	}
	return
}

type User struct {
	BaseUidModelWithDel
	Username     string `gorm:"index;type:varchar(50)" json:"username"`
	PasswordHash string `gorm:"type:varchar(200)" json:"-"`
}
