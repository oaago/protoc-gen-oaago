// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: contract/base/resp.proto

package base

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	github_com_oaago_go_proto_validators "github.com/oaago/go-proto-validators"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Data) Validate() error {
	return nil
}
func (this *BaseResp) Validate() error {
	if this.Data != nil {
		if err := github_com_oaago_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_oaago_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
