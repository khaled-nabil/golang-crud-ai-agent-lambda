package ollama

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat64To32(t *testing.T) {
	t.Run("converts values correctly", func(t *testing.T) {
		input := []float64{0.1, 0.5, 1.0, -0.3}
		result := Float64To32(input)

		require.Len(t, result, len(input))
		for i, v := range input {
			require.InDelta(t, v, float64(result[i]), 1e-6)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Float64To32([]float64{})
		require.Empty(t, result)
	})

	t.Run("preserves length", func(t *testing.T) {
		input := make([]float64, 100)
		for i := range input {
			input[i] = float64(i) * 0.01
		}

		result := Float64To32(input)

		require.Len(t, result, 100)
	})
}
