// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameBasis = "bases"

// Basis mapped from table <bases>
type Basis struct {
	ID               int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title            string    `gorm:"column:title" json:"title"`
	TitleShort1      string    `gorm:"column:title_short1" json:"title_short1"`
	TitleShort2      string    `gorm:"column:title_short2" json:"title_short2"`
	TitleShort3      string    `gorm:"column:title_short3" json:"title_short3"`
	TitleEn          string    `gorm:"column:title_en" json:"title_en"`
	PublicURL        string    `gorm:"column:public_url" json:"public_url"`
	Original         string    `gorm:"column:original" json:"original"`
	TwitterAccount   string    `gorm:"column:twitter_account" json:"twitter_account"`
	TwitterHashTag   string    `gorm:"column:twitter_hash_tag" json:"twitter_hash_tag"`
	Facebook         string    `gorm:"column:facebook" json:"facebook"`
	CoursID          int32     `gorm:"column:cours_id" json:"cours_id"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	NiconicoURL      string    `gorm:"column:niconico_url" json:"niconico_url"`
	WebRadioURL      string    `gorm:"column:web_radio_url" json:"web_radio_url"`
	Sex              int32     `gorm:"column:sex" json:"sex"`
	Sequel           int32     `gorm:"column:sequel" json:"sequel"`
	CityCode         int32     `gorm:"column:city_code" json:"city_code"`
	CityName         string    `gorm:"column:city_name" json:"city_name"`
	ProductCompanies string    `gorm:"column:product_companies" json:"product_companies"`
}

// TableName Basis's table name
func (*Basis) TableName() string {
	return TableNameBasis
}