// Code generated by 'freedom new-project base'
package repository

import (
	"github.com/8treenet/freedom"
	"github.com/jinzhu/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *Default {
			return &Default{}
		})
	})
}

// Default .
type Default struct {
	freedom.Repository
}

// GetIP .
func (repo *Default) GetIP() string {
	//repo.DB().Find()
	repo.Worker.Logger().Info("I'm Repository GetIP")
	return repo.Worker.IrisContext().RemoteAddr()
}

// GetUA - implement DefaultRepoInterface interface
func (repo *Default) GetUA() string {
	repo.Worker.Logger().Info("I'm Repository GetUA")
	return repo.Worker.IrisContext().Request().UserAgent()
}

// db .
func (repo *Default) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	db = db.New()
	db.SetLogger(repo.Worker.Logger())
	return db
}

/*
	// xorm
	func (repo *Default) db() *xorm.Engine {
		var db *xorm.Engine
		if err := repo.FetchDB(&db); err != nil {
			panic(err)
		}
		return db
	}
	func main {
		app.InstallDB(func() interface{} {
			db, _ := xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/xorm?charset=utf8")
			return db
		})
	}
*/
