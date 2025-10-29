package ex_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/ex"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("nillable error", func(t *testing.T) {
		t.Parallel()

		err := ex.New(nil)

		require.NoError(t, err)
	})

	t.Run("standard error", func(t *testing.T) {
		t.Parallel()

		var (
			stdErr     = errors.New("standard error")
			err        = ex.New(stdErr)
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
			err        = ex.New(packageErr)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, cause, causeErr)
		require.ErrorIs(t, got, constErr)
		require.Equal(t, err, packageErr)
		require.NotSame(t, packageErr, err)
	})
}

func TestExpose(t *testing.T) {
	t.Parallel()

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			const (
				baseErr  = ex.Error("base text")
				constErr = ex.Error("some error")
			)

			var (
				err        = ex.New(baseErr).Because(constErr)
				got, cause = ex.Expose(err)
			)

			require.ErrorIs(t, got, baseErr)
			require.ErrorIs(t, cause, constErr)
		})
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		const text = "invalid error type"

		require.PanicsWithValue(t, text, func() {
			_, _ = ex.Expose(nil)
		})
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
	require.EqualError(t, err, "unexpected: unexpected failure")
	require.NoError(t, ex.Unexpected(nil))
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
	require.EqualError(t, err, "critical: critical failure")
	require.NoError(t, ex.Critical(nil))
}

func TestDummy(t *testing.T) {
	t.Parallel()

	var (
		causeErr   = errors.New("dummy failure")
		err        = ex.Dummy(causeErr)
		got, cause = ex.Expose(err)
	)

	require.ErrorIs(t, got, ex.ErrDummy)
	require.ErrorIs(t, cause, causeErr)
	require.EqualError(t, err, "dummy: dummy failure")
	require.NoError(t, ex.Dummy(nil))
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

		require.EqualError(t, err, text)
	})
}

func TestXError(t *testing.T) {
	t.Parallel()

	const baseErr = ex.Error("base error")

	var (
		causeErr = errors.New("root cause")
		xErr     = ex.New(ex.New(baseErr).Because(causeErr))
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
			onlyErr      error = ex.New(baseErr)
			rootErr      error = ex.Error("root error")
			deepCauseErr error = ex.New(xErr.Because(ex.New(ex.Error("something wrong")).Because(rootErr)))
		)

		require.EqualError(t, onlyErr, "base error")
		require.EqualError(t, xErr, "base error: root cause")
		require.EqualError(t, deepCauseErr, "base error: something wrong: root error")
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, baseErr, errors.Unwrap(xErr))
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
				xer:  ex.New(baseErr.Because(causeErr)),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "target is the cause",
				xer:  ex.New(baseErr.Because(causeErr)),
				args: args{target: causeErr},
				want: true,
			},
			{
				name: "target is a different error",
				xer:  ex.New(baseErr.Because(causeErr)),
				args: args{target: errAnother},
				want: false,
			},
			{
				name: "no cause, target matches wrapped error",
				xer:  ex.New(baseErr),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "no cause, target does not match",
				xer:  ex.New(baseErr),
				args: args{target: errAnother},
				want: false,
			},
		}

		for _, test := range tests {
			test := test

			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := errors.Is(test.xer, test.args.target)

				require.Equal(t, test.want, got)
			})
		}
	})
}
