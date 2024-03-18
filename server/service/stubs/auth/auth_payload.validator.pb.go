// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: auth_payload.proto

package pb

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/descriptorpb"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *LoginRequest) Validate() error {
	return nil
}
func (this *AccessToken) Validate() error {
	return nil
}
func (this *InitToken) Validate() error {
	return nil
}
func (this *RefreshToken) Validate() error {
	return nil
}
func (this *RefreshRequest) Validate() error {
	return nil
}
func (this *LoginResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *InitTokenLoginReqEncrypted) Validate() error {
	return nil
}
func (this *InitLoginResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *SSOLoginResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *RefreshResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *HealthCheckResponse) Validate() error {
	return nil
}
func (this *Empty) Validate() error {
	return nil
}
func (this *ErrorBodyResponse) Validate() error {
	return nil
}
func (this *ProductAuthority) Validate() error {
	return nil
}
func (this *LogoutRequest) Validate() error {
	return nil
}
func (this *VerifyTokenReq) Validate() error {
	return nil
}
func (this *VerifySessionReq) Validate() error {
	return nil
}
func (this *VerifyTokenRes) Validate() error {
	for _, item := range this.ProductRoles {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ProductRoles", err)
			}
		}
	}
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	return nil
}
func (this *SetMeRes) Validate() error {
	return nil
}
func (this *FilteredVerifyTokenRes) Validate() error {
	for _, item := range this.ProductRoles {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ProductRoles", err)
			}
		}
	}
	return nil
}
func (this *LoginV2Request) Validate() error {
	return nil
}
func (this *LoginV2Response) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *LoginV2Response_AccessToken) Validate() error {
	return nil
}
func (this *ForgotPasswordRequest) Validate() error {
	return nil
}
func (this *ForgotPasswordResponse) Validate() error {
	return nil
}
func (this *ChangePasswordRequest) Validate() error {
	return nil
}
func (this *ChangePasswordResponse) Validate() error {
	return nil
}
func (this *VerifyUserQuestionRequest) Validate() error {
	return nil
}
func (this *VerifyUserQuestionResponse) Validate() error {
	return nil
}
func (this *VerifyChangePasswordTokenRequest) Validate() error {
	return nil
}
func (this *VerifyChangePasswordTokenResponse) Validate() error {
	return nil
}
func (this *RegisterUserIntoAuthReq) Validate() error {
	return nil
}
func (this *RegisterUserIntoAuthRes) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *Credential) Validate() error {
	return nil
}
func (this *DeleteAccessTokenRequest) Validate() error {
	return nil
}
func (this *DeleteAccessTokenResponse) Validate() error {
	return nil
}
func (this *AuthenticationData) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *AuthenticationData_Response) Validate() error {
	return nil
}
func (this *RequestUserId) Validate() error {
	return nil
}
func (this *QLolaUserValidationResponse) Validate() error {
	return nil
}
func (this *QlolaUserValidationBrigateResponse) Validate() error {
	if this.ResponseData != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ResponseData); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ResponseData", err)
		}
	}
	return nil
}
