package model

// User 用户模型
type User struct {
	ID    int32  `json:"id" gorm:"primaryKey;autoIncrement"` // 主键且自增
	Name  string `json:"name" gorm:"type:varchar(100);not null"`
	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
}
