package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func ProtoToJson(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       " ",
		AnyResolver:  nil,
	}
	return marshaler.MarshalToString(message)
}
