package ex_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/ex"
)

func TestC(t *testing.T) {
	t.Parallel()

	const (
		text     = "const error"
		constErr = ex.C(text)
	)

	assert.IsType(t, constErr, ex.ConstError(""))
	assert.EqualError(t, constErr, text)
}

func TestNew(t *testing.T) {
	t.Parallel()

	const text = "new error"

	var (
		err        = ex.New(text)
		got, cause = ex.Test(err)
	)

	require.NoError(t, cause)

	assert.IsType(t, new(ex.ExtraError), err)
	assert.EqualError(t, got, text)
}

func TestFrom(t *testing.T) {
	t.Parallel()

	t.Run("wrapping a standard error", func(t *testing.T) {
		t.Parallel()

		var (
			stdErr     = errors.New("standard error")
			err        = ex.From(stdErr)
			got, cause = ex.Test(err)
		)

		require.ErrorIs(t, got, stdErr)
		require.NoError(t, cause)

		assert.IsType(t, new(ex.ExtraError), err)
	})

	t.Run("wrapping an existing ExtraError", func(t *testing.T) {
		t.Parallel()

		const constErr = ex.ConstError("original")

		var (
			causeErr   = errors.New("original cause")
			originErr  = constErr.Because(causeErr)
			err        = ex.From(originErr)
			got, cause = ex.Test(err)
		)

		require.ErrorIs(t, cause, causeErr)
		require.ErrorIs(t, got, constErr)

		assert.IsType(t, new(ex.ExtraError), err)
		assert.Equal(t, err, originErr)
		assert.NotSame(t, originErr, err)
	})
}

func TestUnexpected(t *testing.T) {
	t.Parallel()

	var (
		causeErr   = errors.New("critical failure")
		err        = ex.Unexpected(causeErr)
		xer        = new(ex.ExtraError)
		got, cause = ex.Test(err)
	)

	require.ErrorAs(t, err, &xer)
	require.ErrorIs(t, got, ex.ErrUnexpected)
	require.ErrorIs(t, cause, causeErr)

	assert.IsType(t, new(ex.ExtraError), err)
}

func TestMust(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() { ex.Must("all good", nil) })
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		var (
			causeErr = io.EOF
			err      = ex.New("some error").Because(causeErr)
		)

		assert.PanicsWithValue(t, causeErr, func() { ex.Must("value", err) })
	})
}

func TestMustDo(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() { ex.MustDo(nil) })
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		var (
			causeErr = io.EOF
			err      = ex.New("some error").Because(causeErr)
		)

		assert.PanicsWithValue(t, causeErr, func() { ex.MustDo(err) })
	})
}

func TestCause(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	rootCause := errors.New("the actual root cause")

	tests := []struct {
		args args
		want error
		name string
	}{
		{
			name: "nested ExtraError",
			args: args{err: ex.New("level1").Because(ex.New("level2").Because(rootCause))},
			want: rootCause,
		},
		{
			name: "ExtraError with no cause",
			args: args{err: ex.New("level1")},
			want: ex.ConstError("level1"),
		},
		{
			name: "standard error",
			args: args{err: rootCause},
			want: rootCause,
		},
		{
			name: "nil error",
			args: args{err: nil},
			want: nil,
		},
		{
			name: "ExtraError wrapping a standard error",
			args: args{err: ex.From(rootCause)},
			want: rootCause,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := ex.Cause(test.args.err)

			require.ErrorIs(t, got, test.want)
		})
	}
}

func TestConstError(t *testing.T) {
	t.Parallel()

	t.Run("Because", func(t *testing.T) {
		t.Parallel()

		const constErr = ex.ConstError("base error")

		var (
			causeErr   = errors.New("the underlying cause")
			err        = constErr.Because(causeErr)
			xer        = new(ex.ExtraError)
			got, cause = ex.Test(err)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, got, constErr)
		require.ErrorIs(t, cause, causeErr)
	})

	t.Run("Reason", func(t *testing.T) {
		t.Parallel()

		const (
			constErr = ex.ConstError("base error")
			text     = "a specific reason"
		)

		var (
			err        = constErr.Reason(text)
			xer        = new(ex.ExtraError)
			got, cause = ex.Test(err)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, got, constErr)
		require.ErrorIs(t, cause, ex.ConstError(text))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		const text = "test error"

		var err error = ex.ConstError(text)

		assert.Equal(t, text, err.Error())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		const (
			text = "test error"
			err  = ex.ConstError(text)
		)

		assert.Equal(t, text, err.String())
	})
}

func TestExtraError(t *testing.T) {
	t.Parallel()

	const baseErr = ex.ConstError("base error")

	var (
		causeErr = errors.New("root cause")
		xErr     = ex.From(ex.From(baseErr).Because(causeErr))
	)

	t.Run("Because", func(t *testing.T) {
		t.Parallel()

		var (
			newCause   = errors.New("a different cause")
			err        = xErr.Because(newCause)
			xer        = new(ex.ExtraError)
			got, cause = ex.Test(err)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, got, baseErr)
		require.ErrorIs(t, cause, newCause)
	})

	t.Run("Reason", func(t *testing.T) {
		t.Parallel()

		var (
			reasonText = "a new reason"
			err        = xErr.Reason(reasonText)
			xer        = new(ex.ExtraError)
			got, cause = ex.Test(err)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, got, baseErr)
		require.ErrorIs(t, cause, ex.ConstError(reasonText))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		var (
			emptyErr     error = new(ex.ExtraError)
			onlyErr      error = ex.From(baseErr)
			rootErr      error = ex.ConstError("root error")
			deepCauseErr error = ex.From(xErr.Because(ex.From(ex.ConstError("something wrong")).Because(rootErr)))
		)

		assert.Empty(t, emptyErr.Error())
		assert.Equal(t, "base error", onlyErr.Error())
		assert.Equal(t, "base error (root cause)", xErr.Error())
		assert.Equal(t, `base error (root error)`, deepCauseErr.Error())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		var (
			emptyErr     = new(ex.ExtraError)
			onlyErr      = ex.From(baseErr)
			rootErr      = ex.ConstError("root error")
			deepCauseErr = ex.From(xErr.Because(ex.From(ex.ConstError("something wrong")).Because(rootErr)))
		)

		assert.JSONEq(t, `{}`, emptyErr.String())
		assert.JSONEq(t, `{"error":"base error"}`, onlyErr.String())
		assert.JSONEq(t, `{"cause":"root cause","error":"base error"}`, xErr.String())
		assert.JSONEq(t, `{"cause":"something wrong (root error)","error":"base error"}`, deepCauseErr.String())
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, baseErr, xErr.Unwrap())
	})

	t.Run("Is", func(t *testing.T) {
		t.Parallel()

		type args struct {
			target error
		}

		errAnother := errors.New("another")

		tests := []struct {
			args args
			xer  *ex.ExtraError
			name string
			want bool
		}{
			{
				name: "target is the wrapped error",
				xer:  ex.From(baseErr.Because(causeErr)),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "target is the cause",
				xer:  ex.From(baseErr.Because(causeErr)),
				args: args{target: causeErr},
				want: true,
			},
			{
				name: "target is a different error",
				xer:  ex.From(baseErr.Because(causeErr)),
				args: args{target: errAnother},
				want: false,
			},
			{
				name: "no cause, target matches wrapped error",
				xer:  ex.From(baseErr),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "no cause, target does not match",
				xer:  ex.From(baseErr),
				args: args{target: errAnother},
				want: false,
			},
		}

		for _, test := range tests {
			test := test

			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := test.xer.Is(test.args.target)

				assert.Equal(t, test.want, got)
			})
		}
	})
}

func TestTest(t *testing.T) {
	t.Parallel()

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			const (
				text     = "base text"
				constErr = ex.ConstError("some error")
			)

			var (
				err        = ex.New(text).Because(constErr)
				got, cause = ex.Test(err)
			)

			require.ErrorIs(t, got, ex.ConstError(text))
			require.ErrorIs(t, cause, constErr)
		})
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		const text = "invalid error type"

		assert.PanicsWithValue(t, text, func() { _, _ = ex.Test(nil) })
	})
}
