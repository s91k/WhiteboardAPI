package data

type Pixel struct {
	X     int    `gorm:"primaryKey;autoIncrement:false;not null"`
	Y     int    `gorm:"primaryKey;autoIncrement:false;not null"`
	Color string `gorm:"not null;type:varchar(7)"`
}
