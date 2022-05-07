// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTwitterStatus = "twitter_statuses"

// TwitterStatus mapped from table <twitter_statuses>
type TwitterStatus struct {
	BasesID   int32     `gorm:"column:bases_id;primaryKey" json:"bases_id"`
	Follower  int32     `gorm:"column:follower" json:"follower"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName TwitterStatus's table name
func (*TwitterStatus) TableName() string {
	return TableNameTwitterStatus
}
