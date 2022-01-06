package main

import (
	"fmt"
	"io"
	"testing"

	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/spf13/cobra"
)

func TestHandleError(t *testing.T) {
	nc := &cobra.Command{}

	tc := []struct {
		name string
		err  error
		code exitCode
	}{
		{
			name: "test exitcode for nil error",
			err:  nil,
			code: exitOkay,
		},
		{
			name: "test exitcode for any error",
			err:  fmt.Errorf("an error"),
			code: exitError,
		},
		{
			name: "test exitcode for eof error",
			err:  io.EOF,
			code: exitIntrpt,
		},
		{
			name: "test exitcode for auth error",
			err:  cmdutil.ErrNoLogin,
			code: exitAuth,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			code := handleError(nc, c.err)
			if code != c.code {
				t.Errorf("handleError(): want %v, got %v", c.code, code)
			}
		})
	}
}
