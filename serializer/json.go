package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON is a serializer that serializes protobuf messages to JSON.
func ProtobufToJSON(message proto.Message) (string, error) {
	// Use jsonpb will cause the error:
	// protoreflect.ProtoMessage does not implement protoiface.MessageV1
	// (missing ProtoMessage method)
	// Inconsistence between jsonpb and proto.Message
	// proto.Message is from google.golang.org/protobuf/proto
	// jsonpb is from github.com/golang/protobuf 
	//
	// marshaler := jsonpb.Marshaler {
	// 	EnumsAsInts: false,
	// 	EmitDefaults: true,
	// 	Indent: " ",
	// 	OrigName: true,
	// }
	marshaler := protojson.MarshalOptions{
		UseEnumNumbers: false,
		EmitUnpopulated: true,
		Indent: " ",
		UseProtoNames: true,
	}

	json, err := marshaler.Marshal(message)
	return string(json), err
}