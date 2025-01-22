package data

type Whiteboard struct {
	ID           int    `gorm:"primaryKey;autoIncrement:true"`
	Width        int    `gorm:"not null"`
	Height       int    `gorm:"not null"`
	DefaultColor string `gorm:"not null;type:varchar(6)"`
}

type Pixel struct {
	WhiteboardID int        `gorm:"primaryKey;autoIncrement:false"`
	X            int        `gorm:"primaryKey;autoIncrement:false"`
	Y            int        `gorm:"primaryKey;autoIncrement:false"`
	Color        string     `gorm:"not null;type:varchar(6)"`
	Whiteboard   Whiteboard `gorm:"foreignKey:WhiteboardID" json:"-"`
}
