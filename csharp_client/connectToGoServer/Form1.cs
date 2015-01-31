using protocol;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Runtime.Serialization;
using System.Runtime.Serialization.Formatters.Binary;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace connectToGoServer
{
    public partial class Form1 : Form
    {
        private ServerConnector connector;

        private LineAndPointGame game;

        private bool isEnd;

        private int playerIndex = 0;

        private string opptNick = "";

        //private string serverIp = "115.159.40.89";
        private string serverIp = "127.0.0.1";

        private int serverPort = 8990;

        private LoginResultDTO playerInfo;

        private LoginResultDTO opptInfo;

        public delegate void startNewGameDelegate();
        public startNewGameDelegate start;

        public Form1()
        {
            InitializeComponent();
            connector = ServerConnector.GetInstance();
            connector.OnConnectedEvent += ConnectedCallBack;
            connector.OnRecivedMessageEvent += RecivedMessage;
            connector.OnDisconnectedEvent += DisconnectedCallBack;
            Control.CheckForIllegalCrossThreadCalls = false;
            EnableGameUI(false);
            playerIndex = 0;
            opptNick = "";
        }

        private void ConnectedCallBack(IAsyncResult ar)
        {
            lbInfo.Items.Add("connected");
        }

        private void DisconnectedCallBack()
        {
            lbInfo.Items.Add("disconnected");
            EnableGameUI(false);
            this.Close();
        }

        private void EnableGameUI(bool enable) {
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
                    playerInfo = loginResult;
                    EnableGameUI(true);
                    break;
                case (int)MessageType.MSG_TYPE_CREATE_USER_RES:
                    stream = new MemoryStream(msm.message);
                    CreateResultDTO createResult = ProtoBuf.Serializer.Deserialize<CreateResultDTO>(stream);
                    lbInfo.Items.Add(createResult.userId);
                    txtIn.Text = createResult.userId.ToString();
                    break;
                case (int)MessageType.MSG_TYPE_CHAT_MESSAGE_RES:
                    stream = new MemoryStream(msm.message);
                    ChatMsg chat = ProtoBuf.Serializer.Deserialize<ChatMsg>(stream);
                    lbInfo.Items.Add(chat.chatContext);
                    break;
                case (int)MessageType.MSG_TYPE_START_RES:
                    stream = new MemoryStream(msm.message);
                    GameStartDTO gsDto = ProtoBuf.Serializer.Deserialize<GameStartDTO>(stream);
                    playerIndex = gsDto.playerIndex;
                    opptNick = gsDto.opptName;
                    //StartGame();
                    start = new startNewGameDelegate(StartGame);
                    this.Invoke(start);
                    break;
                case (int)MessageType.MSG_TYPE_LINE_A_POINT_RES:
                    stream = new MemoryStream(msm.message);
                    LineAPointDTO lpDto = ProtoBuf.Serializer.Deserialize<LineAPointDTO>(stream);
                    int result = game.Line(lpDto.row, lpDto.col, lpDto.playerIndex);

                    String btnName = String.Format("{0}_{1}", lpDto.row, lpDto.col);
                    Button currentButton = null;
                    foreach (Control c in this.Controls)
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
                    if (game.gameState == 2)
                    {
                        MessageBox.Show("游戏结束");
                        isEnd = true;
                    }
                    break;
            }
        }

        private void btnDisconnect_Click(object sender, EventArgs e)
        {
            connector.CloseConnect();
        }

        private void btnSend_Click(object sender, EventArgs e)
        {
            ChatMsg chat = new ChatMsg();
            chat.chatContext = txtMsg.Text;
            connector.SendMessage<ChatMsg>((int)MessageType.MSG_TYPE_CHAT_MESSGAE_REQ, chat);
        }

        private void btnClick(object sender, System.EventArgs e)
        {
            if (isEnd)
                return;
            Button currentButton = (Button)sender;
            String str = currentButton.Name;
            String[] arr = str.Split(new char[]{'_'});

            int playerId = game.whosTurn;
            int result = game.Line(Int32.Parse(arr[0]), Int32.Parse(arr[1]), playerIndex);
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
                lpDto.playerIndex = playerIndex;
                connector.SendMessage<LineAPointDTO>((int)MessageType.MSG_TYPE_LINE_A_POINT_REQ, lpDto);
            }
            
            updateState();
            if (game.gameState == 2)
            {
                MessageBox.Show("游戏结束");
                isEnd = true;
            }
        }

        private void StartGame()
        {
            game = new LineAndPointGame(2);
            updateState();

            int startX = 490;
            int startY = 66;
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
            Object[] steps = game.gameSteps;
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
                        bt.Click += new System.EventHandler(btnClick);
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
                        bt.Click += new System.EventHandler(btnClick);
                        this.Controls.Add(bt);
                    }
                }
            }
        }

        private void updateState()
        {
            if (game.whosTurn == playerIndex)
            {
                labTurn.Text = "轮到你";
            }
            else
            {
                labTurn.Text = "轮到" + opptNick;
            }
            labPlayer1Socre.Text = game.player1Sorce.ToString();
            labPlayer2Socre.Text = game.player2Sorce.ToString();
        }

        private void btnStart_Click(object sender, EventArgs e)
        {
            playerIndex = 0;
            opptNick = "";

            ChatMsg chat = new ChatMsg();
            chat.chatContext = txtIn.Text;
            connector.SendMessage<ChatMsg>((int)MessageType.MSG_TYPE_SEARCH_A_GAME_REQ, chat);
        }

        private void btnCreateUser_Click(object sender, EventArgs e)
        {
            if (radioSign.Checked)
            {
                CreateUserDTO dto = new CreateUserDTO();
                dto.uName = txtIn.Text;
                connector.SendMessage<CreateUserDTO>((int)MessageType.MSG_TYPE_CREATE_USER_REQ, dto);
            }
            else if (radioLogin.Checked)
            {
                LoginDTO dto = new LoginDTO();
                dto.userId = long.Parse(txtIn.Text);
                connector.SendMessage<LoginDTO>((int)MessageType.MSG_TYPE_LOGIN_REQ, dto);
            }
        }

        private void Form1_Shown(object sender, EventArgs e)
        {
            EnableGameUI(false);
            connector.InitSocket(serverIp, serverPort);
            lbInfo.Items.Add("connecting...");
            radioCheck();
        }

        private void radioLogin_CheckedChanged(object sender, EventArgs e)
        {
            radioCheck();
        }

        private void radioSign_CheckedChanged(object sender, EventArgs e)
        {
            radioCheck();
        }

        private void radioCheck() 
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
