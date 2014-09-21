package controllers

import (
	"github.com/makora/poll/app/models"
	"github.com/makora/poll/app/routes"
	"github.com/revel/revel"
)

type Poll struct {
	App
}

func (p Poll) Index() revel.Result {
	results, err := p.Txn.Select(models.Question{},
		`select * from Question`)
	if err != nil {
		panic(err)
	}

	var questions []*models.Question
	for _, r := range results {
		q := r.(*models.Question)
		questions = append(questions, q)
	}

	return p.Render(questions)
}

func (p Poll) Detail(id int) revel.Result {
	question := p.loadQuestionById(id)
	results, err := p.Txn.Select(models.Choice{},
		`select * from Choice where QuestionId = ?`, question.QuestionId)
	if err != nil {
		panic(err)
	}

	var choices []*models.Choice
	for _, r := range results {
		c := r.(*models.Choice)
		choices = append(choices, c)
	}
	return p.Render(id, question, choices)
}

func (p Poll) Results(id int) revel.Result {
	question := p.loadQuestionById(id)
	results, err := p.Txn.Select(models.Choice{},
		`select * from Choice where QuestionId`, question.QuestionId)
	if err != nil {
		panic(err)
	}

	var choices []*models.Choice
	for _, r := range results {
		c := r.(*models.Choice)
		choices = append(choices, c)
	}

	return p.Render(question, choices)
}

func (p Poll) Vote(id int, choice int) revel.Result {
	result, err := p.Txn.Get(models.Choice{}, choice)
	if err != nil {
		panic(err)
	}
	selectChoice := result.(*models.Choice)
	selectChoice.Vote++
	p.Txn.Exec(`update Choice set Vote = ? where ChoiceId = ?`, selectChoice.Vote, choice)
	return p.Redirect(routes.Poll.Results(id))
}

func (p Poll) loadQuestionById(id int) *models.Question {
	q, err := p.Txn.Get(models.Question{}, id)
	if err != nil {
		panic(err)
	}
	if q == nil {
		return nil
	}
	return q.(*models.Question)
}
