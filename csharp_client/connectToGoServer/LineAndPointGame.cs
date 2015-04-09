using System;

namespace connectToGoServer
{
    public class LineAndPointGame
    {
        // 玩家1的得分
        public int Player1Sorce { get; set; }

        // 玩家2的得分
        public int Player2Sorce { get; set; }

        // 棋盘大小 1:3*3, 2:5*5, 3:7*7
        public int GameType { get; set; }

        // 棋盘数据
        public Object[] GameSteps { get; set; }

        // 轮到谁出手
        public int WhosTurn { get; set; }

        // 游戏状态 1:进行中; 2:已结束
        public int GameState { get; set; }

        public LineAndPointGame(int type)
        {
            Player1Sorce = 0;
            Player2Sorce = 0;
            GameType = type;
            WhosTurn = 1;
            GameState = 1;
            if (type == 1)
            {
                InitSteps(3);
            }
            else if(type == 2)
            {
                InitSteps(5);
            }
            else if (type == 3)
            {
                InitSteps(7);
            }
        }

        public int Line(int rowIndex, int colIndex, int playerId)
        {
            if (playerId != WhosTurn)
            {
                // 不轮到你走
                return 1;
            }
            int[] obj = (int[])GameSteps[rowIndex];
            int state = obj[colIndex];
            if (state != 0)
            {
                // 这里不准走魂淡
                return 2;
            }

            // 能走就走上
            obj[colIndex] = WhosTurn;

            // 判断是否得分
            bool isSorce = false;
            if (rowIndex == 0 || rowIndex % 2 == 0)
            {
                // 偶数
                if (rowIndex + 2 < GameSteps.Length)
                {
                    // 下面的方块
                    var bottomBottom = ((int[])GameSteps[rowIndex + 2])[colIndex];
                    var bottomLeft = ((int[])GameSteps[rowIndex + 1])[colIndex];
                    var bottomRight = ((int[])GameSteps[rowIndex + 1])[colIndex + 1];
                    if (bottomBottom != 0 && bottomLeft != 0 && bottomRight != 0 && bottomBottom == bottomLeft && bottomLeft == bottomRight && WhosTurn == bottomRight)
                    {
                        PlusSorce(playerId);
                        isSorce = true;
                    }
                }
                if (rowIndex - 2 >= 0)
                {
                    // 上面的方块
                    var topLeft = ((int[])GameSteps[rowIndex - 1])[colIndex];
                    var topTop = ((int[])GameSteps[rowIndex - 2])[colIndex];
                    var topRight = ((int[])GameSteps[rowIndex - 1])[colIndex + 1];
                    if (topLeft != 0 && topTop != 0 && topRight != 0 && topLeft == topTop && topTop == topRight && WhosTurn == topRight)
                    {
                        PlusSorce(playerId);
                        isSorce = true;
                    }
                }
            }
            else
            {
                // 奇数
                if (colIndex - 1 >= 0)
                {
                    // 左边的方块
                    var leftLeft = ((int[])GameSteps[rowIndex])[colIndex - 1];
                    var leftTop = ((int[])GameSteps[rowIndex - 1])[colIndex - 1];
                    var leftBottom = ((int[])GameSteps[rowIndex + 1])[colIndex - 1];
                    if (leftLeft != 0 && leftTop != 0 && leftBottom != 0 && leftLeft == leftTop && leftTop == leftBottom && WhosTurn == leftBottom)
                    {
                        PlusSorce(playerId);
                        isSorce = true;
                    }
                }
                if (colIndex + 1 < ((int[])GameSteps[rowIndex]).Length)
                {
                    // 右边的方块
                    var rightRight = ((int[])GameSteps[rowIndex])[colIndex + 1];
                    var rightTop = ((int[])GameSteps[rowIndex - 1])[colIndex];
                    var rightBottom = ((int[])GameSteps[rowIndex + 1])[colIndex];
                    if (rightRight != 0 && rightTop != 0 && rightBottom != 0 && rightRight == rightTop && rightTop == rightBottom && WhosTurn == rightBottom)
                    {
                        PlusSorce(playerId);
                        isSorce = true;
                    }
                }

            }

            // 判断该谁走
            if (isSorce)
            {
                WhosTurn = playerId;
            }
            else
            {
                if (playerId == 1)
                {
                    WhosTurn = 2;
                }
                else if (playerId == 2)
                {
                    WhosTurn = 1;
                }
            }

            // 判断游戏是否结束
            bool haveSpace = false;
            foreach (object t in GameSteps)
            {
                int[] row = (int[])t;
                foreach (int t1 in row)
                {
                    if (t1 == 0)
                    {
                        haveSpace = true;
                        break;
                    }
                }
                if (haveSpace)
                {
                    break;
                }
            }
            GameState = haveSpace ? 1 : 2;

            return 0;
        }

        private void PlusSorce(int player)
        {
            if (player == 1)
            {
                Player1Sorce++;
            }
            else if (player == 2)
            {
                Player2Sorce++;
            }
        }

        public bool IsEndGame()
        {
            return true;
        }

        private void InitSteps(int type)
        {
            GameSteps = new Object[type * 2 - 1];
            for (int i = 0; i < type * 2 - 1; i++)
            {
                int x = type - 1;
                if (i != 0)
                {
                    x = i % 2 == 0 ? type - 1 : type;
                }
                GameSteps[i] = new int[x];
            }
        }

    }
}
