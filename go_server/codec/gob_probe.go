package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

	log "github.com/alecthomas/log4go"
)

// GobProbe Gob的序列化实现
type GobProbe struct{}

// Serialize gob的序列化方法实现
func (gobProbe *GobProbe) Serialize(src interface{}) ([]byte, error) {
	// 序列化
	// vMsg := &VictoriestMsg{MsgType: msgType, MsgContext: src}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(src)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}
	v := buf.Bytes()
	// 序列化后的byte长度
	var length = int32(len(v))

	// 将长度信息写入byte数组
	pkg := new(bytes.Buffer)
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}

	// 写入序列化后的对象
	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

// Deserialize gob的反序列化方法实现
func (gobProbe *GobProbe) Deserialize(src []byte, dst interface{}) (int32, error) {
	var dstObj = dst.(*VictoriestMsg)
	// msg 序列化后的对象
	msg := src[4:]
	buf := bytes.NewBuffer(msg)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&dstObj)
	if err != nil {
		log.Error("when GobProbe.Deserialize:", err.Error())
		println("when GobProbe.Deserialize:", err.Error())
		return -1, err
	}

	return dstObj.MsgType, nil
}
