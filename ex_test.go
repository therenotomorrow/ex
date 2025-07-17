package ex

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	const errText = "new error"

	err := New(errText)

	assert.IsType(t, new(XError), err)
	require.NoError(t, err.cause)
	assert.EqualError(t, err.label, errText)
}

func TestFrom(t *testing.T) {
	t.Parallel()

	t.Run("wrapping a standard error", func(t *testing.T) {
		t.Parallel()

		var (
			stdErr = errors.New("standard error")
			err    = From(stdErr)
		)

		assert.IsType(t, new(XError), err)
		require.ErrorIs(t, err.label, stdErr)
		require.NoError(t, err.cause)
	})

	t.Run("wrapping an existing XError", func(t *testing.T) {
		t.Parallel()

		var (
			causeErr     = errors.New("original cause")
			originalXErr = New("original").Because(causeErr)
			err          = From(originalXErr)
		)

		assert.IsType(t, new(XError), err)
		assert.Equal(t, err, originalXErr)
		require.ErrorIs(t, err.cause, causeErr)
		assert.NotSame(t, originalXErr, err)
	})
}

func TestUnexpected(t *testing.T) {
	t.Parallel()

	var (
		causeErr = errors.New("critical failure")
		err      = Unexpected(causeErr)
		xer      = new(XError)
	)

	require.ErrorAs(t, err, &xer)
	require.ErrorIs(t, xer.label, ErrUnexpected)
	require.ErrorIs(t, xer.cause, causeErr)
}

func TestMust(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			Must("all good", nil)
		})
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		var (
			causeErr = io.EOF
			err      = New("some error").Because(causeErr)
		)

		assert.PanicsWithValue(t, causeErr, func() {
			Must("value", err)
		})
	})
}

func TestMustDo(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			MustDo(nil)
		})
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		var (
			causeErr = io.EOF
			err      = New("some error").Because(causeErr)
		)

		assert.PanicsWithValue(t, causeErr, func() {
			MustDo(err)
		})
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
			name: "nested XError",
			args: args{err: New("level1").Because(New("level2").Because(rootCause))},
			want: rootCause,
		},
		{
			name: "XError with no cause",
			args: args{err: New("level1")},
			want: LError("level1"),
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
			name: "XError wrapping a standard error",
			args: args{err: From(rootCause)},
			want: rootCause,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := Cause(test.args.err)

			require.ErrorIs(t, got, test.want)
		})
	}
}

func TestLError(t *testing.T) {
	t.Parallel()

	t.Run("Because", func(t *testing.T) {
		t.Parallel()

		const constErr = LError("base error")

		var (
			causeErr = errors.New("the underlying cause")
			err      = constErr.Because(causeErr)
			xer      = new(XError)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, xer.label, constErr)
		require.ErrorIs(t, xer.cause, causeErr)
	})

	t.Run("Reason", func(t *testing.T) {
		t.Parallel()

		const (
			constErr   = LError("base error")
			reasonText = "a specific reason"
		)

		var (
			err = constErr.Reason(reasonText)
			xer = new(XError)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, xer.label, constErr)
		require.ErrorIs(t, xer.cause, LError(reasonText))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		const err = LError("test error")

		assert.Equal(t, "test error", err.Error())
	})
}

func TestXError(t *testing.T) { //nolint:funlen // don't want to separate XError.Is tests from Suite
	t.Parallel()

	const baseErr = LError("base error")

	var (
		causeErr = errors.New("root cause")
		xErr     = &XError{label: baseErr, cause: causeErr}
	)

	t.Run("Because", func(t *testing.T) {
		t.Parallel()

		var (
			newCause = errors.New("a different cause")
			err      = xErr.Because(newCause)
			xer      = new(XError)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, xer.label, xErr.label)
		require.ErrorIs(t, xer.cause, newCause)
	})

	t.Run("Reason", func(t *testing.T) {
		t.Parallel()

		var (
			reasonText = "a new reason"
			err        = xErr.Reason(reasonText)
			xer        = new(XError)
		)

		require.ErrorAs(t, err, &xer)
		require.ErrorIs(t, xer.label, xErr.label)
		require.ErrorIs(t, xer.cause, LError(reasonText))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "base error", xErr.Error())
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
			xer  *XError
			name string
			want bool
		}{
			{
				name: "target is the wrapped error",
				xer:  &XError{label: baseErr, cause: causeErr},
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "target is the cause",
				xer:  &XError{label: baseErr, cause: causeErr},
				args: args{target: causeErr},
				want: true,
			},
			{
				name: "target is a different error",
				xer:  &XError{label: baseErr, cause: causeErr},
				args: args{target: errAnother},
				want: false,
			},
			{
				name: "no cause, target matches wrapped error",
				xer:  &XError{label: baseErr, cause: nil},
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "no cause, target does not match",
				xer:  &XError{label: baseErr, cause: nil},
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
