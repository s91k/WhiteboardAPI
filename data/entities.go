package data

type Pixel struct {
	X     int    `gorm:"primaryKey;autoIncrement:false"`
	Y     int    `gorm:"primaryKey;autoIncrement:false"`
	Color string `gorm:"not null;type:varchar(6)"`
}
