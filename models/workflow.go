package models

import (
	"time"
)

//Workflow is a structure to store the service request details
type Workflow struct {
	Id          string    `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Steps       []Step    `json:"steps" binding:"required,gt=0,dive"`
}

//Create a new work flow for a given service flow and return service flow details
func CreateWorkflow(workflowRequest Workflow) Workflow {
	for i := 0; i < len(workflowRequest.Steps); i++ {
		workflowRequest.Steps[i].Id = i + 1
	}
	return newServiceFlow(workflowRequest)
}

func newServiceFlow(workflow Workflow) Workflow {
	return Workflow{Id: workflow.Id, Name: workflow.Name, Description: workflow.Description, Enabled: true, CreatedAt: time.Time{}, UpdatedAt: time.Time{}, Steps: workflow.Steps}
}

type PGWorkflow struct {
	tableName   struct{} `pg:"workflows"`
	Id          string
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Steps       []Step
}

func (workflow Workflow) ToPGWorkflow() PGWorkflow {
	return PGWorkflow{
		Id:          workflow.Id,
		Name:        workflow.Name,
		Description: workflow.Description,
		Enabled:     workflow.Enabled,
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
		Steps:       workflow.Steps,
	}
}

func (pgWorkflow PGWorkflow) ToWorkflow() Workflow {
	return Workflow{
		Id:          pgWorkflow.Id,
		Name:        pgWorkflow.Name,
		Description: pgWorkflow.Description,
		Enabled:     pgWorkflow.Enabled,
		CreatedAt:   pgWorkflow.CreatedAt,
		UpdatedAt:   pgWorkflow.UpdatedAt,
		Steps:       pgWorkflow.Steps,
	}
}
