package model

type Setting struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	BrandName     string `gorm:"size:255;not null" json:"brandName"`
	UrlLogo       string `gorm:"size:255;not null" json:"urlLogo"`
	FileTypeAllow string `gorm:"size:255;not null" json:"fileTypeAllow"`
	IsMaintenance bool   `gorm:"not null" json:"isMaintenance"`
}
