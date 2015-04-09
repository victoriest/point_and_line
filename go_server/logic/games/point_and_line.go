package games

type PointAndLineGame struct {
	// 玩家1的得分
	Player1Score int

	// 玩家2的得分
	Player2Score int

	// 棋盘大小 1:3*3, 2:5*5, 3:7*7
	GameType int

	// 棋盘数据
	GameSteps []interface{}

	// 轮到谁出手
	WhosTurn int

	// 游戏状态 1:进行中; 2:已结束
	GameState int
}

func NewPointAndLineGame(gameType int) *PointAndLineGame {
	game := new(PointAndLineGame)
	game.Player1Score = 0
	game.Player2Score = 0
	game.GameType = gameType
	game.WhosTurn = 1
	var result bool
	switch gameType {
	case 1:
		result = game.initSteps(3)
		break
	case 2:
		result = game.initSteps(5)
		break
	case 3:
		result = game.initSteps(7)
		break
	default:
		return nil
	}
	if !result {
		return nil
	}
	return game
}

func (self *PointAndLineGame) plusSocre(player int) {
	if player == 1 {
		self.Player1Score++
	} else if player == 2 {
		self.Player2Score++
	}
}

func (self *PointAndLineGame) initSteps(gameType int) bool {
	self.GameSteps = make([]interface{}, gameType*2-1)
	for i, _ := range self.GameSteps {
		x := gameType - 1
		if i != 0 {
			if i%2 != 0 {
				x = gameType
			}
		}
		self.GameSteps[i] = make([]int, x)
	}
	return true
}

func (self *PointAndLineGame) IsEndGame() bool {
	return self.GameState == 2
}

func (self *PointAndLineGame) Line(rowIndex int, colIndex, playerId int) int {
	if playerId != WhoTurn {
		// 不论到你走
		return 1
	}

	obj := self.GameSteps[rowIndex]
	state := obj[colIndex]
	if state != 0 {
		// 这里不能走魂淡
		return 2
	}

	// 能走就走上
	obj[colIndex] = WhoTurn

	// 判断是否得分
	isScore := false
	if rowIndex == 0 || rowIndex%2 == 0 {
		// 偶数
		if rowIndex+2 < len(self.GameSteps) {
			// 下面的方块
			bottomBottom := ([]int(self.GameSteps[rowIndex+2]))[colIndex]
			bottomLeft := ([]int(self.GameSteps[rowIndex+1]))[colIndex]
			bottomRight := ([]int(self.GameSteps[rowIndex+1]))[colIndex+1]
			if bottomBottom != 0 && bottomLeft != 0 && bottomRight != 0 && bottomBottom == bottomLeft && bottomLeft == bottomRight && WhosTurn == bottomRight {
				self.plusSocre(playerId)
				isScore = true
			}
		}
		if rowIndex-2 >= 0 {
			// 上面的方块
			topLeft := ([]int(self.GameSteps[rowIndex-1]))[colIndex]
			topTop := ([]int(self.GameSteps[rowIndex-2]))[colIndex]
			topRight := ([]int(self.GameSteps[rowIndex-1]))[colIndex+1]
			if topLeft != 0 && topTop != 0 && topRight != 0 && topLeft == topTop && topTop == topRight && WhosTurn == topRight {
				self.plusSocre(playerId)
				isScore = true
			}
		}
	} else {
		// 奇数
		if colIndex-1 >= 0 {
			// 左边的方块
			leftLeft := ([]int(self.GameSteps[rowIndex]))[colIndex-1]
			leftTop := ([]int(self.GameSteps[rowIndex-1]))[colIndex-1]
			leftBottom := ([]int(self.GameSteps[rowIndex+1]))[colIndex-1]
			if leftLeft != 0 && leftTop != 0 && leftBottom != 0 && leftLeft == leftTop && leftTop == leftBottom && WhosTurn == leftBottom {
				self.plusSocre(playerId)
				isScore = true
			}
		}
		if colIndex+1 < len([]int(self.GameSteps[rowIndex])) {
			// 右边的方块
			rightRight := ([]int(self.GameSteps[rowIndex]))[colIndex+1]
			rightTop := ([]int(self.GameSteps[rowIndex-1]))[colIndex]
			rightBottom := ([]int(self.GameSteps[rowIndex+1]))[colIndex]
			if rightRight != 0 && rightTop != 0 && rightBottom != 0 && rightRight == rightTop && rightTop == rightBottom && WhosTurn == rightBottom {
				self.plusSocre(playerId)
				isScore = true
			}
		}
	}

	// 判断该谁走
	if isScore {
		self.WhosTurn = playerId
	} else {
		if playerId == 1 {
			self.WhosTurn = 2
		} else if playerId == 2 {
			self.WhosTurn = 1
		}
	}

	// 判断游戏是否结束
	haveSpace := false
	for t, _ := range self.GameSteps {
		row := []int(t)
		for t1, _ := range row {
			if t1 == 0 {
				haveSpace = true
				break
			}
		}
		if haveSpace {
			break
		}
	}
	if haveSpace {
		self.GameState = 1
	} else {
		self.GameState = 2
	}

	return 0

}
