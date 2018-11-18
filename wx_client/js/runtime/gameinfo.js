import LineAndPointGame from '../lineAndPointGame'
import Line from '../objects/line'

const screenWidth  = window.innerWidth
const screenHeight = window.innerHeight

let atlas = new Image()
atlas.src = 'images/Common.png'

let GAME_BOARD = null

export default class GameInfo {
  renderGameScore(ctx, score) {
    ctx.fillStyle = "#ffffff"
    ctx.font      = "20px Arial"

    ctx.fillText(
      score,
      10,
      30
    )
  }

  checkWitchLineHasBeenClicked(x, y) {
    if (GAME_BOARD == null) {
      return null
    }

    for (var i = 0; i < GAME_BOARD.length; i++) {
      var rows = GAME_BOARD[i].length
      for (var j = 0; j < rows; j++) {
        var line = GAME_BOARD[i][j]
        if(line.isClicked(x, y)){
          return [i, j]
        }
      }
    }
    return null
  }

  renderGameBoard(ctx, game) {
    if (GAME_BOARD == null) {
      GAME_BOARD = new Array(game.gameSteps.length)
      for (var i = 0; i < game.gameSteps.length; i++) {
        GAME_BOARD[i] = new Array(game.gameSteps[i].length)
      }
      var startX = screenWidth * 0.05
      var startY = screenHeight / 2 - screenWidth / 2
      var w = screenWidth / (game.gameSteps.length / 2)
      var h = w / 5
      var rows = game.gameSteps.length
      for (var i = 0; i < rows; i++) {
        var seed = i / 2
        var colObj = game.gameSteps[i]
        var cols = colObj.length
        for (var j = 0; j < cols; j++) {
          var line = null
          // 偶数
          if (i == 0 || i % 2 == 0) {
            line = new Line(startX + j * w + h, startY + seed * w, w - h, h)
          }
          // 奇数
          else {
            line = new Line(startX + j * w, startY + seed * w - h, h, w - h)
          }
          GAME_BOARD[i][j] = line
        }
      }
    }

    for (var i = 0; i < GAME_BOARD.length; i++) {
      var rows = GAME_BOARD[i].length
      for (var j = 0; j < rows; j++) {
        var line = GAME_BOARD[i][j]
        line.drawToCanvas(ctx, game.gameSteps[i][j])
      }
    }
  }

  reanderGameStart(ctx) {
    ctx.drawImage(atlas, 0, 0, 119, 108, screenWidth / 2 - 150, screenHeight / 2 - 100, 300, 300)

    ctx.fillStyle = "#ffffff"
    ctx.font = "20px Arial"

    ctx.fillText(
      'Point And Line',
      screenWidth / 2 - 55,
      screenHeight / 2 - 100 + 50
    )

    ctx.drawImage(
      atlas,
      120, 6, 39, 24,
      screenWidth / 2 - 60,
      screenHeight / 2 - 100 + 180,
      120, 40
    )

    ctx.fillText(
      '寻找对手',
      screenWidth / 2 - 40,
      screenHeight / 2 - 100 + 205
    )

    /**
     * 寻找对手按钮区域
     */
    this.btnArea = {
      startX: screenWidth / 2 - 40,
      startY: screenHeight / 2 - 100 + 180,
      endX: screenWidth / 2 + 50,
      endY: screenHeight / 2 - 100 + 255
    }
  }

  reanderSearch(ctx) {
    ctx.drawImage(atlas, 0, 0, 119, 108, screenWidth / 2 - 150, screenHeight / 2 - 100, 300, 300)

    ctx.fillStyle = "#ffffff"
    ctx.font = "20px Arial"

    ctx.fillText(
      'Point And Line',
      screenWidth / 2 - 55,
      screenHeight / 2 - 100 + 50
    )

    ctx.fillText(
      '正在寻找...',
      screenWidth / 2 - 40,
      screenHeight / 2 - 100 + 130
    )

    ctx.drawImage(
      atlas,
      120, 6, 39, 24,
      screenWidth / 2 - 60,
      screenHeight / 2 - 100 + 180,
      120, 40
    )

    ctx.fillText(
      '取消搜索',
      screenWidth / 2 - 40,
      screenHeight / 2 - 100 + 205
    )

    /**
     * 寻找对手按钮区域
     */
    this.btnArea = {
      startX: screenWidth / 2 - 40,
      startY: screenHeight / 2 - 100 + 180,
      endX: screenWidth / 2 + 50,
      endY: screenHeight / 2 - 100 + 255
    }
  }

  renderGameOver(ctx, score) {
    GAME_BOARD = null
    ctx.drawImage(atlas, 0, 0, 119, 108, screenWidth / 2 - 150, screenHeight / 2 - 100, 300, 300)

    ctx.fillStyle = "#ffffff"
    ctx.font    = "20px Arial"

    ctx.fillText(
      '游戏结束',
      screenWidth / 2 - 40,
      screenHeight / 2 - 100 + 50
    )

    ctx.fillText(
      '得分: ' + score,
      screenWidth / 2 - 40,
      screenHeight / 2 - 100 + 130
    )

    ctx.drawImage(
      atlas,
      120, 6, 39, 24,
      screenWidth / 2 - 60,
      screenHeight / 2 - 100 + 180,
      120, 40
    )

    ctx.fillText(
      '再来一局',
      screenWidth / 2 - 40,
      screenHeight / 2 - 100 + 205
    )

    /**
     * 重新开始按钮区域
     * 方便简易判断按钮点击
     */
    this.btnArea = {
      startX: screenWidth / 2 - 40,
      startY: screenHeight / 2 - 100 + 180,
      endX  : screenWidth / 2  + 50,
      endY  : screenHeight / 2 - 100 + 255
    }
  }
}

