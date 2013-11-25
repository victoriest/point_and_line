package probe

/**
 * 序列化与反序列化的包
 * 序列化后的格式为: int32的包长度信息+int32的消息类型标识+包内容
 */
import (
	"bufio"
	"bytes"
	log "code.google.com/p/log4go"
	"encoding/binary"
	// "encoding/gob"
	"encoding/json"
)

// 序列化接口
type Serializable interface {
	/**
	 * param  : src    - 需要序列化的参数
	 *          msgType- 需要序列化的对象标识
	 * return : []byte - 序列化后的byte数组
	 *          error  - 错误信息, 如果成功则为nil
	 */
	Serialize(src interface{}, msgType int32) ([]byte, error)

	/**
	* param  : src         - 序列化过的对象
	*          interface{} - 反序列化后的对象
	* return : error       - 错误信息, 如果成功则为nil
				msgType     - 反序列化后的对象标识
	*/
	Deserialize(src []byte, dst interface{}) (int32, error)
}

type Codecable interface{}

// Josn的序列化实现
type JsonProbe struct{}

func (self JsonProbe) Serialize(src interface{}, msgType int32) ([]byte, error) {
	var v []byte
	var err error

	v, err = json.Marshal(src)
	if err != nil {
		log.Error("when Encoding:", err.Error())
		return nil, err
	}

	var length int32 = int32(len(v))

	pkg := new(bytes.Buffer)

	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		log.Error("when Encoding length:", err.Error())
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, msgType)
	if err != nil {
		log.Error("when Encoding msgType:", err.Error())
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		log.Error("when Encoding msgContext:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

func (self JsonProbe) Deserialize(src []byte, dst interface{}) (int32, error) {
	msg := src[8:]
	var msgType int32
	data := bytes.NewBuffer(src[4:8])
	err := binary.Read(data, binary.LittleEndian, &msgType)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return -1, err
	}

	err = json.Unmarshal(msg, &dst)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return -1, err
	}

	return msgType, nil
}

func (self JsonProbe) DeserializeByReader(reader *bufio.Reader) (interface{}, int32, error) {
	buff, _ := reader.Peek(4)
	data := bytes.NewBuffer(buff)
	var length int32
	err := binary.Read(data, binary.LittleEndian, &length)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return nil, -1, err
	}

	if int32(reader.Buffered()) < length+8 {
		log.Error("int32(reader.Buffered()) < length + 8")
		return nil, -1, err
	}

	pack := make([]byte, int(8+length))
	_, err = reader.Read(pack)
	if err != nil {
		log.Error("when Deserialize:", err.Error())
		return nil, -1, err
	}
	var dst interface{}
	msgType, _ := self.Deserialize(pack, &dst)

	return dst, msgType, nil
}

// // Gob的序列化实现
// type GobProbe struct{}

// // gob的序列化方法实现
// func (self GobProbe) Serialize(src interface{}) (v []byte, err error) {
// 	// 序列化
// 	buf := new(bytes.Buffer)
// 	enc := gob.NewEncoder(buf)
// 	err = enc.Encode(src)
// 	if err != nil {
// 		log.Error("when GobProbe.Encoding:", err.Error())
// 		return nil, err
// 	}
// 	v = buf.Bytes()
// 	// 序列化后的byte长度
// 	var length int32 = int32(len(v))

// 	// 将长度信息写入byte数组
// 	pkg := new(bytes.Buffer)
// 	err = binary.Write(pkg, binary.LittleEndian, length)
// 	if err != nil {
// 		log.Error("when GobProbe.Encoding:", err.Error())
// 		return nil, err
// 	}
// 	err = binary.Write(pkg, binary.LittleEndian, v)
// 	if err != nil {
// 		log.Error("when GobProbe.Encoding:", err.Error())
// 		return nil, err
// 	}

// 	return pkg.Bytes(), nil
// }

// // gob的反序列化方法实现
// func (self GobProbe) Deserialize(src []byte, dst interface{}) error {
// 	msg := src[4:]
// 	buf := bytes.NewBuffer(msg)
// 	dec := gob.NewDecoder(buf)
// 	err := dec.Decode(dst)

// 	if err != nil {
// 		log.Error("when GobProbe.Deserialize:", err.Error())
// 		return err
// 	}

// 	return nil
// }
