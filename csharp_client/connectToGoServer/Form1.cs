using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace connectToGoServer
{
    public partial class Form1 : Form
    {
        private ServerConnector connector;

        public Form1()
        {
            InitializeComponent();
            connector = ServerConnector.GetInstance();
            connector.OnConnectedEvent += ConnectedCallBack;
            connector.OnRecivedMessageEvent += RecivedMessage;
            connector.OnDisconnectedEvent += DisconnectedCallBack;
            Control.CheckForIllegalCrossThreadCalls = false;
            
        }

        private void ConnectedCallBack(IAsyncResult ar){
            //MessageBox.Show("connected");
            lbInfo.Items.Add("connected");
        }

        private void RecivedMessage(byte[] data)
        {
            //MessageBox.Show(BitConverter.ToString(data));
            lbInfo.Items.Add(BitConverter.ToString(data));
        }

        private void DisconnectedCallBack() {
            //MessageBox.Show("disconnected");
            lbInfo.Items.Add("disconnected");
        }

        private void btnConnect_Click(object sender, EventArgs e)
        {
            connector.InitSocket(txtIp.Text, Int32.Parse(txtPort.Text));
        }

        private void btnDisconnect_Click(object sender, EventArgs e)
        {
            connector.CloseConnect();
        }
    }
}
