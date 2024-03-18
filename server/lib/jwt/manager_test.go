package jwtmanager

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
)

type JwtManagerTestSuite struct {
	suite.Suite
	ctx    context.Context
	logger *logrus.Logger
}

func (s *JwtManagerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.logger = logrus.New()
	s.logger.SetLevel(logrus.DebugLevel)
}

func TestInitJwtManager(t *testing.T) {
	suite.Run(t, new(JwtManagerTestSuite))
}

func (s *JwtManagerTestSuite) TestJwtManager_NewJWTManager() {
	type expectation struct {
		out *JWTManager
	}

	tests := map[string]struct {
		expected expectation
	}{
		"Success": {
			expected: expectation{
				out: &JWTManager{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			out := NewJWTManager("", "", 1*time.Hour, nil, nil)

			if out == nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
			}
		})
	}
}

func (s *JwtManagerTestSuite) TestJwtManager_Generate() {
	type expectation struct {
		out string
	}

	tests := map[string]struct {
		expected expectation
	}{
		"Success": {
			expected: expectation{
				out: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			manager := &JWTManager{
				tokenDuration: 1 * time.Hour,
				secretKey:     "unit-test-123456",
			}
			out, err := manager.Generate()

			if err != nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", nil, err)
			} else {
				splitOut := strings.Split(out, ".")

				if tt.expected.out != splitOut[0] {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, splitOut[0])
				}
			}
		})
	}
}

func (s *JwtManagerTestSuite) TestJwtManager_Verify() {
	type expectation struct {
		out *UserClaims
		err error
	}

	tests := map[string]struct {
		tokenAccess string
		expected    expectation
	}{
		"Success": {
			expected: expectation{
				out: &UserClaims{},
				err: nil,
			},
		},
		"Failed": {
			tokenAccess: "neUQdspTBcXAXntqHlay4gqkwy4hHKq4f8sV5S3b6Ko",
			expected: expectation{
				out: nil,
				err: errors.New("invalid token: token contains an invalid number of segments"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			manager := &JWTManager{
				tokenDuration: 1 * time.Hour,
				secretKey:     "unit-test-123456",
			}

			if tt.tokenAccess == "" {
				tt.tokenAccess, _ = manager.Generate()
			}

			out, err := manager.Verify(tt.tokenAccess)

			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				} else if tt.expected.err.Error() != err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err, err)
				}
			} else {
				if out == nil {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
				}
			}
		})
	}
}

func (s *JwtManagerTestSuite) TestJwtManager_GetMeFromJWT() {
	type expectation struct {
		out *CurrentUser
		err error
	}

	tests := map[string]struct {
		ctx         context.Context
		tokenAccess string
		expected    expectation
	}{
		// "Success": {
		// 	expected: expectation{
		// 		out: &UserClaims{
		// 			StandardClaims: jwt.StandardClaims{
		// 				ExpiresAt: 1695030127,
		// 			},
		// 		},
		// 		err: nil,
		// 	},
		// },
		"TokenEmpty": {
			ctx:         s.ctx,
			tokenAccess: "",
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unauthenticated desc = session is empty"),
			},
		},
		"InvalidToken": {
			ctx:         s.ctx,
			tokenAccess: "neUQdspTBcXAXntqHlay4gqkwy4hHKq4f8sV5S3b6Ko",
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unauthenticated desc = Session expired"),
			},
		},
		"Failed": {
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = PermissionDenied desc = Authority Denied"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			manager := &JWTManager{
				tokenDuration: 1 * time.Hour,
				secretKey:     "unit-test-123456",
				logger:        s.logger,
			}
			if tt.tokenAccess == "" && scenario != "TokenEmpty" {
				tt.tokenAccess, _ = manager.Generate()
				tt.ctx = metadata.NewIncomingContext(s.ctx, metadata.New(map[string]string{"authorization": "bearer " + tt.tokenAccess}))
			}

			out, err := manager.GetMeFromJWT(tt.ctx, tt.tokenAccess)

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
