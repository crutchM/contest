package services

import (
	. "Contest/internal/domain"
	"Contest/internal/enums"
	"Contest/internal/storage"
	"Contest/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ITestService interface {
	RunTest(taskID int, language enums.Language, code string) (TestsResult, error)
	GetTest(id int) (Test, error)
	AddTest(test Test) error
	DeleteTest(id int) error
	UpdateTest(id int, newTest Test) error
	GetTests() ([]Test, error)
	GetTestsByTaskID(taskID int) ([]Test, error)
}

type TestService struct {
	compileService ICompileService
	testRepository storage.Repository[Test]
}

var (
	ProgramError         = errors.New("Program error")
	TimeLimitError       = errors.New("Time limit error")
	UnknownLanguageError = errors.New("Unknown language")
	TestsNotFoundError   = errors.New("Tests not found")
)

func NewTestService(compileService ICompileService, testRepository storage.Repository[Test]) *TestService {
	return &TestService{
		compileService: compileService,
		testRepository: testRepository,
	}
}

func (s *TestService) RunTest(taskID int, language enums.Language, code string) (TestsResult, error) {
	var fileName string
	var err error

	switch language {
	case enums.CPP:
		fileName, err = s.compileService.CompileCPP(code)
	case enums.CSharp:
		panic("IMPLEMENT ME PLEASE")
	case enums.Python:
		panic("GIVE ME DIE PLEASE!")
	default:
		return TestsResult{}, UnknownLanguageError
	}
	if err != nil {
		return TestsResult{}, fmt.Errorf("In TestService(RunTest): %w", err)
	}

	file, err := os.Open(fileName)
	defer os.Remove(fileName)
	defer file.Close()
	if err != nil {
		return TestsResult{}, fmt.Errorf("In TestService(RunTest): %w", err)
	}

	tests, err := s.testRepository.FindItemsByCondition(func(item Test) bool {
		return item.TaskID == taskID
	})

	if len(tests) == 0 {
		return TestsResult{}, TestsNotFoundError
	}

	if err != nil {
		return TestsResult{}, fmt.Errorf("In TestService(RunTest): %w", err)
	}

	points := 0
	for _, test := range tests {
		timeout := time.Millisecond * 10000 //TODO
		maxMemoryKB := 1024 * 1024          //TODO
		output, err := runCompiledCodeWithInput(fileName, test.Input, timeout, maxMemoryKB)
		if err != nil {
			if errors.Is(err, TimeLimitError) {
				return TestsResult{
					ResultCode:  enums.TimeLimit,
					Description: "",
					Points:      points,
				}, nil
			} else if errors.Is(err, ProgramError) {
				return TestsResult{
					ResultCode:  enums.RuntimeError,
					Description: fmt.Sprintf("Error Info: %s Output: %s", err.Error(), output),
					Points:      points,
				}, nil
			} else {
				return TestsResult{}, fmt.Errorf("In TestService(RunTests): %w")
			}
		}
		if output == test.ExpectedResult {
			points += test.Points
		} else {
			return TestsResult{
				ResultCode:  enums.IncorrectAnswer,
				Description: fmt.Sprintf("Test Failed: %d", test.ID),
				Points:      points,
			}, nil
		}
	}

	return TestsResult{
		ResultCode:  enums.Succes,
		Description: "",
		Points:      points,
	}, err
}

func runCompiledCodeWithInput(fileName string, input string, timeout time.Duration, maxMemoryKB int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "./"+fileName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("In TestService(runCompiledCodeWithInput): %w", err)
	}
	defer stdin.Close()

	fmt.Fprintln(stdin, input)

	output, err := cmd.CombinedOutput()
	if errors.Is(err, context.DeadlineExceeded) {
		return "", TimeLimitError
	}
	if err != nil {
		return string(output), fmt.Errorf("%w: %w", ProgramError, err)
	}

	return strings.ReplaceAll(string(output), "\n", ""), nil
}

var (
	ErrNotFound = errors.New("Not found in database")
)

func (s *TestService) GetTest(id int) (Test, error) {
	test, err := s.testRepository.FindItemByID(id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return Test{}, ErrNotFound
		}
		return Test{}, fmt.Errorf("In TestService(GetTest): %w", err)
	}
	return test, nil
}

func (s *TestService) AddTest(test Test) error {
	err := s.testRepository.AddItem(test)
	if err != nil {
		return fmt.Errorf("In TestService(AddTest): %w", err)
	}
	return nil
}

func (s *TestService) DeleteTest(id int) error {
	err := s.testRepository.DeleteItem(id)
	if err != nil {
		return fmt.Errorf("In TestService(DeleteTest): %w", err)
	}
	return nil
}

func (s *TestService) UpdateTest(id int, newTest Test) error {
	err := s.testRepository.UpdateItem(id, newTest)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("In TestService(UpdateTest): %w", err)
	}
	return nil
}

func (s *TestService) GetTests() ([]Test, error) {
	tests, err := s.testRepository.GetTable()
	if err != nil {
		return nil, fmt.Errorf("In TestService(GetTests): %w", err)
	}
	return tests, nil
}

func (s *TestService) GetTestsByTaskID(taskID int) ([]Test, error) {
	tests, err := s.testRepository.FindItemsByCondition(
		func(item Test) bool {
			return item.TaskID == taskID
		})
	if err != nil {
		return nil, fmt.Errorf("In TestService(GetTestsByTaskID): %w", err)
	}
	return tests, nil
}
