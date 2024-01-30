package domain

type Test struct {
	ID             int    `json:"id,string"`
	TaskID         int    `json:"taskID,string"`
	Input          string `json:"input"`
	ExpectedResult string `json:"expectedResult"`
	Points         int    `json:"points,string"`
}
