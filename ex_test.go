package ex_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/ex"
)

func TestNew(t *testing.T) {
	t.Parallel()

	const text = "new error"

	var (
		err        = ex.New(text)
		got, cause = ex.Expose(err)
	)

	require.NoError(t, cause)
	require.ErrorIs(t, got, ex.Error(text))
}

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("nillable error", func(t *testing.T) {
		t.Parallel()

		err := ex.Cast(nil)

		require.NoError(t, err)
	})

	t.Run("standard error", func(t *testing.T) {
		t.Parallel()

		var (
			stdErr     = errors.New("standard error")
			err        = ex.Cast(stdErr)
			got, cause = ex.Expose(err)
		)

		require.NoError(t, cause)
		require.ErrorIs(t, got, stdErr)
	})

	t.Run("package error", func(t *testing.T) {
		t.Parallel()

		const constErr = ex.Error("original")

		var (
			causeErr   = errors.New("original cause")
			packageErr = constErr.Because(causeErr)
			err        = ex.Cast(packageErr)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, cause, causeErr)
		require.ErrorIs(t, got, constErr)

		assert.Equal(t, err, packageErr)
		assert.NotSame(t, packageErr, err)
	})
}

func TestExpose(t *testing.T) {
	t.Parallel()

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			const (
				text     = "base text"
				constErr = ex.Error("some error")
			)

			var (
				err        = ex.New(text).Because(constErr)
				got, cause = ex.Expose(err)
			)

			require.ErrorIs(t, got, ex.Error(text))
			require.ErrorIs(t, cause, constErr)
		})
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		const text = "invalid error type"

		assert.PanicsWithValue(t, text, func() { _, _ = ex.Expose(nil) })
	})
}

func TestMust(t *testing.T) {
	t.Parallel()

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			const text = "some error"

			got := ex.Must(text, nil)

			assert.Equal(t, text, got)
		})
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		var (
			causeErr = io.EOF
			err      = ex.New("some error").Because(causeErr)
		)

		assert.PanicsWithValue(t, err, func() { _ = ex.Must("value", err) })
	})
}

func TestUnexpected(t *testing.T) {
	t.Parallel()

	var (
		causeErr   = errors.New("unexpected failure")
		err        = ex.Unexpected(causeErr)
		got, cause = ex.Expose(err)
	)

	require.ErrorIs(t, got, ex.ErrUnexpected)
	require.ErrorIs(t, cause, causeErr)
}

func TestCritical(t *testing.T) {
	t.Parallel()

	var (
		causeErr   = errors.New("critical failure")
		err        = ex.Critical(causeErr)
		got, cause = ex.Expose(err)
	)

	require.ErrorIs(t, got, ex.ErrCritical)
	require.ErrorIs(t, cause, causeErr)
}

func TestDummy(t *testing.T) {
	t.Parallel()

	var (
		causeErr   = errors.New("critical failure")
		err        = ex.Dummy(causeErr)
		got, cause = ex.Expose(err)
	)

	require.ErrorIs(t, got, ex.ErrDummy)
	require.ErrorIs(t, cause, causeErr)
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("Because", func(t *testing.T) {
		t.Parallel()

		const constErr = ex.Error("base error")

		var (
			causeErr   = errors.New("the underlying cause")
			err        = constErr.Because(causeErr)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, got, constErr)
		require.ErrorIs(t, cause, causeErr)
	})

	t.Run("Reason", func(t *testing.T) {
		t.Parallel()

		const (
			constErr = ex.Error("base error")
			text     = "a specific reason"
		)

		var (
			err        = constErr.Reason(text)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, got, constErr)
		require.ErrorIs(t, cause, ex.Error(text))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		const text = "test error"

		var err error = ex.Error(text)

		assert.Equal(t, text, err.Error())
	})
}

func TestXError(t *testing.T) {
	t.Parallel()

	const baseErr = ex.Error("base error")

	var (
		causeErr = errors.New("root cause")
		xErr     = ex.Cast(ex.Cast(baseErr).Because(causeErr))
	)

	t.Run("Because", func(t *testing.T) {
		t.Parallel()

		var (
			newCause   = errors.New("a different cause")
			err        = xErr.Because(newCause)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, got, baseErr)
		require.ErrorIs(t, cause, newCause)
	})

	t.Run("Reason", func(t *testing.T) {
		t.Parallel()

		var (
			reasonText = "a new reason"
			err        = xErr.Reason(reasonText)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, got, baseErr)
		require.ErrorIs(t, cause, ex.Error(reasonText))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		var (
			onlyErr      error = ex.Cast(baseErr)
			rootErr      error = ex.Error("root error")
			deepCauseErr error = ex.Cast(xErr.Because(ex.Cast(ex.Error("something wrong")).Because(rootErr)))
		)

		assert.Equal(t, "base error", onlyErr.Error())
		assert.Equal(t, "base error: root cause", xErr.Error())
		assert.Equal(t, `base error: something wrong: root error`, deepCauseErr.Error())
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, baseErr, errors.Unwrap(xErr))
	})

	t.Run("Is", func(t *testing.T) {
		t.Parallel()

		type args struct {
			target error
		}

		errAnother := errors.New("another")

		tests := []struct {
			args args
			xer  ex.XError
			name string
			want bool
		}{
			{
				name: "target is the wrapped error",
				xer:  ex.Cast(baseErr.Because(causeErr)),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "target is the cause",
				xer:  ex.Cast(baseErr.Because(causeErr)),
				args: args{target: causeErr},
				want: true,
			},
			{
				name: "target is a different error",
				xer:  ex.Cast(baseErr.Because(causeErr)),
				args: args{target: errAnother},
				want: false,
			},
			{
				name: "no cause, target matches wrapped error",
				xer:  ex.Cast(baseErr),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "no cause, target does not match",
				xer:  ex.Cast(baseErr),
				args: args{target: errAnother},
				want: false,
			},
		}

		for _, test := range tests {
			test := test

			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := errors.Is(test.xer, test.args.target)

				assert.Equal(t, test.want, got)
			})
		}
	})
}
