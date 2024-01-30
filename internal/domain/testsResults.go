package domain

import (
	"Contest/internal/enums"
)

type TestsResult struct {
	ResultCode  enums.TestResultCode
	Description string
	Points      int
}
