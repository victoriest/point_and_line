using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace connectToGoServer
{
    public class LineAndPointGame
    {
        // 玩家1的得分
        public int player1Sorce { get; set; }

        // 玩家2的得分
        public int player2Sorce { get; set; }

        // 棋盘大小 1:3*3, 2:5*5, 3:7*7
        public int gameType { get; set; }

        // 棋盘数据
        public Object[] gameSteps { get; set; }

        // 轮到谁出手
        public int whosTurn { get; set; }

        // 游戏状态 1:进行中; 2:已结束
        public int gameState { get; set; }

        public LineAndPointGame(int type)
        {
            player1Sorce = 0;
            player2Sorce = 0;
            gameType = type;
            whosTurn = 1;
            gameState = 1;
            if (type == 1)
            {
                initSteps(3);
            }
            else if(type == 2)
            {
                initSteps(5);
            }
            else if (type == 3)
            {
                initSteps(7);
            }
        }

        public int Line(int rowIndex, int colIndex, int playerId)
        {
            if (playerId != whosTurn)
            {
                // 不轮到你走
                return 1;
            }
            int[] obj = (int[])gameSteps[rowIndex];
            int state = obj[colIndex];
            if (state != 0)
            {
                // 这里不准走魂淡
                return 2;
            }

            // 能走就走上
            obj[colIndex] = whosTurn;

            // 判断是否得分
            bool isSorce = false;
            if (rowIndex == 0 || rowIndex % 2 == 0)
            {
                // 偶数
                if (rowIndex + 2 < gameSteps.Length)
                {
                    // 下面的方块
                    var bottomBottom = ((int[])gameSteps[rowIndex + 2])[colIndex];
                    var bottomLeft = ((int[])gameSteps[rowIndex + 1])[colIndex];
                    var bottomRight = ((int[])gameSteps[rowIndex + 1])[colIndex + 1];
                    if (bottomBottom != 0 && bottomLeft != 0 && bottomRight != 0)
                    {
                        plusSorce(playerId);
                        isSorce = true;
                    }
                }
                if (rowIndex - 2 >= 0)
                {
                    // 上面的方块
                    var topLeft = ((int[])gameSteps[rowIndex - 1])[colIndex];
                    var topTop = ((int[])gameSteps[rowIndex - 2])[colIndex];
                    var topRight = ((int[])gameSteps[rowIndex - 1])[colIndex + 1];
                    if (topLeft != 0 && topTop != 0 && topRight != 0)
                    {
                        plusSorce(playerId);
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
                    var leftLeft = ((int[])gameSteps[rowIndex])[colIndex - 1];
                    var leftTop = ((int[])gameSteps[rowIndex - 1])[colIndex - 1];
                    var leftBottom = ((int[])gameSteps[rowIndex + 1])[colIndex - 1];
                    if (leftLeft != 0 && leftTop != 0 && leftBottom != 0)
                    {
                        plusSorce(playerId);
                        isSorce = true;
                    }
                }
                if (colIndex + 1 < ((int[])gameSteps[rowIndex]).Length)
                {
                    // 右边的方块
                    var rightRight = ((int[])gameSteps[rowIndex])[colIndex];
                    var rightTop = ((int[])gameSteps[rowIndex - 1])[colIndex];
                    var rightBottom = ((int[])gameSteps[rowIndex + 1])[colIndex];
                    if (rightRight != 0 && rightTop != 0 && rightBottom != 0)
                    {
                        plusSorce(playerId);
                        isSorce = true;
                    }
                }

            }

            // 判断该谁走
            if (isSorce)
            {
                whosTurn = playerId;
            }
            else
            {
                if (playerId == 1)
                {
                    whosTurn = 2;
                }
                else if (playerId == 2)
                {
                    whosTurn = 1;
                }
            }

            // 判断游戏是否结束
            bool haveSpace = false;
            for (int i = 0; i < gameSteps.Length; i++)
            {
                int[] row = (int[])gameSteps[i];
                for (int j = 0; j < row.Length; j++)
                {
                    if (row[j] == 0)
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
            if (haveSpace)
            {
                gameState = 1;
            }
            else
            {
                gameState = 2;
            }

            return 0;
        }

        private void plusSorce(int player)
        {
            if (player == 1)
            {
                player1Sorce++;
            }
            else if (player == 2)
            {
                player2Sorce++;
            }
        }

        public int getWhosTurn()
        {
            return whosTurn;
        }

        public bool isEndGame()
        {
            return true;
        }

        private void initSteps(int type)
        {
            gameSteps = new Object[type * 2 - 1];
            for (int i = 0; i < type * 2 - 1; i++)
            {
                int x = type - 1;
                if (i != 0)
                {
                    x = i % 2 == 0 ? type - 1 : type;
                }
                gameSteps[i] = new int[x];
            }
        }

    }
}
