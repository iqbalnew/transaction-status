package aes

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AesTestSuite struct {
	suite.Suite
	ctx       context.Context
	customAes *CustomAES
}

func (s *AesTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.customAes = NewCustomAES("unit-test-123")
}

func TestInitAes(t *testing.T) {
	suite.Run(t, new(AesTestSuite))
}

func (s *AesTestSuite) TestAes_NewCustomAES() {
	type expectation struct {
		out *CustomAES
	}

	tests := map[string]struct {
		passphrase string
		expected   expectation
	}{
		"Success": {
			passphrase: "unit-test-123",
			expected: expectation{
				out: &CustomAES{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			out := NewCustomAES(tt.passphrase)

			if out == nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
			}
		})
	}
}

func (s *AesTestSuite) TestAes_HashIDEncode() {
	type expectation struct {
		out string
		err error
	}

	tests := map[string]struct {
		hashId   string
		expected expectation
	}{
		"Success": {
			hashId: "1001",
			expected: expectation{
				out: "pln",
				err: nil,
			},
		},
		"Failed": {
			hashId: "failed",
			expected: expectation{
				out: "",
				err: errors.New("strconv.Atoi: parsing \"failed\": invalid syntax"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			out, err := s.customAes.HashIDEncode(tt.hashId)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			} else {
				if out != tt.expected.out {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
				}
			}
		})
	}
}

func (s *AesTestSuite) TestAes_HashIDDecode() {
	type expectation struct {
		out string
		err error
	}

	tests := map[string]struct {
		hash     string
		expected expectation
	}{
		"Success": {
			hash: "pln",
			expected: expectation{
				out: "1001",
				err: nil,
			},
		},
		"ZeroLength": {
			hash: "",
			expected: expectation{
				out: "0",
				err: nil,
			},
		},
		"Failed": {
			hash: "failed",
			expected: expectation{
				out: "",
				err: errors.New("mismatch between encode and decode: failed start XOuzB7 re-encoded. result: [4 81314]"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			out, err := s.customAes.HashIDDecode(tt.hash)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			} else {
				if out != tt.expected.out {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
				}
			}
		})
	}
}

func (s *AesTestSuite) TestAes_Encrypt() {
	type expectation struct {
		out string
		err error
	}

	tests := map[string]struct {
		passphrase string
		text       string
		expected   expectation
	}{
		"Success": {
			passphrase: "unit-test-123456",
			text:       "unit-test-success",
			expected: expectation{
				out: "cannot predict because encrypt random value",
				err: nil,
			},
		},
		"Failed": {
			passphrase: "",
			text:       "failed",
			expected: expectation{
				out: "",
				err: errors.New("crypto/aes: invalid key size 0"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			customAes := NewCustomAES(tt.passphrase)
			_, err := customAes.Encrypt(tt.text)

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

func (s *AesTestSuite) TestAes_Decrypt() {
	type expectation struct {
		out string
		err error
	}

	tests := map[string]struct {
		passphrase string
		text       string
		expected   expectation
	}{
		"Success": {
			passphrase: "unit-test-123456",
			text:       "17be333fe4f80c16484d8afe-87d7cfdaefeabd71445db795b00bb71d1e07d0502c9478648a29c2144d09134256",
			expected: expectation{
				out: "unit-test-success",
				err: nil,
			},
		},
		"FailedEmptyText": {
			passphrase: "unit-test-123456",
			text:       "",
			expected: expectation{
				out: "",
				err: nil,
			},
		},
		"FailedFirstKey": {
			passphrase: "unit-test-123456",
			text:       "failed-87d7cfdaefeabd71445db795b00bb71d1e07d0502c9478648a29c2144d09134256",
			expected: expectation{
				out: "",
				err: errors.New("encoding/hex: invalid byte: U+0069 'i'"),
			},
		},
		"FailedSecondKey": {
			passphrase: "unit-test-123456",
			text:       "17be333fe4f80c16484d8afe-failed",
			expected: expectation{
				out: "",
				err: errors.New("encoding/hex: invalid byte: U+0069 'i'"),
			},
		},
		"FailedEmptyKey": {
			passphrase: "",
			text:       "17be333fe4f80c16484d8afe-87d7cfdaefeabd71445db795b00bb71d1e07d0502c9478648a29c2144d09134256",
			expected: expectation{
				out: "",
				err: errors.New("crypto/aes: invalid key size 0"),
			},
		},
		"Failed": {
			passphrase: "unit-test-123465",
			text:       "17be333fe4f80c16484d8afe-87d7cfdaefeabd71445db795b00bb71d1e07d0502c9478648a29c2144d09134256",
			expected: expectation{
				out: "",
				err: errors.New("cipher: message authentication failed"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			customAes := NewCustomAES(tt.passphrase)
			out, err := customAes.Decrypt(tt.text)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			} else {
				if out != tt.expected.out {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
				}
			}
		})
	}
}
