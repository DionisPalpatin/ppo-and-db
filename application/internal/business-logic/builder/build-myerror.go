package myerror_builder

import (
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
)


type MyErrorBuilder struct {
	MyError *bl.MyError
}

func NewMyErrorBuilder() *MyErrorBuilder {
	return &MyErrorBuilder{
		MyError: &bl.MyError{
			ErrNum: bl.Ok,
			FuncName: "",
			Module: "",
		},
	}
}

func (b *MyErrorBuilder) WithErrNum(err_num int) *MyErrorBuilder {
	b.MyError.ErrNum = err_num
	return b
}

func (b *MyErrorBuilder) WithFuncName(func_name string) *MyErrorBuilder {
	b.MyError.FuncName = func_name
	return b
}

func (b *MyErrorBuilder) WithModule(module string) *MyErrorBuilder {
	b.MyError.Module = module
	return b
}