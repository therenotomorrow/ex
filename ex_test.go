package ex_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/ex"
)

func TestConv(t *testing.T) {
	t.Parallel()

	t.Run("nillable error", func(t *testing.T) {
		t.Parallel()

		err := ex.Conv(nil)

		require.NoError(t, err)
	})

	t.Run("standard error", func(t *testing.T) {
		t.Parallel()

		var (
			stdErr     = errors.New("standard error")
			err        = ex.Conv(stdErr)
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
			err        = ex.Conv(packageErr)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, cause, causeErr)
		require.ErrorIs(t, got, constErr)
		require.Equal(t, err, packageErr)
		require.NotSame(t, packageErr, err)
	})

	t.Run("nil cause", func(t *testing.T) {
		t.Parallel()

		const baseErr = ex.Error("base error")

		var (
			xErr       = ex.Conv(baseErr)
			wrapped    = ex.Conv(xErr)
			got, cause = ex.Expose(wrapped)
		)

		require.NotNil(t, xErr)
		require.NotNil(t, wrapped)
		require.ErrorIs(t, got, baseErr)
		require.NoError(t, cause)
	})

	t.Run("wrapped xerror", func(t *testing.T) {
		t.Parallel()

		const (
			baseErr  = ex.Error("base")
			causeErr = ex.Error("cause")
		)

		var (
			original   = ex.Conv(baseErr).Because(causeErr)
			wrapped    = ex.Conv(original)
			got, cause = ex.Expose(wrapped)
		)

		require.NotSame(t, original, wrapped)
		require.ErrorIs(t, got, baseErr)
		require.ErrorIs(t, cause, causeErr)
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("empty error", func(t *testing.T) {
		t.Parallel()

		var (
			err        = ex.New("")
			got, cause = ex.Expose(err)
		)

		require.NoError(t, got)
		require.NoError(t, cause)
	})

	t.Run("usual error", func(t *testing.T) {
		t.Parallel()

		const constErr = "something went wrong"

		var (
			err        = ex.New(constErr)
			got, cause = ex.Expose(err)
		)

		require.NoError(t, cause)
		require.ErrorIs(t, got, ex.Error(constErr))
	})
}

func TestPanic(t *testing.T) {
	t.Parallel()

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			ex.Panic(nil)
		})
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		const text = "critical: super fail"

		err := errors.New("super fail")

		require.PanicsWithError(t, text, func() {
			ex.Panic(err)
		})
	})
}

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("no skip", func(t *testing.T) {
		t.Parallel()

		ex.Skip(nil)
	})

	t.Run("with skip", func(t *testing.T) {
		t.Parallel()

		err := errors.New("super fail")

		ex.Skip(err) // nothing happens here - it's just a mark
	})
}

func TestExpose(t *testing.T) {
	t.Parallel()

	t.Run("xerror passed", func(t *testing.T) {
		t.Parallel()

		const (
			baseErr  = ex.Error("base text")
			constErr = ex.Error("some error")
		)

		var (
			err        = ex.Conv(baseErr).Because(constErr)
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, got, baseErr)
		require.ErrorIs(t, cause, constErr)
	})

	t.Run("nothing passed", func(t *testing.T) {
		t.Parallel()

		got, cause := ex.Expose(nil)

		require.NoError(t, got)
		require.NoError(t, cause)
	})

	t.Run("standard passed", func(t *testing.T) {
		t.Parallel()

		var (
			err        = errors.New("not an ex")
			got, cause = ex.Expose(err)
		)

		require.ErrorIs(t, got, err)
		require.NoError(t, cause)
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

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			require.NoError(t, ex.Critical(nil))
		})
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		const text = "critical: critical failure"

		causeErr := errors.New("critical failure")

		require.PanicsWithError(t, text, func() {
			_ = ex.Critical(causeErr)
		})
	})
}

func TestUnknown(t *testing.T) {
	t.Parallel()

	var (
		causeErr   = errors.New("unknown failure")
		err        = ex.Unknown(causeErr)
		got, cause = ex.Expose(err)
	)

	require.ErrorIs(t, got, ex.ErrUnknown)
	require.ErrorIs(t, cause, causeErr)
	require.EqualError(t, err, "unknown: unknown failure")
	require.NoError(t, ex.Unknown(nil))
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
		xErr     = ex.Conv(ex.Conv(baseErr).Because(causeErr))
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
			onlyErr      error = ex.Conv(baseErr)
			rootErr      error = ex.Error("root error")
			deepCauseErr error = ex.Conv(xErr.Because(ex.Conv(ex.Error("something wrong")).Because(rootErr)))
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
				xer:  ex.Conv(baseErr.Because(causeErr)),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "target is the cause",
				xer:  ex.Conv(baseErr.Because(causeErr)),
				args: args{target: causeErr},
				want: true,
			},
			{
				name: "target is a different error",
				xer:  ex.Conv(baseErr.Because(causeErr)),
				args: args{target: errAnother},
				want: false,
			},
			{
				name: "no cause, target matches wrapped error",
				xer:  ex.Conv(baseErr),
				args: args{target: baseErr},
				want: true,
			},
			{
				name: "no cause, target does not match",
				xer:  ex.Conv(baseErr),
				args: args{target: errAnother},
				want: false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				got := errors.Is(test.xer, test.args.target)

				require.Equal(t, test.want, got)
			})
		}
	})
}

func TestDeepErrorChain(t *testing.T) {
	t.Parallel()

	const (
		depth   = 15
		rootErr = ex.Error("root")
	)

	var err error = ex.Error("root")
	// build a deep chain of errors
	for i := 1; i <= depth; i++ {
		level := "level " + string(rune('a'+i))
		err = ex.Conv(ex.Error(level)).Because(err)
	}

	require.EqualError(t, err, ""+
		"level p: "+
		"level o: "+
		"level n: "+
		"level m: "+
		"level l: "+
		"level k: "+
		"level j: "+
		"level i: "+
		"level h: "+
		"level g: "+
		"level f: "+
		"level e: "+
		"level d: "+
		"level c: "+
		"level b: "+
		"root")

	require.ErrorIs(t, err, rootErr)

	unwrapped := errors.Unwrap(err)

	require.Error(t, unwrapped)
	require.NotEqual(t, err, unwrapped)
}

func TestErrorChainWithMixedTypes(t *testing.T) {
	t.Parallel()

	var (
		stdErr1 = errors.New("standard error 1")
		exErr1  = ex.Error("ex error 1")
		stdErr2 = errors.New("standard error 2")
		exErr2  = ex.Error("ex error 2")
	)

	err := ex.Conv(exErr2).Because(
		ex.Conv(stdErr2).Because(
			ex.Conv(exErr1).Because(stdErr1),
		),
	)

	// finds all errors in the chain
	require.ErrorIs(t, err, exErr2)
	require.ErrorIs(t, err, stdErr2)
	require.ErrorIs(t, err, exErr1)
	require.ErrorIs(t, err, stdErr1)

	require.EqualError(t, err, "ex error 2: standard error 2: ex error 1: standard error 1")
}
