package model

type ImageData struct {
	Id  int64
	URL string `xorm:"'url'"`
}
