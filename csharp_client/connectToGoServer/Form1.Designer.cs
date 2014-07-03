namespace connectToGoServer
{
    partial class Form1
    {
        /// <summary>
        /// 必需的设计器变量。
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// 清理所有正在使用的资源。
        /// </summary>
        /// <param name="disposing">如果应释放托管资源，为 true；否则为 false。</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows 窗体设计器生成的代码

        /// <summary>
        /// 设计器支持所需的方法 - 不要
        /// 使用代码编辑器修改此方法的内容。
        /// </summary>
        private void InitializeComponent()
        {
            this.btnConnect = new System.Windows.Forms.Button();
            this.txtIp = new System.Windows.Forms.TextBox();
            this.txtPort = new System.Windows.Forms.TextBox();
            this.btnDisconnect = new System.Windows.Forms.Button();
            this.lbInfo = new System.Windows.Forms.ListBox();
            this.txtMsg = new System.Windows.Forms.TextBox();
            this.btnSend = new System.Windows.Forms.Button();
            this.btnStart = new System.Windows.Forms.Button();
            this.label3 = new System.Windows.Forms.Label();
            this.label4 = new System.Windows.Forms.Label();
            this.labPlayer1Socre = new System.Windows.Forms.Label();
            this.labPlayer2Socre = new System.Windows.Forms.Label();
            this.labTurn = new System.Windows.Forms.Label();
            this.label1 = new System.Windows.Forms.Label();
            this.txtNick = new System.Windows.Forms.TextBox();
            this.SuspendLayout();
            // 
            // btnConnect
            // 
            this.btnConnect.Location = new System.Drawing.Point(308, 12);
            this.btnConnect.Name = "btnConnect";
            this.btnConnect.Size = new System.Drawing.Size(75, 23);
            this.btnConnect.TabIndex = 0;
            this.btnConnect.Text = "connect";
            this.btnConnect.UseVisualStyleBackColor = true;
            this.btnConnect.Click += new System.EventHandler(this.btnConnect_Click);
            // 
            // txtIp
            // 
            this.txtIp.Location = new System.Drawing.Point(171, 14);
            this.txtIp.Name = "txtIp";
            this.txtIp.Size = new System.Drawing.Size(92, 21);
            this.txtIp.TabIndex = 1;
            this.txtIp.Text = "192.168.0.31";
            // 
            // txtPort
            // 
            this.txtPort.Location = new System.Drawing.Point(269, 14);
            this.txtPort.Name = "txtPort";
            this.txtPort.Size = new System.Drawing.Size(33, 21);
            this.txtPort.TabIndex = 2;
            this.txtPort.Text = "8990";
            // 
            // btnDisconnect
            // 
            this.btnDisconnect.Location = new System.Drawing.Point(389, 12);
            this.btnDisconnect.Name = "btnDisconnect";
            this.btnDisconnect.Size = new System.Drawing.Size(75, 23);
            this.btnDisconnect.TabIndex = 5;
            this.btnDisconnect.Text = "disconnect";
            this.btnDisconnect.UseVisualStyleBackColor = true;
            this.btnDisconnect.Click += new System.EventHandler(this.btnDisconnect_Click);
            // 
            // lbInfo
            // 
            this.lbInfo.FormattingEnabled = true;
            this.lbInfo.ItemHeight = 12;
            this.lbInfo.Location = new System.Drawing.Point(14, 44);
            this.lbInfo.Name = "lbInfo";
            this.lbInfo.Size = new System.Drawing.Size(450, 352);
            this.lbInfo.TabIndex = 6;
            // 
            // txtMsg
            // 
            this.txtMsg.Location = new System.Drawing.Point(12, 404);
            this.txtMsg.Name = "txtMsg";
            this.txtMsg.Size = new System.Drawing.Size(369, 21);
            this.txtMsg.TabIndex = 7;
            this.txtMsg.Text = "hello wrold";
            // 
            // btnSend
            // 
            this.btnSend.Location = new System.Drawing.Point(389, 402);
            this.btnSend.Name = "btnSend";
            this.btnSend.Size = new System.Drawing.Size(75, 23);
            this.btnSend.TabIndex = 8;
            this.btnSend.Text = "send";
            this.btnSend.UseVisualStyleBackColor = true;
            this.btnSend.Click += new System.EventHandler(this.btnSend_Click);
            // 
            // btnStart
            // 
            this.btnStart.Location = new System.Drawing.Point(490, 12);
            this.btnStart.Name = "btnStart";
            this.btnStart.Size = new System.Drawing.Size(63, 48);
            this.btnStart.TabIndex = 9;
            this.btnStart.Text = "start";
            this.btnStart.UseVisualStyleBackColor = true;
            this.btnStart.Click += new System.EventHandler(this.btnStart_Click);
            // 
            // label3
            // 
            this.label3.AutoSize = true;
            this.label3.Location = new System.Drawing.Point(559, 17);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(65, 12);
            this.label3.TabIndex = 10;
            this.label3.Text = "玩家1得分:";
            // 
            // label4
            // 
            this.label4.AutoSize = true;
            this.label4.Location = new System.Drawing.Point(559, 48);
            this.label4.Name = "label4";
            this.label4.Size = new System.Drawing.Size(65, 12);
            this.label4.TabIndex = 11;
            this.label4.Text = "玩家2得分:";
            // 
            // labPlayer1Socre
            // 
            this.labPlayer1Socre.AutoSize = true;
            this.labPlayer1Socre.Location = new System.Drawing.Point(630, 17);
            this.labPlayer1Socre.Name = "labPlayer1Socre";
            this.labPlayer1Socre.Size = new System.Drawing.Size(11, 12);
            this.labPlayer1Socre.TabIndex = 12;
            this.labPlayer1Socre.Text = "0";
            // 
            // labPlayer2Socre
            // 
            this.labPlayer2Socre.AutoSize = true;
            this.labPlayer2Socre.Location = new System.Drawing.Point(630, 48);
            this.labPlayer2Socre.Name = "labPlayer2Socre";
            this.labPlayer2Socre.Size = new System.Drawing.Size(11, 12);
            this.labPlayer2Socre.TabIndex = 13;
            this.labPlayer2Socre.Text = "0";
            // 
            // labTurn
            // 
            this.labTurn.AutoSize = true;
            this.labTurn.Location = new System.Drawing.Point(683, 30);
            this.labTurn.Name = "labTurn";
            this.labTurn.Size = new System.Drawing.Size(29, 12);
            this.labTurn.TabIndex = 14;
            this.labTurn.Text = "轮到";
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(12, 17);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(59, 12);
            this.label1.TabIndex = 15;
            this.label1.Text = "你的昵称:";
            // 
            // txtNick
            // 
            this.txtNick.Location = new System.Drawing.Point(67, 14);
            this.txtNick.Name = "txtNick";
            this.txtNick.Size = new System.Drawing.Size(98, 21);
            this.txtNick.TabIndex = 16;
            // 
            // Form1
            // 
            this.AcceptButton = this.btnSend;
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 12F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(847, 436);
            this.Controls.Add(this.txtNick);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.labTurn);
            this.Controls.Add(this.labPlayer2Socre);
            this.Controls.Add(this.labPlayer1Socre);
            this.Controls.Add(this.label4);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.btnStart);
            this.Controls.Add(this.btnSend);
            this.Controls.Add(this.txtMsg);
            this.Controls.Add(this.lbInfo);
            this.Controls.Add(this.btnDisconnect);
            this.Controls.Add(this.txtPort);
            this.Controls.Add(this.txtIp);
            this.Controls.Add(this.btnConnect);
            this.Name = "Form1";
            this.Text = "Victoriest";
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Button btnConnect;
        private System.Windows.Forms.TextBox txtIp;
        private System.Windows.Forms.TextBox txtPort;
        private System.Windows.Forms.Button btnDisconnect;
        private System.Windows.Forms.ListBox lbInfo;
        private System.Windows.Forms.TextBox txtMsg;
        private System.Windows.Forms.Button btnSend;
        private System.Windows.Forms.Button btnStart;
        private System.Windows.Forms.Label label3;
        private System.Windows.Forms.Label label4;
        private System.Windows.Forms.Label labPlayer1Socre;
        private System.Windows.Forms.Label labPlayer2Socre;
        private System.Windows.Forms.Label labTurn;
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.TextBox txtNick;
    }
}

