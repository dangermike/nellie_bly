package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/dangermike/nelly_bly/die"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gonum.org/v1/gonum/stat"
)

const (
	giveUp       = 200
	targetTrials = 1 << 20
)

type popItem struct {
	ix  int
	cnt int
}

type statSummary struct {
	Min     float64
	Ntile25 float64
	Ntile50 float64
	Ntile75 float64
	Ntile90 float64
	Ntile95 float64
	Ntile99 float64
	Max     float64
	Mean    float64
	StdDev  float64
}

type trialSummary struct {
	totalGames int
	totalTurns int
	players    int
	turns      *statSummary
	captures   *statSummary
	popular    []popItem
}

func newTrialSummary(
	players int,
	totalGames int,
	totalTurns int,
	minTurns int,
	maxTurns int,
	turnWeights []float64,
	capWeights []float64,
	popularity []popItem,
) *trialSummary {
	turnMean, turnStdev := stat.MeanStdDev(counts, turnWeights)
	capMean, capStdev := stat.MeanStdDev(counts, capWeights)
	safePop := make([]popItem, len(popularity))
	copy(safePop, popularity)
	// reverse-sort
	sort.Slice(safePop, func(i int, j int) bool {
		return safePop[i].cnt > safePop[j].cnt
	})

	return &trialSummary{
		players:    players,
		totalGames: totalGames,
		totalTurns: totalTurns,
		turns: &statSummary{
			Min:     float64(minTurns),
			Ntile25: stat.Quantile(0.25, stat.Empirical, counts, turnWeights),
			Ntile50: stat.Quantile(0.5, stat.Empirical, counts, turnWeights),
			Ntile75: stat.Quantile(0.75, stat.Empirical, counts, turnWeights),
			Ntile90: stat.Quantile(0.90, stat.Empirical, counts, turnWeights),
			Ntile95: stat.Quantile(0.95, stat.Empirical, counts, turnWeights),
			Ntile99: stat.Quantile(0.99, stat.Empirical, counts, turnWeights),
			Max:     float64(maxTurns),
			Mean:    turnMean,
			StdDev:  turnStdev,
		},
		captures: &statSummary{
			Min:     stat.Quantile(0.0, stat.Empirical, counts, capWeights),
			Ntile25: stat.Quantile(0.25, stat.Empirical, counts, capWeights),
			Ntile50: stat.Quantile(0.5, stat.Empirical, counts, capWeights),
			Ntile75: stat.Quantile(0.75, stat.Empirical, counts, capWeights),
			Ntile90: stat.Quantile(0.90, stat.Empirical, counts, capWeights),
			Ntile95: stat.Quantile(0.95, stat.Empirical, counts, capWeights),
			Ntile99: stat.Quantile(0.99, stat.Empirical, counts, capWeights),
			Max:     stat.Quantile(1.0, stat.Empirical, counts, capWeights),
			Mean:    capMean,
			StdDev:  capStdev,
		},
		popular: safePop,
	}
}

// Player is one who plays
type Player struct {
	turns    int
	position int
	skips    int // pending lost turns
	captured int
}

func (p *Player) Reset() {
	p.turns = 0
	p.position = 0
	p.skips = 0
	p.captured = 0
}

func makeCounts(cnt int) []float64 {
	ret := make([]float64, cnt)
	for i := 0; i < cnt; i++ {
		ret[i] = float64(i)
	}
	return ret
}

var counts = makeCounts(giveUp + 1)

func main() {
	summaries := []*trialSummary{}
	p := message.NewPrinter(getLanguage())

	for cycle := 0; cycle < 5; cycle++ {
		playerCnt := 1 << uint(cycle*2)
		// playerCnt := 1 + cycle
		trials := targetTrials / playerCnt
		var totalGames int
		var totalTurns int
		minTurns := giveUp + 1
		maxTurns := 0
		turnWeights := make([]float64, giveUp+1)
		capWeights := make([]float64, giveUp+1)
		popularity := make([]popItem, len(nellyboard))
		for i := 0; i < len(popularity); i++ {
			popularity[i] = popItem{ix: i}
		}
		var globalLock sync.Mutex // this is how we are going to update those global things safely

		chTrials := make(chan int)
		var wg sync.WaitGroup

		for t := 0; t < runtime.NumCPU()*2; t++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				dieSix := die.New(6)
				players := make([]*Player, playerCnt)
				for i := 0; i < playerCnt; i++ {
					players[i] = &Player{}
				}
				checks := make([]int, 2) // squares to check
				localGames := 0
				localTurns := 0
				localMinTurns := giveUp + 1
				localMaxTurns := 0
				localTurnWeights := make([]float64, giveUp+1)
				localCapWeights := make([]float64, giveUp+1)
				localPopularity := make([]int, len(nellyboard))
				for range chTrials {
					for i := 0; i < playerCnt; i++ {
						players[i].Reset()
					}

					var winner *Player = nil
					for turn := 0; winner == nil; turn++ {
						player := players[turn%playerCnt]
						checks = checks[:0]

						player.turns++
						if player.skips > 0 {
							player.skips--
							continue
						}

						newPos := player.position + dieSix.Roll()
						if newPos >= len(nellyboard) {
							continue
						}
						player.position = newPos
						localPopularity[player.position]++
						if !nellyboard[player.position].IsSafe() {
							checks = append(checks, player.position)
						}

						// according to the rules, you only follow the offset once
						if nellyboard[player.position].Offset != 0 {
							player.position += nellyboard[player.position].Offset
							localPopularity[player.position]++
							if !nellyboard[player.position].IsSafe() {
								checks = append(checks, player.position)
							}
						}
						for _, ix := range checks {
							if !nellyboard[ix].IsSafe() {
								for _, p := range players {
									if p.position == newPos && p != player {
										p.captured++
										p.position = 0
									}
								}
							}
						}

						// handle turn-altering instructions
						if nellyboard[player.position].Turns > 0 {
							// skip me n times
							player.skips += nellyboard[player.position].Turns
						} else if nellyboard[player.position].Turns < 0 {
							// skip everyone else n times
							for _, p := range players {
								if p != player {
									p.skips += nellyboard[player.position].Turns
								}
							}
						}
						if player.position == len(nellyboard)-1 {
							winner = player
						}
					}
					localGames++
					localTurns += winner.turns
					localTurnWeights[winner.turns] += 1
					if winner.turns < localMinTurns {
						localMinTurns = winner.turns
					}
					if winner.turns > localMaxTurns {
						localMaxTurns = winner.turns
					}
					for _, p := range players {
						localCapWeights[p.captured]++
					}
				}

				globalLock.Lock()
				totalGames += localGames
				totalTurns += localTurns
				if localMinTurns < minTurns {
					minTurns = localMinTurns
				}
				if localMaxTurns > maxTurns {
					maxTurns = localMaxTurns
				}
				for i := 0; i < giveUp+1; i++ {
					turnWeights[i] += localTurnWeights[i]
					capWeights[i] += localCapWeights[i]
				}
				for i := 0; i < len(nellyboard); i++ {
					popularity[i].cnt += localPopularity[i]
				}
				globalLock.Unlock()
			}()
		}

		for i := 0; i < trials; i++ {
			chTrials <- i
		}
		close(chTrials)
		wg.Wait()
		ts := newTrialSummary(playerCnt, totalGames, totalTurns, minTurns, maxTurns, turnWeights, capWeights, popularity)
		summaries = append(summaries, ts)
		p.Printf("%d games with %d players took %d turns (%0.3f avg, %0.3f stdev)\n",
			ts.totalGames,
			ts.players,
			ts.totalTurns,
			ts.turns.Mean,
			ts.turns.StdDev,
		)
	}
	p.Println()
	p.Println("turns")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{s.players}
	}, "         ", "% 3d player", "\n")
	p.Print("---------")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{"---------"}
	}, "", " %s", "\n")
	formatSummaryStats(p, summaries, func(s *trialSummary) *statSummary {
		return s.turns
	})
	p.Println()
	p.Println("captures")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{s.players}
	}, "         ", "% 3d player", "\n")
	p.Print("---------")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{"---------"}
	}, "", " %s", "\n")
	formatSummaryStats(p, summaries, func(s *trialSummary) *statSummary {
		return s.captures
	})

	p.Println()
	p.Println("most popular squares")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{fmt.Sprintf("%d player", s.players)}
	}, "   ", " %-40s", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{"----------------------------------------"}
	}, "   ", " %s", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{fmt.Sprintf("%11s %3s %-20s %3s", "cnt", "loc", "name", "ofs")}
	}, "   ", " %-40s", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{"----------- --- -------------------- ---"}
	}, "   ", " %s", "\n")
	for i := 0; i < 20; i++ {
		summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
			square := nellyboard[s.popular[i].ix]
			return []interface{}{
				s.popular[i].cnt,
				square.Day,
				square.Name,
				square.Offset,
			}
		}, p.Sprintf("%-2d ", i+1), " %11d %3d %-20s %+3d", "\n")
	}
}

func getLanguage() language.Tag {
	s, ok := os.LookupEnv("LANGUAGE")
	if !ok {
		return language.English
	}
	lang, err := language.Parse(s)
	if err != nil {
		return language.English
	}
	return lang
}

func summariesPrintAll(p *message.Printer, summaries []*trialSummary, selector func(*trialSummary) []interface{}, prefix, format, suffix string) {
	p.Print(prefix)
	for _, s := range summaries {
		p.Printf(format, selector(s)...)
	}
	p.Print(suffix)
}

func formatSummaryStats(p *message.Printer, summaries []*trialSummary, selctor func(s *trialSummary) *statSummary) {
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Min}
	}, "Min:     ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Ntile25}
	}, "25th:    ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Ntile50}
	}, "50th:    ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Ntile75}
	}, "75th:    ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Ntile90}
	}, "90th:    ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Ntile95}
	}, "95th:    ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Ntile99}
	}, "99th:    ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Max}
	}, "Max:     ", "% 10.0f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).Mean}
	}, "Mean:    ", "% 10.3f", "\n")
	summariesPrintAll(p, summaries, func(s *trialSummary) []interface{} {
		return []interface{}{selctor(s).StdDev}
	}, "Stdev:   ", "% 10.3f", "\n")
}
