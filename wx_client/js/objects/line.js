import Sprite from '../base/sprite'
import LineAndPointGame from '../lineAndPointGame'

const screenWidth = window.innerWidth
const screenHeight = window.innerHeight

const STYLE_0_IMG = {x: 486, y: 9, width: 6, height: 56}
const STYLE_1_IMG = {x: 436, y: 9, width: 15, height: 56}
const STYLE_2_IMG = {x: 461, y: 9, width: 15, height: 56}

export default class Line {
  constructor(x, y, width, height, gameType = 2) {
    this.img = new Image()
    this.img.src = 'images/Common.png'
    this.width = width
    this.height = height
    this.x = x
    this.y = y
    this.gameType = gameType
  }

  isClicked(clickedX, clickedY) {
    // TODO
  }

  drawToCanvas(ctx, style) {
    var whatToDraw = STYLE_0_IMG
    if (style == 1) {
      whatToDraw = STYLE_1_IMG
    } else if (style == 2) {
      whatToDraw = STYLE_2_IMG
    }
    ctx.drawImage(
      this.img,
      whatToDraw.x, whatToDraw.y, whatToDraw.width, whatToDraw.height,
      this.x, this.y, this.width, this.height
    )
  }

}
