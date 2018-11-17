import Pool from './base/pool'

let instance

/**
 * 全局状态管理器
 */
export default class LineAndPointGame {
  constructor() {
    if ( instance )
      return instance

    instance = this

    this.pool = new Pool()
    this.player1Score = 0
    this.player2Score = 0
    this.gameType = 2
    this.gameSteps = []
    this.whosTurn = 1
    // -1: 初始状态; 0: 寻找中; 1: 正在进行, 2: 已结束
    this.gameState = -1

    this.reset(2)
  }

  reset(gt = 2) {
    this.player1Score = 0
    this.player2Score = 0
    this.gameType = gt
    this.gameSteps = []
    this.whosTurn = 1
    this.gameState = -1

    if (gt == 1) {
      this.initSteps(3);
    }
    else if (gt == 2) {
      this.initSteps(5);
    }
    else if (gt == 3) {
      this.initSteps(7);
    }
  }

  drawALine(rowIndex, colIndex, playerId) {
    // not your turn
    if(playerId != this.whosTurn) {
      return 1
    }
    var obj = this.gameSteps[rowIndex]
    var state = obj[colIndex]

    // cant step here
    if(state != 0) {
      return 2
    }

    // put the step
    obj[colIndex] = this.whosTurn

    var isScore = false
    if(rowIndex == 0 || rowIndex % 2 == 0) {
      // even
      if(rowIndex + 2 < this.gameSteps.length) {
        var bottomBottom = this.gameSteps[rowIndex + 2][colIndex]
        var bottomLeft = this.gameSteps[rowIndex + 1][colIndex]
        var bottomRight = this.gameSteps[rowIndex + 1][colIndex]
        if(bottomBottom != 0 && bottomLeft != 0 && bottomRight != 0 && 
          bottomBottom == bottomLeft && bottomLeft == bottomRight && this.whosTurn == bottomRight) {
          this.plusScore(playerId)
          isScore = true
        }
      }
      if(rowIndex = 2 >= 0) {
        var topLeft = this.gameSteps[rowIndex - 1][colIndex]
        var topTop = this.gameSteps[rowIndex - 2][colIndex]
        var topRight = this.gameSteps[rowIndex - 1][colIndex + 1]
        if(topLeft != 0 && topTop != 0 && topRight != 0 &&
          topLeft == topTop && topTop == topRight && this.whosTurn == topRight) {
          this.plusScore(playerId)
          isScore = true
        }
      }
    }
    else {
      // odd
      if(colIndex - 1 >= 0) {
        var leftLeft = this.gameSteps[rowIndex][colIndex - 1]
        var leftTop = thsi.gameSteps[rowIndex - 1][colIndex - 1]
        var leftBottom = this.gameSteps[rowIndex + 1][colIndex - 1]
        if(leftLeft != 0 && leftTop != 0 && leftBottom != 0 && 
          leftLeft == leftTop && leftTop == leftBottom && this.whosTurn == leftBottom) {
          this.plusScore(playerId)
          isScore = true
        }
      }
      if(colIndex + 1 < this.gameSteps[rowIndex].length) {
        var rightRight = this.gameSteps[rowIndex][colIndex + 1]
        var rightTop = this.gameSteps[rowIndex - 1][colIndex]
        var rightBottom = this.gameSteps[rowIndex + 1][colIndex]
        if (rightRight != 0 && rightTop != 0 && rightBottom != 0 && 
          rightRight == rightTop && rightTop == rightBottom && this.whosTurn == rightBottom) {
          PlusSorce(playerId);
          isSorce = true;
        }
      }
    }

    // 判断该谁走
    if (isSorce) {
      this.whosTurn = playerId
    }
    else {
      if (playerId == 1) {
        this.whosTurn = 2
      }
      else if (playerId == 2) {
        this.whosTurn = 1
      }
    }

    // 判断游戏是否结束
    var haveSpace = false
    for(var i = 0; i < this.gameSteps.length; i++) {
      var row = this.gameSteps[i]
      for(var j = 0; j < row.length; j++) {
        var cell = row[j]
        if(cell == 0) {
          haveSpace = true
          break
        }
      }
      if(haveSpace) break
    }

    this.gameState = haveSpace ? 1 : 2

    return 0
  }

  plusScore(player) {
    if(player == 1) {
      this.player1Score++
    }
    else if(player == 2) {
      this.player2Score++
    }
  }

  initSteps(boardLength) {
    var arrLength = boardLength * 2 - 1
    this.gameSteps = new Array(arrLength)
    for(var i = 0; i < arrLength; i++) {
      var x = boardLength - 1
      if(i != 0) {
        x = i % 2 == 0 ? boardLength - 1 : boardLength
      }
      this.gameSteps[i] = new Array(x)
    }
  }

}
