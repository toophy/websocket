package mud

import (
	"time"
)

// Match
type Match struct {
	Name    string
	Game    string
	Time    int32
	Step    int32   // 段位
	WinRate float32 // 胜率
	Kda     float32 // K/D/A 值 (杀/死/助), (K+A)/D, 如果D为0, 则D为1
	Elo     int32   // Elo积分
	Use     bool
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

func (m *MatchSys) AskMatch(name, game string, step int32, winRate, kda float32, elo int32) {
	if _, ok := m.Matchs[game]; !ok {
		matchSys.Matchs[game] = make(map[string]*Match, 0)
	}
	if _, ok := m.Matchs[game][name]; !ok {
		m.Matchs[game][name] = &Match{
			Name:    name,
			Game:    game,
			Time:    int32(time.Now().Unix()),
			Step:    step,
			WinRate: winRate,
			Kda:     kda,
			Elo:     elo,
			Use:     false}
	}
}

func (m *MatchSys) CancelMatch(name string) {
	if _, ok := m.Matchs[name]; ok {
		delete(m.Matchs, name)
	}
}

func (m *MatchSys) Update() {
	/*
	   优先使用 段位 找对手
	   其次 Elo 积分
	   再次 胜率
	   最后 Kda值

	   1. 段位池
	   2. Elo积分池
	*/
}
