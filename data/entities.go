package data

type Pixel struct {
	X     int    `gorm:"primaryKey;autoIncrement:false;size:6"`
	Y     int    `gorm:"primaryKey;autoIncrement:false;size:6"`
	Color string `gorm:"not null"`
}
