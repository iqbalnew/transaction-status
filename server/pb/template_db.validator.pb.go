// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: template_db.proto

package pb

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Templates) Validate() error {
	if this.CreatedDt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedDt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedDt", err)
		}
	}
	if this.UpdatedDt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UpdatedDt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UpdatedDt", err)
		}
	}
	return nil
}
