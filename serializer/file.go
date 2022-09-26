package serializer

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
)

// WriteProtobufToBinaryFile is a serializer that serializes protobuf messages to a file.
func WriteProtobufToBinaryFile(message proto.Message, file string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to data: %w", err)
	}
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write data to file: %w", err)
	}

	return nil
}

// ReadProtobufFromBinaryFile is a deserializer that deserializes protobuf messages from a file.
func ReadProtobufFromBinaryFile(file string, message proto.Message) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read data from file: %w", err)
	}
	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data to proto message: %w", err)
	}

	return nil
}

// WriteProtobufToJSONFile is a serializer that serializes protobuf messages to a JSON file.
func WriteProtobufToJSONFile(message proto.Message, file string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %w", err)
	}
	err = ioutil.WriteFile(file, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON to file: %w", err)
	}

	return nil
}