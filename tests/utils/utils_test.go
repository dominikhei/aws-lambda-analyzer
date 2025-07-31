package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			name: "profile only",
			opts: sdktypes.ConfigOptions{
				Profile: "my-profile",
			},
			wantErr: false,
		},
		{
			name: "region only",
			opts: sdktypes.ConfigOptions{
				Region: "us-west-1",
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
			wantMin:    1,
			wantMax:    20,
			wantMean:   10.5,
			wantMedian: 10,
			expectP95:  true,
			expectP99:  false,
			expectConf: false,
		},
		{
			name:       "large slice with p99 and confInt",
			input:      generateSlice(100),
			wantErr:    false,
			wantMin:    1,
			wantMax:    100,
			wantMean:   50.5,
			wantMedian: 50,
			expectP95:  true,
			expectP99:  true,
			expectConf: true,
		},
		{
			name:       "slice with confInt only",
			input:      generateSlice(30),
			wantErr:    false,
			wantMin:    1,
			wantMax:    30,
			wantMean:   15.5,
			wantMedian: 15,
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

// Helper function to generate the slices used as input.
func generateSlice(n int) []float64 {
	s := make([]float64, n)
	for i := 0; i < n; i++ {
		s[i] = float64(i + 1)
	}
	return s
}

// This mock client mocks the actual lambda client and matches the client interface defined in interfaces.
type MockLambdaClient struct {
	mock.Mock
}

func (m *MockLambdaClient) GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*lambda.GetFunctionOutput), args.Error(1)
}

func TestFunctionExists(t *testing.T) {
	tests := []struct {
		name         string
		functionName string
		setupMock    func(*MockLambdaClient)
		want         bool
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "function exists",
			functionName: "test-function",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function"
				})).Return(&lambda.GetFunctionOutput{}, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:         "function does not exist",
			functionName: "nonexistent-function",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "nonexistent-function"
				})).Return(nil, &types.ResourceNotFoundException{
					Type:    aws.String("User"),
					Message: aws.String("Function not found"),
				})
			},
			want:    false,
			wantErr: false,
		},
		{
			name:         "access denied error",
			functionName: "restricted-function",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "restricted-function"
				})).Return(nil, errors.New("AccessDeniedException: User is not authorized"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "AccessDeniedException: User is not authorized",
		},
		{
			name:         "generic error",
			functionName: "error-function",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "error-function"
				})).Return(nil, errors.New("internal server error"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "internal server error",
		},
		{
			name:         "empty function name",
			functionName: "",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == ""
				})).Return(nil, errors.New("ValidationException: Function name cannot be empty"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "ValidationException: Function name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockLambdaClient)
			tt.setupMock(mockClient)
			ctx := context.Background()

			got, err := utils.FunctionExists(ctx, mockClient, tt.functionName)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestQualifierExists(t *testing.T) {
	tests := []struct {
		name         string
		functionName string
		qualifier    string
		setupMock    func(*MockLambdaClient)
		want         bool
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "version qualifier exists",
			functionName: "test-function",
			qualifier:    "1",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function" && *input.Qualifier == "1"
				})).Return(&lambda.GetFunctionOutput{}, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:         "alias qualifier exists",
			functionName: "test-function",
			qualifier:    "PROD",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function" && *input.Qualifier == "PROD"
				})).Return(&lambda.GetFunctionOutput{}, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:         "qualifier does not exist",
			functionName: "test-function",
			qualifier:    "999",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function" && *input.Qualifier == "999"
				})).Return(nil, &types.ResourceNotFoundException{
					Type:    aws.String("User"),
					Message: aws.String("The resource you requested does not exist."),
				})
			},
			want:    false,
			wantErr: false,
		},
		{
			name:         "function does not exist",
			functionName: "nonexistent-function",
			qualifier:    "1",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "nonexistent-function" && *input.Qualifier == "1"
				})).Return(nil, &types.ResourceNotFoundException{
					Type:    aws.String("User"),
					Message: aws.String("Function not found: arn:aws:lambda:us-east-1:123456789012:function:nonexistent-function:1"),
				})
			},
			want:    false,
			wantErr: false,
		},
		{
			name:         "access denied error",
			functionName: "restricted-function",
			qualifier:    "PROD",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "restricted-function" && *input.Qualifier == "PROD"
				})).Return(nil, errors.New("AccessDeniedException: User is not authorized"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "AccessDeniedException: User is not authorized",
		},
		{
			name:         "invalid qualifier format",
			functionName: "test-function",
			qualifier:    "invalid-qualifier!",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function" && *input.Qualifier == "invalid-qualifier!"
				})).Return(nil, errors.New("ValidationException: 1 validation error detected"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "ValidationException: 1 validation error detected",
		},
		{
			name:         "empty function name",
			functionName: "",
			qualifier:    "1",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "" && *input.Qualifier == "1"
				})).Return(nil, errors.New("ValidationException: Function name cannot be empty"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "ValidationException: Function name cannot be empty",
		},
		{
			name:         "$LATEST qualifier",
			functionName: "test-function",
			qualifier:    "$LATEST",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function" && *input.Qualifier == "$LATEST"
				})).Return(&lambda.GetFunctionOutput{}, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:         "generic error",
			functionName: "test-function",
			qualifier:    "1",
			setupMock: func(m *MockLambdaClient) {
				m.On("GetFunction", mock.Anything, mock.MatchedBy(func(input *lambda.GetFunctionInput) bool {
					return *input.FunctionName == "test-function" && *input.Qualifier == "1"
				})).Return(nil, errors.New("internal server error"))
			},
			want:        false,
			wantErr:     true,
			expectedErr: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockLambdaClient)
			tt.setupMock(mockClient)
			ctx := context.Background()

			got, err := utils.QualifierExists(ctx, mockClient, tt.functionName, tt.qualifier)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			mockClient.AssertExpectations(t)
		})
	}
}
