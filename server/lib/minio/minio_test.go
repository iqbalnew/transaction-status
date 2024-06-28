package minio

import (
	"context"
	"errors"
	"testing"

	miniomock "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/minio/mock"
	miniowrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/minio/wrapper"

	"github.com/stretchr/testify/suite"
)

type MinioTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *MinioTestSuite) SetupSuite() {
	s.ctx = context.Background()
}

func TestInitMinio(t *testing.T) {
	suite.Run(t, new(MinioTestSuite))
}

func (s *MinioTestSuite) TestMinio_NewMinio() {
	type expectation struct {
		out *MinioObj
		err error
	}

	tests := map[string]struct {
		config   *Config
		mcw      miniowrapper.MinioInterface
		expected expectation
	}{
		"SuccessCreateNewMinio": {
			mcw:    &miniomock.MinioMock{},
			config: &Config{},
			expected: expectation{
				out: &MinioObj{
					config: &Config{},
				},
				err: nil,
			},
		},
		"FailedCreateNewMinio": {
			mcw:    nil,
			config: &Config{},
			expected: expectation{
				out: nil,
				err: errors.New("failed create minio object"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			_, err := NewMinio(tt.config, tt.mcw)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *MinioTestSuite) TestMinio_Connect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		minioObj *MinioObj
		expected expectation
	}{
		"Success": {
			minioObj: &MinioObj{
				config: &Config{
					Endpoint: "success",
				},
				mcw: &miniomock.MinioMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			minioObj: &MinioObj{
				config: &Config{
					Endpoint: "failed-connect",
				},
				mcw: &miniomock.MinioMock{},
			},
			expected: expectation{
				err: errors.New("failed to connect minio"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.minioObj.Connect()

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *MinioTestSuite) TestMinio_ReadFile() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		fileName string
		expected expectation
	}{
		"Failed": {
			fileName: "failed-object.xlsx",
			expected: expectation{
				err: errors.New("failed get object"),
			},
		},
	}

	minioObj := &MinioObj{
		config: &Config{
			Timeout: int64(10),
		},
		mcw: &miniomock.MinioMock{},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			_, err := minioObj.ReadFile(s.ctx, tt.fileName)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *MinioTestSuite) TestMinio_UploadFile() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		fileName string
		data     []byte
		expected expectation
	}{
		"Success": {
			fileName: "unit-test-minio.txt",
			data:     []byte{},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			fileName: "failed-object.xlsx",
			expected: expectation{
				err: errors.New("failed put object"),
			},
		},
	}

	minioObj := &MinioObj{
		config: &Config{
			Timeout: int64(10),
		},
		mcw: &miniomock.MinioMock{},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			_, err := minioObj.UploadFile(s.ctx, tt.fileName, tt.data)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *MinioTestSuite) TestMinio_DeleteFile() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		fileName string
		expected expectation
	}{
		"Success": {
			fileName: "unit-test-minio.txt",
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			fileName: "failed-object.xlsx",
			expected: expectation{
				err: errors.New("failed remove object"),
			},
		},
	}

	minioObj := &MinioObj{
		config: &Config{
			Timeout: int64(10),
		},
		mcw: &miniomock.MinioMock{},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := minioObj.DeleteFile(s.ctx, tt.fileName)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}
