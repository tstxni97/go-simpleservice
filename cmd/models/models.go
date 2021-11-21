package migrations

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// // gorm.Model definition
// type Model struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// }

type User struct {
	gorm.Model
	Name         string
	Email        string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreditCards  []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
	gorm.Model
	Number    string
	UserRefer uint
}

type Dog struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;polymorphicValue:master"`
}

type Toy struct {
	ID        int
	Name      string
	OwnerID   int
	OwnerType string
}
