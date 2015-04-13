using System.Globalization;
using protocol;
using System;
using System.Drawing;
using System.IO;
using System.Windows.Forms;

namespace connectToGoServer
{
    public partial class Form1 : Form
    {
        private readonly ServerConnector _connector;

        private LineAndPointGame _game;

        private bool _isEnd;

        private int _playerIndex;

        private string _opptNick = "";

        //private string serverIp = "115.159.40.89";
        private const string ServerIp = "127.0.0.1";

        private const int ServerPort = 8990;

        private LoginResultDTO _playerInfo;

        private LoginResultDTO _opptInfo;

        public delegate void StartNewGameDelegate();
        public StartNewGameDelegate Start;

        public Form1(LoginResultDTO opptInfo)
        {
            this._opptInfo = opptInfo;
            InitializeComponent();
            _connector = ServerConnector.GetInstance();
            _connector.OnConnectedEvent += ConnectedCallBack;
            _connector.OnRecivedMessageEvent += RecivedMessage;
            _connector.OnDisconnectedEvent += DisconnectedCallBack;
            Control.CheckForIllegalCrossThreadCalls = false;
            EnableGameUi(false);
            _playerIndex = 0;
            _opptNick = "";
        }

        private void ConnectedCallBack(IAsyncResult ar)
        {
            lbInfo.Items.Add("connected");
        }

        private void DisconnectedCallBack()
        {
            lbInfo.Items.Add("disconnected");
            EnableGameUi(false);
            this.Close();
        }

        private void EnableGameUi(bool enable) {
            //txtNick.Enabled = !enable;
            txtMsg.Enabled = enable;
            btnSend.Enabled = enable;
            btnStart.Enabled = enable;

            txtIn.Enabled = !enable;
            btnIn.Enabled = !enable;
            radioSign.Enabled = !enable;
            radioLogin.Enabled = !enable;
        }

        private void RecivedMessage(byte[] data)
        {
            var tmp = new byte[data.Length - 4];
            var lengthByte = new byte[4];
            Array.Copy(data, 0, lengthByte, 0, 4);
            Array.Copy(data, 4, tmp, 0, tmp.Length);
            Stream stream = new MemoryStream(tmp);
            MobileSuiteModel msm = ProtoBuf.Serializer.Deserialize<MobileSuiteModel>(stream);
            switch (msm.type) 
            {
                case (int)MessageType.MSG_TYPE_LOGIN_RES:
                    if (msm.message == null)
                    {
                        MessageBox.Show("无此用户");
                        return;
                    }
                    stream = new MemoryStream(msm.message);
                    LoginResultDTO loginResult = ProtoBuf.Serializer.Deserialize<LoginResultDTO>(stream);
                    txtIn.Text = loginResult.uName;
                    _playerInfo = loginResult;
                    EnableGameUi(true);
                    break;
                case (int)MessageType.MSG_TYPE_LOGOUT_RES:
                case (int)MessageType.MSG_TYPE_END_GAME_RES:
                    MessageBox.Show("玩家退出了游戏");
                    break;
                case (int)MessageType.MSG_TYPE_CREATE_USER_RES:
                    stream = new MemoryStream(msm.message);
                    CreateResultDTO createResult = ProtoBuf.Serializer.Deserialize<CreateResultDTO>(stream);
                    lbInfo.Items.Add(createResult.userId);
                    txtIn.Text = createResult.userId.ToString(CultureInfo.InvariantCulture);
                    break;
                case (int)MessageType.MSG_TYPE_CHAT_MESSAGE_RES:
                    stream = new MemoryStream(msm.message);
                    ChatMsg chat = ProtoBuf.Serializer.Deserialize<ChatMsg>(stream);
                    lbInfo.Items.Add(chat.chatContext);
                    break;
                case (int)MessageType.MSG_TYPE_START_RES:
                    stream = new MemoryStream(msm.message);
                    GameStartDTO gsDto = ProtoBuf.Serializer.Deserialize<GameStartDTO>(stream);
                    _playerIndex = gsDto.playerIndex;
                    _opptNick = gsDto.opptName;
                    //StartGame();
                    Start = new StartNewGameDelegate(StartGame);
                    Invoke(Start);
                    break;
                case (int)MessageType.MSG_TYPE_LINE_A_POINT_TO_REQUEST_RES:
                    break;
                case (int)MessageType.MSG_TYPE_LINE_A_POINT_RES:
                    stream = new MemoryStream(msm.message);
                    LineAPointDTO lpDto = ProtoBuf.Serializer.Deserialize<LineAPointDTO>(stream);
                    int result = _game.Line(lpDto.row, lpDto.col, lpDto.playerIndex);

                    String btnName = String.Format("{0}_{1}", lpDto.row, lpDto.col);
                    Button currentButton = null;
                    foreach (Control c in Controls)
                    {
                        if (c.Name == btnName)
                        {
                            currentButton = (Button)c;
                            break;
                        }
                    }

                    if (result == 0 && lpDto.playerIndex == 1)
                    {
                        Color color = Color.FromArgb(255, 0, 0);
                        currentButton.BackColor = color;
                    }
                    else if (result == 0)
                    {
                        Color color = Color.FromArgb(0, 0, 255);
                        currentButton.BackColor = color;
                    }
            
                    updateState();
                    if (_game.GameState == 2)
                    {
                        MessageBox.Show("游戏结束");
                        _isEnd = true;
                    }
                    break;
            }
        }

        private void btnDisconnect_Click(object sender, EventArgs e)
        {
            _connector.CloseConnect();
        }

        private void btnSend_Click(object sender, EventArgs e)
        {
            ChatMsg chat = new ChatMsg();
            chat.userId = _playerInfo.userId;
            chat.uName = _playerInfo.uName;
            chat.chatContext = String.Format("{0}:{1}", _playerInfo.uName, txtMsg.Text);
            _connector.SendMessage<ChatMsg>((int)MessageType.MSG_TYPE_CHAT_MESSGAE_REQ, chat);
        }

        private void BtnClick(object sender, EventArgs e)
        {
            if (_isEnd)
                return;
            Button currentButton = (Button)sender;
            String str = currentButton.Name;
            String[] arr = str.Split(new char[]{'_'});

            int playerId = _game.WhosTurn;
            int result = _game.Line(Int32.Parse(arr[0]), Int32.Parse(arr[1]), _playerIndex);
            if (result == 1)
            {
                MessageBox.Show("还没轮到你走魂淡!");
                return;
            }
            if (result == 2)
            {
                MessageBox.Show("不许走那里!");
                return;
            }
            if (result == 0 && playerId == 1)
            {
                Color color = Color.FromArgb(255, 0, 0);
                currentButton.BackColor = color;
            }
            else if (result == 0)
            {
                Color color = Color.FromArgb(0, 0, 255);
                currentButton.BackColor = color;
            }

            if (result == 0)
            {
                LineAPointDTO lpDto = new LineAPointDTO();
                lpDto.row = Int32.Parse(arr[0]);
                lpDto.col = Int32.Parse(arr[1]);
                lpDto.playerIndex = _playerIndex;
                _connector.SendMessage<LineAPointDTO>((int)MessageType.MSG_TYPE_LINE_A_POINT_REQ, lpDto);
            }
            
            updateState();
            if (_game.GameState == 2)
            {
                MessageBox.Show("游戏结束");
                _isEnd = true;
            }
        }

        private void StartGame()
        {
            _game = new LineAndPointGame(2);
            updateState();

            const int startX = 490;
            const int startY = 66;
            //Size btnSize = new Size(24, 24);
            //for (int i = 0; i < 5; i++) 
            //{
            //    for (int j = 0; j < 5; j++)
            //    {
            //        Button bt = new Button();
            //        bt.Location = new Point(startX + i * 80, startY + j * 80);
            //        bt.Size = btnSize;
            //        bt.Text = "";
            //        this.Controls.Add(bt);
            //    }
            //}
            Object[] steps = _game.GameSteps;
            int rows = steps.Length;
            for (int i = 0; i < rows; i++)
            {
                // 偶数
                if (i == 0 || i % 2 == 0)
                {
                    int seed = i / 2;
                    int[] colObj = (int[])steps[i];
                    int cols = colObj.Length;
                    for (int j = 0; j < cols; j++)
                    {
                        Button bt = new Button();
                        bt.Location = new Point(startX + j * 80 + 24, startY + seed * 80);
                        bt.Size = new Size(80 - 24, 24);
                        bt.Name = String.Format("{0}_{1}", i, j);
                        bt.Click += new EventHandler(BtnClick);
                        this.Controls.Add(bt);
                    }

                }
                // 奇数
                else
                {
                    int seed = i / 2;
                    int[] colObj = (int[])steps[i];
                    int cols = colObj.Length;
                    for (int j = 0; j < cols; j++)
                    {
                        Button bt = new Button();
                        bt.Location = new Point(startX + j * 80, startY + seed * 80 + 24);
                        bt.Size = new Size(24, 80 - 24);
                        bt.Name = String.Format("{0}_{1}", i, j);
                        bt.Click += new EventHandler(BtnClick);
                        this.Controls.Add(bt);
                    }
                }
            }
        }

        private void updateState()
        {
            if (_game.WhosTurn == _playerIndex)
            {
                labTurn.Text = "轮到你";
            }
            else
            {
                labTurn.Text = "轮到" + _opptNick;
            }
            labPlayer1Socre.Text = _game.Player1Sorce.ToString(CultureInfo.InvariantCulture);
            labPlayer2Socre.Text = _game.Player2Sorce.ToString(CultureInfo.InvariantCulture);
        }

        private void btnStart_Click(object sender, EventArgs e)
        {
            _playerIndex = 0;
            _opptNick = "";

            ChatMsg chat = new ChatMsg();
            chat.chatContext = txtIn.Text;
            _connector.SendMessage<ChatMsg>((int)MessageType.MSG_TYPE_SEARCH_A_GAME_REQ, chat);
        }

        private void btnCreateUser_Click(object sender, EventArgs e)
        {
            if (radioSign.Checked)
            {
                CreateUserDTO dto = new CreateUserDTO();
                dto.uName = txtIn.Text;
                dto.pwd = txtPwd.Text;
                _connector.SendMessage<CreateUserDTO>((int)MessageType.MSG_TYPE_CREATE_USER_REQ, dto);
            }
            else if (radioLogin.Checked)
            {
                LoginDTO dto = new LoginDTO();
                dto.userId = 1;
                dto.uName = txtIn.Text;
                dto.pwd = txtPwd.Text;
                _connector.SendMessage<LoginDTO>((int)MessageType.MSG_TYPE_LOGIN_REQ, dto);
            }
        }

        private void Form1_Shown(object sender, EventArgs e)
        {
            EnableGameUi(false);
            _connector.InitSocket(ServerIp, ServerPort);
            lbInfo.Items.Add("connecting...");
            RadioCheck();
        }

        private void radioLogin_CheckedChanged(object sender, EventArgs e)
        {
            RadioCheck();
        }

        private void radioSign_CheckedChanged(object sender, EventArgs e)
        {
            RadioCheck();
        }

        private void RadioCheck() 
        {
            if (radioLogin.Checked)
            {
                lblIn.Text = "登陆id:";
                btnIn.Text = "登陆";
            }
            else if (radioSign.Checked)
            {
                lblIn.Text = "注册昵称:";
                btnIn.Text = "注册";
            }
        }
    }
}
