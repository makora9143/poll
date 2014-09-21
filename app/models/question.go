package models

import (
	"github.com/revel/revel"
	"time"
)

type Question struct {
	QuestionId   int
	QuestionText string
	PubTime      time.Time
}

func (q *Question) Validate(v *revel.Validation) {
	v.Check(q.QuestionText,
		revel.Required{},
		revel.MaxSize{200},
	)
	v.Required(q.PubTime)
}
