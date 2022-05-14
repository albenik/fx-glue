package fxglue

import (
	"fmt"
	"io"

	"go.uber.org/fx"
)

type ErrorHandlerFunc func(err error)

func (fn ErrorHandlerFunc) HandleError(err error) {
	fn(err)
}

func PrintErrorHandler(w io.Writer) fx.ErrorHandler {
	return ErrorHandlerFunc(func(err error) {
		str, verr := fx.VisualizeError(err)
		if verr != nil {
			fmt.Fprintln(w, verr)
			return
		}
		fmt.Fprintln(w, str)
	})
}
