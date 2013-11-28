package probe

/**
 * 序列化与反序列化的包
 * 序列化后的格式为: int32的包长度信息+int32的消息类型标识+包内容
 */

const SEND_BRO_GLOBAL string = "BRO_GLOBAL"

const SEND_BRO_GROUP string = "BRO_GROUP"

const SEND_TO_SERVER string = "TO_SERVER"

type VictoriestMsg struct {
	MsgType    int32
	MsgContext interface{}
}

// 序列化接口
type ISerializable interface {
	/**
	 * param  : src    - 需要序列化的参数
	 *          msgType- 需要序列化的对象标识
	 * return : []byte - 序列化后的byte数组
	 *          error  - 错误信息, 如果成功则为nil
	 */
	Serialize(src *VictoriestMsg) ([]byte, error)

	/**
	 * param  : src            - 序列化过的对象
	 *          dst            - 反序列化后的对象
	 * return : error          - 错误信息, 如果成功则为nil
	 *          msgType        - 反序列化后的对象标识
	 */
	Deserialize(src []byte, dst *VictoriestMsg) (int32, error)
}
