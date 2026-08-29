package api

import "github.com/golang/protobuf/proto"

func (x *RegistrationReq) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *RegistrationReq) Unmarshal(marshalled []byte) error {
	return proto.Unmarshal(marshalled, x)
}

func (x *RegistrationResp) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *RegistrationResp) Unmarshal(marshalled []byte) error {
	return proto.Unmarshal(marshalled, x)
}

func (x *NodeUpdateReq) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *NodeUpdateReq) Unmarshal(marshalled []byte) error {
	return proto.Unmarshal(marshalled, x)
}

func (x *NodeUpdateResp) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *NodeUpdateResp) Unmarshal(marshalled []byte) error {
	return proto.Unmarshal(marshalled, x)
}
