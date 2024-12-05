package bl

import (
	"fmt"
)

type MyError struct {
	ErrNum   int
	FuncName string
	Module   string
}

const (
	Ok = iota

	NoSuchUser = iota
	NoSuchTeam = iota
	NoSuchNote = iota
	NoSuchColl = iota
	NoSuchSect = iota

	EmptyResult = iota

	ErrAccessDenied    = iota
	OperationError     = iota
	ErrSearchParameter = iota
	ErrInParameter     = iota

	ErrNoFile   = iota
	ErrReadFile = iota

	DatabaseError = iota
)

func CreateError(errNum int, funcName string, module string) *MyError {
	myErr := new(MyError)

	myErr.ErrNum = errNum
	myErr.FuncName = funcName
	myErr.Module = module

	return myErr
}

func (m *MyError) Error() string {
	return fmt.Sprintf("Error in function %s (module %s): errnum %d", m.FuncName, m.Module, m.ErrNum)
}

func (m *MyError) ConcatenateFields() string {
	return fmt.Sprintf("Error in function %s (module %s): errnum %d", m.FuncName, m.Module, m.ErrNum)
}

func (m *MyError) ConcatenateWithExternalErr(extErr error) string {
	return fmt.Sprintf("Error in function %s (module %s): %s", m.FuncName, m.Module, extErr.Error())
}
