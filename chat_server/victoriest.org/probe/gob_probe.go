package probe

import (
	"bytes"
	log "code.google.com/p/log4go"
	"encoding/binary"
	"encoding/gob"
)

// Gob的序列化实现
type GobProbe struct{}

// gob的序列化方法实现
func (self GobProbe) Serialize(src interface{}, msgType int32) ([]byte, error) {
	// 序列化
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(src)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}
	v := buf.Bytes()
	// 序列化后的byte长度
	var length int32 = int32(len(v))

	// 将长度信息写入byte数组
	pkg := new(bytes.Buffer)
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		log.Error("when GobProbe.Encoding:", err.Error())
		return nil, err
	}

	// 写入消息类型标识
	err = binary.Write(pkg, binary.LittleEndian, msgType)
	if err != nil {
		log.Error("when Encoding msgType:", err.Error())
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

// gob的反序列化方法实现
func (self GobProbe) Deserialize(src []byte, dst interface{}) (int32, error) {
	// msgType 对象类型的标示
	var msgType int32
	data := bytes.NewBuffer(src[4:8])
	err := binary.Read(data, binary.LittleEndian, &msgType)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		println("when Deserialize:", err.Error())
		return -1, err
	}

	// msg 序列化后的对象
	msg := src[8:]
	buf := bytes.NewBuffer(msg)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&dst)
	if err != nil {
		log.Error("when GobProbe.Deserialize:", err.Error())
		println("when GobProbe.Deserialize:", err.Error())
		return -1, err
	}

	return msgType, nil
}
