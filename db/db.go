package db

import (
	"fmt"
	"kua-mei-a-api/model"
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var engine *xorm.Engine

func init() {
	if err := getDBEngine(); err != nil {
		log.Fatal(err)
	}
	engine.CreateTables(model.ImageData{})
}

func getDBEngine() error {
	var err error
	engine, err = xorm.NewEngine("postgres", datastoreName)
	if err != nil {
		return err
	}
	engine.ShowSQL() //菜鸟必备

	err = engine.Ping()
	if err != nil {
		return err
	}
	fmt.Println("connect postgresql success")
	return nil
}

// Count 查詢總數
func Count() (int64, error) {
	urls := new(model.ImageData)
	return engine.Where("id >?", 1).Count(urls)
}

// SelectOne 可以用Get查询单个元素
func SelectOne(id int64) *model.ImageData {
	var image model.ImageData
	engine.Id(id).Get(&image)
	//engine.Alias("u").Where("u.id=?",id).Get(&user)
	return &image
}

// SessionInsert 使用Session來批量處理
func SessionInsert(images []*model.ImageData) bool {
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	for _, image := range images {
		_, err := session.Insert(image)
		if err != nil {
			session.Rollback()
			log.Println(err)
		}
	}

	err := session.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Insert 添加
func Insert(image *model.ImageData) bool {
	rows, err := engine.Insert(image)
	if err != nil {
		log.Println(err)
		return false
	}
	if rows == 0 {
		return false
	}
	return true
}

// Delete 删除
func Delete(image *model.ImageData) bool {
	rows, err := engine.Delete(image)
	if err != nil {
		log.Println(err)
		return false
	}
	if rows == 0 {
		return false
	}
	return true
}

// Update 更新
func Update(image *model.ImageData) bool {
	//Update(bean interface{}, condiBeans ...interface{}) bean是需要更新的bean,condiBeans是条件
	rows, err := engine.Update(image, model.ImageData{Id: image.Id})
	if err != nil {
		log.Println(err)
		return false
	}
	if rows > 0 {
		return true
	}
	return false
}
