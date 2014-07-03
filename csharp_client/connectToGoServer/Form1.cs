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
            EnableGameUI(true);
        }

        private void DisconnectedCallBack()
        {
            lbInfo.Items.Add("disconnected");
            EnableGameUI(false);
        }

        private void EnableGameUI(bool enable) {
            txtNick.Enabled = !enable;
            txtIp.Enabled = !enable;
            txtPort.Enabled = !enable;
            btnConnect.Enabled = !enable;

            btnDisconnect.Enabled = enable;
            txtMsg.Enabled = enable;
            btnSend.Enabled = enable;
            btnStart.Enabled = enable;
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
                case (int)MessageType.MSG_TYPE_CHAT_MESSGAE:
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

        private void btnConnect_Click(object sender, EventArgs e)
        {
            connector.InitSocket(txtIp.Text, Int32.Parse(txtPort.Text));
        }

        private void btnDisconnect_Click(object sender, EventArgs e)
        {
            connector.CloseConnect();
        }

        private void btnSend_Click(object sender, EventArgs e)
        {
            MobileSuiteModel msm = new MobileSuiteModel();
            msm.type = (int)MessageType.MSG_TYPE_CHAT_MESSGAE;

            ChatMsg chat = new ChatMsg();
            chat.chatContext = txtMsg.Text;

            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize<ChatMsg>(ms, chat);
                msm.message = ms.ToArray();
                ms.Close();
            }
            byte[] bytes;
            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize(ms, msm);
                bytes = ms.ToArray();
                ms.Close();
            }

            int length = bytes.Length;
            byte[] data = new byte[length + 4];
            byte[] lengthBytes = BitConverter.GetBytes(length);
            Array.Copy(lengthBytes, data, 4);
            Array.Copy(bytes, 0, data, 4, length);

            connector.SendMessage(data);
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
                MobileSuiteModel msm = new MobileSuiteModel();
                msm.type = (int)MessageType.MSG_TYPE_LINE_A_POINT_REQ;

                LineAPointDTO lpDto = new LineAPointDTO();
                lpDto.row = Int32.Parse(arr[0]);
                lpDto.col = Int32.Parse(arr[1]);
                lpDto.playerIndex = playerIndex;

                using (MemoryStream ms = new MemoryStream())
                {
                    ProtoBuf.Serializer.Serialize<LineAPointDTO>(ms, lpDto);
                    msm.message = ms.ToArray();
                    ms.Close();
                }
                byte[] bytes;
                using (MemoryStream ms = new MemoryStream())
                {
                    ProtoBuf.Serializer.Serialize(ms, msm);
                    bytes = ms.ToArray();
                    ms.Close();
                }

                int length = bytes.Length;
                byte[] data = new byte[length + 4];
                byte[] lengthBytes = BitConverter.GetBytes(length);
                Array.Copy(lengthBytes, data, 4);
                Array.Copy(bytes, 0, data, 4, length);

                connector.SendMessage(data);
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

            MobileSuiteModel msm = new MobileSuiteModel();
            msm.type = (int)MessageType.MSG_TYPE_SEARCH_A_GAME_REQ;

            ChatMsg chat = new ChatMsg();
            chat.chatContext = txtNick.Text;

            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize<ChatMsg>(ms, chat);
                msm.message = ms.ToArray();
                ms.Close();
            }
            byte[] bytes;
            using (MemoryStream ms = new MemoryStream())
            {
                ProtoBuf.Serializer.Serialize(ms, msm);
                bytes = ms.ToArray();
                ms.Close();
            }

            int length = bytes.Length;
            byte[] data = new byte[length + 4];
            byte[] lengthBytes = BitConverter.GetBytes(length);
            Array.Copy(lengthBytes, data, 4);
            Array.Copy(bytes, 0, data, 4, length);

            connector.SendMessage(data);
        }
    }
}
