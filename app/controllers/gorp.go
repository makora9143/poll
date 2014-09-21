package controllers

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/makora/poll/app/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/db/app"
	"time"
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

	t := Dbm.AddTable(models.Question{}).SetKeys(true, "QuestionId")
	setColumnSizes(t, map[string]int{
		"QuestionText": 200,
	})

	t = Dbm.AddTable(models.Choice{}).SetKeys(true, "ChoiceId")
	t.ColMap("Question").Transient = true
	setColumnSizes(t, map[string]int{
		"ChoiceText": 200,
	})

	Dbm.TraceOn("[gorp]", revel.INFO)
	Dbm.CreateTables()

	questions := []*models.Question{
		&models.Question{0, "What's new?", time.Now()},
	}
	for _, question := range questions {
		if err := Dbm.Insert(question); err != nil {
			panic(err)
		}
	}
	choices := []*models.Choice{
		&models.Choice{ChoiceId: 0, ChoiceText: "Not much", Vote: 0, Question: questions[0]},
		&models.Choice{0, questions[0].QuestionId, "The sky", 0, questions[0]},
		&models.Choice{0, questions[0].QuestionId, "Just hacking again", 0, questions[0]},
	}
	for _, choice := range choices {
		if err := Dbm.Insert(choice); err != nil {
			panic(err)
		}
	}

}

type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

func (g *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	g.Txn = txn
	return nil
}

func (g *GorpController) Commit() revel.Result {
	if g.Txn == nil {
		return nil
	}
	if err := g.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	g.Txn = nil
	return nil
}

func (g *GorpController) Rollback() revel.Result {
	if g.Txn == nil {
		return nil
	}
	if err := g.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	g.Txn = nil
	return nil
}
