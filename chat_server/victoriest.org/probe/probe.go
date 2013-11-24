package probe

/**
 * 序列化与反序列化的包
 * 序列化后的格式为: int32的包长度信息+包内容
 */
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// 序列化接口
type Serializable interface {
	/**
	 * param  : src    - 需要序列化的参数
	 * return : []byte - 序列化后的byte数组
	 *          error  - 错误信息, 如果成功则为nil
	 */
	Serialize(src interface{}) ([]byte, error)

	/**
	 * param  : src         - 序列化过的对象
	 *          interface{} - 反序列化后的对象
	 * return : error       - 错误信息, 如果成功则为nil
	 */
	Deserialize(src []byte, dst interface{}) error

	/**
	 * param  : reader      - 读取序列化数组的reader
	 * return : interface{} - 反序列化后的对象
	 *          error       - 错误信息, 如果成功则为nil
	 */
	//DeserializeByReader(reader *bufio.Reader) (interface{}, error)
}

type Codecable interface{}

// Josn的序列化实现
type JsonProbe struct{}

func (self JsonProbe) Serialize(src interface{}) ([]byte, error) {
	var v []byte
	var err error

	v, err = json.Marshal(src)
	if err != nil {
		fmt.Println("when Encoding:", err.Error())
		return nil, err
	}

	var length int32 = int32(len(v))

	pkg := new(bytes.Buffer)
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		fmt.Println("when Encoding:", err.Error())
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		fmt.Println("when Encoding:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

func (self JsonProbe) Deserialize(src []byte, dst interface{}) error {
	msg := src[4:]
	err := json.Unmarshal(msg, &dst)
	if err != nil {
		fmt.Println("when Deserialize:", err.Error())
		return err
	}
	return nil
}

func (self JsonProbe) DeserializeByReader(reader *bufio.Reader, dst interface{}) error {
	buff, _ := reader.Peek(4)
	data := bytes.NewBuffer(buff)
	var length int32
	err := binary.Read(data, binary.LittleEndian, &length)
	if err != nil {
		fmt.Println("when Deserialize:", err.Error())
		return err
	}

	if int32(reader.Buffered()) < length+4 {
		fmt.Println("int32(reader.Buffered()) < length+4")
		_, err := reader.Peek(int(4 + length))
		if err != nil {
			return err
		}
	}
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		fmt.Println("when Deserialize:", err.Error())
		return err
	}

	return self.Deserialize(pack, dst)
}

// Gob的序列化实现
type GobProbe struct{}

// gob的序列化方法实现
func (self GobProbe) Serialize(src interface{}) (v []byte, err error) {
	// 序列化
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(src)
	if err != nil {
		fmt.Println("when GobProbe.Encoding:", err.Error())
		return nil, err
	}
	v = buf.Bytes()
	// 序列化后的byte长度
	var length int32 = int32(len(v))

	// 将长度信息写入byte数组
	pkg := new(bytes.Buffer)
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		fmt.Println("when GobProbe.Encoding:", err.Error())
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, v)
	if err != nil {
		fmt.Println("when GobProbe.Encoding:", err.Error())
		return nil, err
	}

	return pkg.Bytes(), nil
}

// gob的反序列化方法实现
func (self GobProbe) Deserialize(src []byte, dst interface{}) error {
	msg := src[4:]
	buf := bytes.NewBuffer(msg)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(dst)

	if err != nil {
		fmt.Println("when GobProbe.Deserialize:", err.Error())
		return err
	}

	return nil
}

func (self GobProbe) DeserializeByReader(reader *bufio.Reader) error {
	buff, _ := reader.Peek(4)
	data := bytes.NewBuffer(buff)
	var length int32
	err := binary.Read(data, binary.LittleEndian, &length)
	if err != nil {
		fmt.Println("when Deserialize2:", err.Error())
		return err
	}

	if int32(reader.Buffered()) < length+4 {
		fmt.Println("int32(reader.Buffered()) < length+4")
		_, err := reader.Peek(int(4 + length))
		if err != nil {
			return err
		}
	}
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		fmt.Println("when Deserialize3:", err.Error())
		return err
	}

	var dst interface{}
	return self.Deserialize(pack, dst)
}
