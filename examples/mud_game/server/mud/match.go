package mud

import (
	"math"
	"time"
)

// Match
type Match struct {
	Name  string
	Game  string
	Time  int32
	Score int64 // 复合积分, step*10^10 + elo*10^6 + kda*10^3 + winRate
	Use   bool
}

type MatchSys struct {
	Matchs map[string]map[string]*Match
}

var (
	matchSys *MatchSys
)

func init() {
	matchSys = &MatchSys{}
	matchSys.Matchs = make(map[string]map[string]*Match, 0)
}

func GetMatchSys() *MatchSys {
	return matchSys
}

func (m *MatchSys) AskMatch(name, game string, step, elo int32, winRate, kda float64) {
	if _, ok := m.Matchs[game]; !ok {
		matchSys.Matchs[game] = make(map[string]*Match, 0)
	}
	if _, ok := m.Matchs[game][name]; !ok {
		m.Matchs[game][name] = &Match{
			Name:  name,
			Game:  game,
			Time:  int32(time.Now().Unix()),
			Score: (int64(step) * 10000000000) + (int64(elo) * 1000000) + (int64(math.Ceil(kda*10)) * 1000) + (int64(math.Ceil(winRate * 10))),
			Use:   false}
	}
}

func (m *MatchSys) CancelMatch(name string) {
	if _, ok := m.Matchs[name]; ok {
		delete(m.Matchs, name)
	}
}

func (m *MatchSys) matchGame(gameType string, accName string, count int) (ret []string) {
	if count < 2 {
		return
	}

	/*
	   优先使用 段位 找对手
	   其次 Elo 积分
	   再次 胜率
	   最后 Kda值
	   1. 段位池
	   2. Elo积分池
	   3. KDA
	   4. 胜率
	*/

	if _, ok := m.Matchs[gameType]; ok {

		gamePlayerCount := len(m.Matchs[gameType])
		if count > gamePlayerCount {
			return
		}

		if _, ok2 := m.Matchs[gameType][accName]; ok2 {
			// Time  int32
			for k := range m.Matchs[gameType] {
				if k != accName {
					ret = append(ret, k)
				}
			}
			if len(ret) < count-1 {
				ret = []string{}
			} else {
				ret = append(ret, accName)
			}
		}
	}
	return
}

func (m *MatchSys) Update() {
	for k := range m.Matchs {
		for {
			loopOver := true
			for j := range m.Matchs[k] {
				ret := m.matchGame(k, j, 2)
				if len(ret) >= 2 {
					// 成功匹配一个房间, 房间内的所有匹配删除, 房间生成, 制作一个房间, 并发放给hall去处理房间问题
					go GetHall().OnMatchOver(ret)
					for i := 0; i < len(ret); i++ {
						delete(m.Matchs[k], ret[i])
					}
					loopOver = false
					break
				}
			}
			if loopOver {
				break
			}
		}
	}
}
