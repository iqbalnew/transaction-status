package ftp

import (
	"context"
	"errors"
	"testing"

	ftpmock "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/ftp/mock"
	ftpwrapper "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/ftp/wrapper"

	"github.com/stretchr/testify/suite"
)

type FtpTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *FtpTestSuite) SetupSuite() {
	s.ctx = context.Background()
}

func (s *FtpTestSuite) SetupTest() {
	ftpmock.NewFtpMock(":9322")
}

func TestInitFtp(t *testing.T) {
	suite.Run(t, new(FtpTestSuite))
}

func (s *FtpTestSuite) TestFtp_NewFtp() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		config   Config
		expected expectation
	}{
		"Success": {
			config: Config{},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			ftpObj := NewFtp(tt.config, &ftpmock.FtpMock{}, nil)

			if ftpObj == nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", "not nil", ftpObj)
			}
		})
	}
}

func (s *FtpTestSuite) TestFtp_InitConnection() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		ftpObj   *FtpObj
		expected expectation
	}{
		"Success": {
			ftpObj: &FtpObj{
				fw: &ftpwrapper.FtpWrapper{},
				config: Config{
					Host:    ":9322",
					Timeout: 30,
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			ftpObj: &FtpObj{
				fw: &ftpwrapper.FtpWrapper{},
				config: Config{
					Host:    ":9999",
					Timeout: 30,
				},
			},
			expected: expectation{
				err: errors.New("dial tcp :9999: connectex: No connection could be made because the target machine actively refused it."),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.ftpObj.InitConnection()

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *FtpTestSuite) TestFtp_Connect() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		ftpObj   FtpObj
		expected expectation
	}{
		"Success": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
				config: Config{
					User:     "test",
					Password: "test",
				},
			},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.ftpObj.Connect()

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *FtpTestSuite) TestFtp_Close() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		ftpObj   FtpObj
		expected expectation
	}{
		"Success": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
		"Failed": {
			ftpObj: FtpObj{},
			expected: expectation{
				err: errors.New("connection already closed"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.ftpObj.Close()

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *FtpTestSuite) TestFtp_ListFiles() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		ftpObj   FtpObj
		path     string
		expected expectation
	}{
		"Success": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			expected: expectation{
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			entries, err := tt.ftpObj.ListFiles(tt.path)

			if err != nil && entries == nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *FtpTestSuite) TestFtp_ReadFile() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		ftpObj   FtpObj
		path     string
		filename string
		expected expectation
	}{
		"Success": {
			ftpObj: FtpObj{
				fw: &ftpwrapper.FtpWrapper{},
				config: Config{
					Host:     ":9322",
					User:     "coba",
					Password: "coba123",
				},
			},
			path:     "/ada/nih/",
			filename: "ada.csv",
			expected: expectation{
				err: nil,
			},
		},
		"FailedChangeDir": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			path:     "",
			filename: "",
			expected: expectation{
				err: errors.New("path not found"),
			},
		},
		"FailedReadFile": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			path:     "/ada/nih/",
			filename: "",
			expected: expectation{
				err: errors.New("file not found"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			tt.ftpObj.InitConnection()
			tt.ftpObj.Connect()
			_, err := tt.ftpObj.ReadFile(tt.path, tt.filename)

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}

func (s *FtpTestSuite) TestFtp_WriteFile() {
	type expectation struct {
		err error
	}

	tests := map[string]struct {
		ftpObj   FtpObj
		path     string
		filename string
		content  []byte
		expected expectation
	}{
		"Success": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			path:     "/path/to/file/",
			filename: "file.txt",
			content:  []byte("This is content"),
			expected: expectation{
				err: nil,
			},
		},
		"FailedChangeDir": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			path:     "",
			filename: "",
			content:  nil,
			expected: expectation{
				err: errors.New("path not found"),
			},
		},
		"FailedWriteData": {
			ftpObj: FtpObj{
				conn: &ftpmock.FtpMock{},
			},
			path:     "/ada/nih/",
			filename: "",
			content:  []byte("This is content"),
			expected: expectation{
				err: errors.New("failed write file"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			err := tt.ftpObj.WriteFile(tt.path, tt.filename, tt.content)

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			}
		})
	}
}
