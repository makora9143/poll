package models

import (
	"github.com/coopernurse/gorp"
	"github.com/revel/revel"
)

type Choice struct {
	ChoiceId   int
	QuestionId int
	ChoiceText string
	Vote       int
	Question   *Question
}

func (c Choice) Validate(v *revel.Validation) {
	v.Check(c.ChoiceText,
		revel.Required{},
		revel.MaxSize{200},
	)
	v.Required(c.Vote)
	v.Required(c.Question)
}

func (c *Choice) PreInsert(_ gorp.SqlExecutor) error {
	c.QuestionId = c.Question.QuestionId
	return nil
}
