package mud

import (
	"math"
)

const (
	deno    = 400
	novice  = 1000.0 //0-999 初学者
	someExp = 1500.0 //1000-1499 有经验者
	skill   = 2000.0 //1500-1999 熟练者
	expert  = 2200.0 //2000-2199 专家
	master  = 2400.0 //2200-2399 大师 >2400 宗师
)

/**
* @brief         自适应k值
*
* @param score   积分
* @param isWin   胜负因子
*
* @returns       k值结果
 */
func AdaptationK(score float64, isWin bool) (ret uint) {
	ret = 10
	if score < novice {
		if isWin {
			ret = 32 + 32
		} else {
			ret = 32
		}
	} else if score < someExp {
		if isWin {
			ret = 32 + 16
		} else {
			ret = 32
		}
	} else if score < skill {
		ret = 32
	} else if score < expert {
		ret = 20
	} else if score < master {
		ret = 15
	}
	return
}

/**
* @brief         计算积分结果
*
* @param a       a的积分
* @param b       b的积分
* @param isWin   a对b的胜负关系 true-a胜 false-a败
* @param isDraw  a对b是否平局   true-平（不查看isWin字段） false-有胜负（查看isWin字段）
*
* @returns       a的积分终值
 */
func CalcResult(a, b float64, isWin, isDraw bool) float64 {
	return a + GetScoreChg(a, b, isWin, isDraw)
}

/**
* @brief    计算a相对b的胜率
*
* @param a  a的积分
* @param b  b的积分
*
* @returns  胜率
 */
func GetWinRate(a, b float64) float64 {
	return 1 / (1 + math.Pow(10, (b-a)/deno))
}

/**
* @brief        计算积分变化
*
* @param a      a的积分
* @param b      b的积分
* @param isWin  a对b的胜负关系 true-a胜 false-a败
* @param isDraw a对b是否平局   true-平（不查看isWin字段） false-有胜负（查看isWin字段）
*
* @returns      a的积分变化
 */
func GetScoreChg(a, b float64, isWin, isDraw bool) float64 {
	w := 0.5
	if !isDraw {
		if isWin {
			w = 1.0
		} else {
			w = 0.0
		}
	}
	winRate := GetWinRate(a, b)
	return (w - winRate) * float64(AdaptationK(a, isWin))
}

// GetKDA 通过k/D/A值获取KDA值
func GetKDA(kill, dead, assist uint32) float32 {
	if dead == 0 {
		dead = 1
	}
	return float32(kill+assist) / float32(dead)
}
