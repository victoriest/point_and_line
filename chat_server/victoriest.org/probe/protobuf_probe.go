package probe

import (
	pb "./victoriest.org/protobuf"
	"bufio"
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	log "code.google.com/p/log4go"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type ProtobufProbe struct{}

/**
 * param  : src    - 需要序列化的参数
 *          msgType- 需要序列化的对象标识
 * return : []byte - 序列化后的byte数组
 *          error  - 错误信息, 如果成功则为nil
 */
func (self ProtobufProbe) Serialize(src *VictoriestMsg) ([]byte, error) {
	var v []bytes
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
func (self ProtobufProbe) Deserialize(src []byte, dst *VictoriestMsg) (int32, error) {
	// msg 序列化后的对象
	msg := src[4:]

	err := proto.Unmarshal(msg, dst)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return -1, err
	}

	return dst.MsgType, nil
}
