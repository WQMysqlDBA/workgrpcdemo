package serializer

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
	"workgrpc/pb"
	"workgrpc/sample"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()
	binaryFile := "./test.bin"
	jsonFile := "./test.json"
	cpu1 := sample.NewCPU()
	fmt.Println(cpu1)
	err := WriteProtoBufToBinaryFile(cpu1, binaryFile)
	require.NoError(t, err)
	cpu2 := &pb.CPU{}
	err = ReadProtobufFromBinaryFile(binaryFile, cpu2)
	require.NoError(t, err)
	require.True(t, proto.Equal(cpu1, cpu2))
	err = WriteProtobuToJson(cpu1, jsonFile)
	require.NoError(t, err)
}
