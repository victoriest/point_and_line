//------------------------------------------------------------------------------
// <auto-generated>
//     This code was generated by a tool.
//
//     Changes to this file may cause incorrect behavior and will be lost if
//     the code is regenerated.
// </auto-generated>
//------------------------------------------------------------------------------

// Generated from: MobileSuite.proto
namespace protocol
{
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MobileSuiteModel")]
  public partial class MobileSuiteModel : global::ProtoBuf.IExtensible
  {
    public MobileSuiteModel() {}
    
    private int _type;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"type", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int type
    {
      get { return _type; }
      set { _type = value; }
    }
    private byte[] _message = null;
    [global::ProtoBuf.ProtoMember(2, IsRequired = false, Name=@"message", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public byte[] message
    {
      get { return _message; }
      set { _message = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"ChatMsg")]
  public partial class ChatMsg : global::ProtoBuf.IExtensible
  {
    public ChatMsg() {}
    
    private string _chatContext;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"chatContext", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string chatContext
    {
      get { return _chatContext; }
      set { _chatContext = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"GameStartDTO")]
  public partial class GameStartDTO : global::ProtoBuf.IExtensible
  {
    public GameStartDTO() {}
    
    private string _opptName;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"opptName", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string opptName
    {
      get { return _opptName; }
      set { _opptName = value; }
    }
    private int _playerIndex;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"playerIndex", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int playerIndex
    {
      get { return _playerIndex; }
      set { _playerIndex = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"LineAPointDTO")]
  public partial class LineAPointDTO : global::ProtoBuf.IExtensible
  {
    public LineAPointDTO() {}
    
    private int _row;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"row", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int row
    {
      get { return _row; }
      set { _row = value; }
    }
    private int _col;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"col", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int col
    {
      get { return _col; }
      set { _col = value; }
    }
    private int _playerIndex;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"playerIndex", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int playerIndex
    {
      get { return _playerIndex; }
      set { _playerIndex = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"CreateUserDTO")]
  public partial class CreateUserDTO : global::ProtoBuf.IExtensible
  {
    public CreateUserDTO() {}
    
    private string _name;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"name", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string name
    {
      get { return _name; }
      set { _name = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"CreateResultDTO")]
  public partial class CreateResultDTO : global::ProtoBuf.IExtensible
  {
    public CreateResultDTO() {}
    
    private int _userId;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"userId", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int userId
    {
      get { return _userId; }
      set { _userId = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"LoginDTO")]
  public partial class LoginDTO : global::ProtoBuf.IExtensible
  {
    public LoginDTO() {}
    
    private int _userId;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"userId", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int userId
    {
      get { return _userId; }
      set { _userId = value; }
    }
    private string _name;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"name", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string name
    {
      get { return _name; }
      set { _name = value; }
    }
    private int _Round;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"Round", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int Round
    {
      get { return _Round; }
      set { _Round = value; }
    }
    private int _WinCount;
    [global::ProtoBuf.ProtoMember(4, IsRequired = true, Name=@"WinCount", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int WinCount
    {
      get { return _WinCount; }
      set { _WinCount = value; }
    }
    private int _Rank;
    [global::ProtoBuf.ProtoMember(5, IsRequired = true, Name=@"Rank", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int Rank
    {
      get { return _Rank; }
      set { _Rank = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"LogoutDTO")]
  public partial class LogoutDTO : global::ProtoBuf.IExtensible
  {
    public LogoutDTO() {}
    
    private int _userId;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"userId", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int userId
    {
      get { return _userId; }
      set { _userId = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
    [global::ProtoBuf.ProtoContract(Name=@"MessageType")]
    public enum MessageType
    {
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_CHAT_MESSGAE", Value=1)]
      MSG_TYPE_CHAT_MESSGAE = 1,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_SEARCH_A_GAME_REQ", Value=3)]
      MSG_TYPE_SEARCH_A_GAME_REQ = 3,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_SEARCH_A_GAME_RES", Value=4)]
      MSG_TYPE_SEARCH_A_GAME_RES = 4,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_STOP_SEARCH_REQ", Value=5)]
      MSG_TYPE_STOP_SEARCH_REQ = 5,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_STOP_SEARCH_RES", Value=6)]
      MSG_TYPE_STOP_SEARCH_RES = 6,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_START_RES", Value=8)]
      MSG_TYPE_START_RES = 8,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_LINE_A_POINT_REQ", Value=9)]
      MSG_TYPE_LINE_A_POINT_REQ = 9,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_LINE_A_POINT_RES", Value=10)]
      MSG_TYPE_LINE_A_POINT_RES = 10,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_END_GAME_REQ", Value=11)]
      MSG_TYPE_END_GAME_REQ = 11,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_CREATE_USER_REQ", Value=101)]
      MSG_TYPE_CREATE_USER_REQ = 101,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_CREATE_USER_RES", Value=102)]
      MSG_TYPE_CREATE_USER_RES = 102,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_LOGIN_REQ", Value=103)]
      MSG_TYPE_LOGIN_REQ = 103,
            
      [global::ProtoBuf.ProtoEnum(Name=@"MSG_TYPE_LOGIN_RES", Value=104)]
      MSG_TYPE_LOGIN_RES = 104
    }
  
}