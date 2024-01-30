package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ICompileService interface {
	CompileCPP(code string) (string, error)
}

type CompileService struct {
}

var (
	CompileError = errors.New("Compile Error")
)

func NewCompileSevice() *CompileService {
	return &CompileService{}
}

func (c *CompileService) CompileCPP(code string) (string, error) {
	fileName := time.Now().String()
	file, err := os.Create(fileName + ".cpp")
	defer file.Close()
	if err != nil {
		return "", err
	}

	file.WriteString(code)

	cmd := exec.Command("g++", "-o", fileName+".exe", "-x", "c++", fileName+".cpp")
	err = cmd.Run()

	if err != nil {
		return "", fmt.Errorf("%w: %w", CompileError, err.Error())
	}

	cmd = exec.Command("chmod", "+x", fileName+".exe")
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("In CompileService(CompileCPP): %w", err)
	}

	return fileName + ".exe", nil
}
