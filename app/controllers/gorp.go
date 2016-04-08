package controllers

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	r "github.com/revel/revel"
	"github.com/revel/modules/db/app"
	"dashboard/app/models"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTable(models.User{}).SetKeys(true, "UserID")
	t.ColMap("Password").Transient = true
	setColumnSizes(t, map[string]int{
		"Username": 60,
		"Email":    100,
	})

	t = Dbm.AddTable(models.Hotel{}).SetKeys(true, "HotelId")
	setColumnSizes(t, map[string]int{
		"Name":    50,
		"Address": 100,
		"City":    40,
		"State":   6,
		"Zip":     6,
		"Country": 40,
	})

	t = Dbm.AddTable(models.Booking{}).SetKeys(true, "BookingId")
	t.ColMap("User").Transient = true
	t.ColMap("Hotel").Transient = true
	t.ColMap("CheckInDate").Transient = true
	t.ColMap("CheckOutDate").Transient = true
	setColumnSizes(t, map[string]int{
		"CardNumber": 16,
		"NameOnCard": 50,
	})


	t = Dbm.AddTable(models.Product{}).SetKeys(true, "ProductID")
	setColumnSizes(t, map[string]int{
		"Name":    50,
		"Description": 100,
		"Quantity":    10,
	})


	t = Dbm.AddTable(models.Store{}).SetKeys(true, "StoreID")
	t.ColMap("Brand").Transient = true
	setColumnSizes(t, map[string]int{
		"Name":    50,
		"Description": 100,
		"Address":    100,
	})

	t = Dbm.AddTable(models.Brand{}).SetKeys(true, "BrandID")
	setColumnSizes(t, map[string]int{
		"Name":    50,
	})

	t = Dbm.AddTable(models.AssoBrandProduct{}).SetKeys(true, "AssoBrandProductID")
	t.ColMap("Brand").Transient = true
	t.ColMap("Product").Transient = true
	setColumnSizes(t, map[string]int{
		"BrandID":    	50,
		"ProductID":	50,
	})


	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTables()


	brands := []*models.Brand{
		&models.Brand{1, "Pret A Manger"},
	}
	for _, brand := range brands {
		if err := Dbm.Insert(brand); err != nil {
			panic(err)
		}
	}


	stores := []*models.Store{
		&models.Store{StoreID:1, Name:"Pret A Manger Farringdon",BrandID: 1, Description: "desc",
			Address: "Unit 1, Nexus Place, 24 Farringdon Street, London EC4A 4JA", Latitude:51.515620, Longitude: -0.104707},
		&models.Store{StoreID:2, Name:"Pret A Manger Warren Street",BrandID: 1, Description: "desc",
			Address: "14 Warren St, London W1T 5LL", Latitude:51.524167, Longitude: -0.138822},
	}
	for _, store := range stores {
		if err := Dbm.Insert(store); err != nil {
			panic(err)
		}
	}

	/*bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}

	hotels := []*models.Hotel{
		&models.Hotel{0, "Marriott Courtyard", "Tower Pl, Buckhead", "Atlanta", "GA", "30305", "USA", 120},
		&models.Hotel{0, "W Hotel", "Union Square, Manhattan", "New York", "NY", "10011", "USA", 450},
		&models.Hotel{0, "Hotel Rouge", "1315 16th St NW", "Washington", "DC", "20036", "USA", 250},
	}
	for _, hotel := range hotels {
		if err := Dbm.Insert(hotel); err != nil {
			panic(err)
		}
	}*/
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
