package enums

type TestResultCode string

const (
	TimeLimit       TestResultCode = "TL"
	MemoryLimit     TestResultCode = "ML"
	CompileError    TestResultCode = "CE"
	RuntimeError    TestResultCode = "RE"
	Succes          TestResultCode = "SC"
	IncorrectAnswer TestResultCode = "IA"
)
