package PoliteDog

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type DogError struct {
	err error
	fn  DogErrorFunc
}

func NewError() *DogError {
	return &DogError{}
}

// DogErrorFunc 暴露给开发者的自定义错误处理函数
type DogErrorFunc func(der *DogError)

func (der *DogError) Result(fn DogErrorFunc) {
	der.fn = fn
}

func (der *DogError) Error() string {
	return der.err.Error()
}

func (der *DogError) Put(err error) {
	if err != nil {
		der.err = err
		panic(der)
	}
}

func (der *DogError) ExecuteError() {
	der.fn(der)
}

// Recovery 异常统一处理
func Recovery(ctx *Context) {
	defer func() {
		err := recover()
		if err != nil {
			var der *DogError
			if errors.As(err.(error), &der) {
				der.fn(der)
			}

			ctx.e.logger.Error(getStackFrame(err))
			ctx.Data(http.StatusInternalServerError, nil)
		}
	}()

	ctx.Next()
}

// 打印栈帧
func getStackFrame(err any) string {
	var sb strings.Builder
	var pcs [32]uintptr

	sb.WriteString(fmt.Sprintf("?%v", err))

	n := runtime.Callers(3, pcs[:])
	for _, pc := range pcs[0:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)

		sb.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}

	sb.WriteString("\n")

	return sb.String()
}
