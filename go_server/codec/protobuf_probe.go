package codec

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"go_server/log"
	"go_server/protocol"
)

type ProtobufProbe struct{}

/**
 * param  : src    - 需要序列化的参数
 *          msgType- 需要序列化的对象标识
 * return : []byte - 序列化后的byte数组
 *          error  - 错误信息, 如果成功则为nil
 */
func (self ProtobufProbe) Serialize(src *protocol.MobileSuiteModel) ([]byte, error) {
	var v []byte
	var err error

	v, err = proto.Marshal(src)
	if err != nil {
		log.Error("when encoding:", err.Error())
		return nil, err
	}

	var length int32 = int32(len(v))
	pkg := new(bytes.Buffer)

	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		log.Error("when write length:", err.Error())
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		log.Error("when Serialize:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil

}

/**
 * param  : src            - 序列化过的对象
 *          dst            - 反序列化后的对象
 * return : error          - 错误信息, 如果成功则为nil
 *          msgType        - 反序列化后的对象标识
 */
func (self ProtobufProbe) Deserialize(src []byte, dst *protocol.MobileSuiteModel) (int32, error) {
	// msg 序列化后的对象
	msg := src[4:]

	err := proto.Unmarshal(msg, dst)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return -1, err
	}

	return *dst.Type, nil
}

func (self ProtobufProbe) DeserializeByReader(reader *bufio.Reader) (*protocol.MobileSuiteModel, int32, error) {
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		log.Error("when deserializeByReader:", err.Error())
		return nil, -1, err
	}

	if int32(reader.Buffered()) < length+4 {
		log.Error("int32(reader.Buffered()) < length + 4")
		return nil, -1, err
	}

	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		log.Error("when deserializeByReader:", err.Error())
		return nil, -1, err
	}
	var dst protocol.MobileSuiteModel
	var msgType int32
	msgType, err = self.Deserialize(pack, &dst)
	log.Debug(length, msgType, dst)
	return &dst, msgType, nil
}
