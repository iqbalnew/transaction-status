// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: transaction_status_payload.proto

package pb

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Pagination) Validate() error {
	return nil
}
func (this *GeneralBodyResponse) Validate() error {
	return nil
}
func (this *HealthCheckResponse) Validate() error {
	return nil
}
func (this *UserAuthority) Validate() error {
	return nil
}
func (this *GetAllTemplatesRequest) Validate() error {
	if this.Pagination != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Pagination); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Pagination", err)
		}
	}
	return nil
}
func (this *GetAllTemplatesResponse) Validate() error {
	for _, item := range this.Data {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
			}
		}
	}
	if this.Pagination != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Pagination); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Pagination", err)
		}
	}
	return nil
}
func (this *GetTemplateDetailRequest) Validate() error {
	return nil
}
func (this *GetTemplateDetailResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *SaveTemplateRequest) Validate() error {
	if this.Template != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Template); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Template", err)
		}
	}
	return nil
}
func (this *UpdateTemplateRequest) Validate() error {
	if this.Template != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Template); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Template", err)
		}
	}
	return nil
}
func (this *DeleteTemplateRequest) Validate() error {
	return nil
}