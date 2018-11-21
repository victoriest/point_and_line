// import Player     from './player/index'
// import Enemy      from './npc/enemy'
// import Music      from './runtime/music'

import BackGround from './runtime/background'
import GameInfo   from './runtime/gameinfo'
import LineAndPointGame from './lineAndPointGame'

let CTX   = canvas.getContext('2d')
let GAME = new LineAndPointGame()

/**
 * 游戏主函数
 */
export default class Main {
  constructor() {
    // 维护当前requestAnimationFrame的id
    this.aniId    = 0
    this.restart()

  }

  restart() {
    GAME.reset()
    wx.onSocketOpen(function(data) {console.log(data)})
    wx.onSocketMessage(function (data) { console.log(data) })
    wx.onSocketClose(function (data) { console.log(data) })
    wx.connectSocket({
      url: 'ws://127.0.0.1:9090/ws'
    })

    // canvas.removeEventListener(
    //   'touchstart',
    //   this.touchHandler
    // )
    this.bg = new BackGround(CTX)
    // this.player   = new Player(CTX)
    this.gameInfo = new GameInfo()
    // this.music    = new Music()

    this.bindLoop     = this.loop.bind(this)
    this.hasEventBind = false

    // 清除上一局的动画
    window.cancelAnimationFrame(this.aniId)

    this.aniId = window.requestAnimationFrame(
      this.bindLoop,
      canvas
    )

    if (!this.hasEventBind) {
      this.hasEventBind = true
      this.touchHandler = this.touchEventHandler.bind(this)
      canvas.addEventListener('touchstart', this.touchHandler)
    }
  }

  // 游戏结束后的触摸事件处理逻辑
  touchEventHandler(e) {
    console.log(e)
     e.preventDefault()

    let x = e.touches[0].clientX
    let y = e.touches[0].clientY

    if(GAME.gameState == 1) {
      var s = this.gameInfo.checkWitchLineHasBeenClicked(x, y)
      if(s != null) {
        GAME.drawALine(s[0], s[1], GAME.whosTurn)
        console.log(s, GAME.player1Score, GAME.player2Score)
      }
      return
    }

    var area = this.gameInfo.btnArea
    var clicked = false
    if (x >= area.startX && x <= area.endX && y >= area.startY && y <= area.endY) {
          clicked = true
    }
    if(!clicked) return

    switch (GAME.gameState) {
      case -1:

        protobuf.load("./MobileSuite.json", function (err, root) {
          if (err) throw err;
          // Obtain a message type
          var MobileSuiteModel = root.lookupType("protocol.MobileSuiteModel");
          var LoginDTO = root.lookupType("protocol.LoginDTO");

          var loginDto = LoginDTO.create({ userId: 1, uName: "est", pwd: "123123" });
          var ms = { "type": 103, message: loginDto }
          wx.sendSocketMessage({
            data: ms,
          });
        });
        // GAME.gameState = 0
        break;
      case 0:
        GAME.gameState = 1
        break;
      case 1:
        break;
      case 2:
        GAME.gameState = -1
        break;
    }
  }

  /**
   * canvas重绘函数
   * 每一帧重新绘制所有的需要展示的元素
   */
  render() {
    CTX.clearRect(0, 0, canvas.width, canvas.height)
    // console.log(GAME)

    switch (GAME.gameState) {
      case -1:
        this.gameInfo.reanderGameStart(CTX)
        break;
      case 0:
        this.gameInfo.reanderSearch(CTX)
        break;
      case 1:
        // this.bg.render(CTX)
        // this.gameInfo.renderGameScore(CTX, GAME.player1Score)
        this.gameInfo.renderGameBoard(CTX, GAME)
        break;
      case 2:
        this.gameInfo.renderGameOver(CTX, GAME.player1Score)
        break;
      default:
    }


    // databus.bullets
    //       .concat(databus.enemys)
    //       .forEach((item) => {
    //           item.drawToCanvas(ctx)
    //         })

    // this.player.drawToCanvas(ctx)

    // databus.animations.forEach((ani) => {
    //   if ( ani.isPlaying ) {
    //     ani.aniRender(ctx)
    //   }
    // })

    // this.gameinfo.renderGameScore(ctx, databus.score)

    // // 游戏结束停止帧循环
    // if ( databus.gameOver ) {
    //   this.gameinfo.renderGameOver(ctx, databus.score)

    //   if ( !this.hasEventBind ) {
    //     this.hasEventBind = true
    //     this.touchHandler = this.touchEventHandler.bind(this)
    //     canvas.addEventListener('touchstart', this.touchHandler)
    //   }
    // }
  }

  // 实现游戏帧循环
  loop() {
    // databus.frame++
    // this.update()
    this.render()

    this.aniId = window.requestAnimationFrame(
      this.bindLoop,
      canvas
    )
  }
}

  // 游戏逻辑更新主函数
  // update() {
    // console.log("update")
    // if ( databus.gameOver )
    //   return;

    // this.bg.update()

    // databus.bullets
    //        .concat(databus.enemys)
    //        .forEach((item) => {
    //           item.update()
    //         })

    // this.enemyGenerate()

    // this.collisionDetection()

    // if ( databus.frame % 20 === 0 ) {
    //   this.player.shoot()
    //   this.music.playShoot()
    // }
  // }

/**
 * 随着帧数变化的敌机生成逻辑
 * 帧数取模定义成生成的频率
 */
// enemyGenerate() {
//   // if ( databus.frame % 30 === 0 ) {
//   //   let enemy = databus.pool.getItemByClass('enemy', Enemy)
//   //   enemy.init(6)
//   //   databus.enemys.push(enemy)
//   // }
// }

  // // 全局碰撞检测
  // collisionDetection() {
  //   let that = this

  //   databus.bullets.forEach((bullet) => {
  //     for ( let i = 0, il = databus.enemys.length; i < il;i++ ) {
  //       let enemy = databus.enemys[i]

  //       if ( !enemy.isPlaying && enemy.isCollideWith(bullet) ) {
  //         enemy.playAnimation()
  //         that.music.playExplosion()

  //         bullet.visible = false
  //         databus.score  += 1

  //         break
  //       }
  //     }
  //   })

  //   for ( let i = 0, il = databus.enemys.length; i < il;i++ ) {
  //     let enemy = databus.enemys[i]

  //     if ( this.player.isCollideWith(enemy) ) {
  //       databus.gameOver = true

  //       break
  //     }
  //   }
  // }