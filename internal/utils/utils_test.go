package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dominikhei/serverless-statistics/internal/utils"
	sdktypes "github.com/dominikhei/serverless-statistics/types"
)

func TestToLoadOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    sdktypes.ConfigOptions
		wantErr bool
	}{
		{
			name: "profile and region only",
			opts: sdktypes.ConfigOptions{
				Profile: "my-profile",
				Region:  "us-west-1",
			},
			wantErr: false,
		},
		{
			name: "with credentials",
			opts: sdktypes.ConfigOptions{
				AccessKeyID:     "AKIA...",
				SecretAccessKey: "SECRET...",
			},
			wantErr: false,
		},
		{
			name: "missing secret",
			opts: sdktypes.ConfigOptions{
				AccessKeyID: "AKIA...",
			},
			wantErr: true,
		},
		{
			name: "missing key",
			opts: sdktypes.ConfigOptions{
				SecretAccessKey: "SECRET...",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := utils.ToLoadOptions(tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, opts)
			} else {
				require.NoError(t, err)
				require.NotNil(t, opts)
			}
		})
	}
}

func TestCalcSummaryStats(t *testing.T) {
	tests := []struct {
		name       string
		input      []float64
		wantErr    bool
		wantMin    float64
		wantMax    float64
		wantMean   float64
		wantMedian float64
		expectP95  bool
		expectP99  bool
		expectConf bool
	}{
		{
			name:    "empty slice",
			input:   []float64{},
			wantErr: true,
		},
		{
			name:       "small slice no percentiles/confInt",
			input:      []float64{1, 2, 3, 4, 5},
			wantErr:    false,
			wantMin:    1,
			wantMax:    5,
			wantMean:   3,
			wantMedian: 3,
			expectP95:  false,
			expectP99:  false,
			expectConf: false,
		},
		{
			name:       "medium slice with p95",
			input:      generateSlice(20),
			wantErr:    false,
			expectP95:  true,
			expectP99:  false,
			expectConf: false,
		},
		{
			name:       "large slice with p99 and confInt",
			input:      generateSlice(100),
			wantErr:    false,
			expectP95:  true,
			expectP99:  true,
			expectConf: true,
		},
		{
			name:       "slice with confInt only",
			input:      generateSlice(30),
			wantErr:    false,
			expectP95:  true,
			expectP99:  false,
			expectConf: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.CalcSummaryStats(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantMin, got.Min)
			require.Equal(t, tt.wantMax, got.Max)
			require.InDelta(t, tt.wantMean, got.Mean, 0.0001)
			require.InDelta(t, tt.wantMedian, got.Median, 0.0001)

			if tt.expectP95 {
				require.NotNil(t, got.P95)
			} else {
				require.Nil(t, got.P95)
			}

			if tt.expectP99 {
				require.NotNil(t, got.P99)
			} else {
				require.Nil(t, got.P99)
			}

			if tt.expectConf {
				require.NotNil(t, got.ConfInt95)
			} else {
				require.Nil(t, got.ConfInt95)
			}
		})
	}
}

func generateSlice(n int) []float64 {
	s := make([]float64, n)
	for i := 0; i < n; i++ {
		s[i] = float64(i + 1)
	}
	return s
}
