package tcp

import (
	"fmt"
	"os"
)

func Pack(message string) ([]byte, error) {
	length := len(message)

	head := Int322byte(int32(length))
	return append(head, []byte(message)...), nil
}

func UnPack(byteSlice []byte) (string, error) {
	head := byteSlice[:4]
	length := Byte2int32(head)

	content := byteSlice[4 : length+4]

	return string(content), nil
}

func Byte2int32(byteSlice []byte) int32 {
	if len(byteSlice) != 4 {
		fmt.Println("包头长度不对")
		os.Exit(1)
	}
	var int32Num int32
	for i := 0; i < 4; i++ {
		int32Num += int32(byteSlice[i]) << (8 * (3 - i))
	}
	return int32Num
}

func Int322byte(num int32) []byte {
	var ret [4]byte
	ret[3] = uint8(num >> 0)
	ret[2] = uint8(num >> 8)
	ret[1] = uint8(num >> 16)
	ret[0] = uint8(num >> 24)
	//fmt.Println(ret)
	return ret[:]
}
