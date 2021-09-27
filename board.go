package main

// Space is a position on the game board
type Space struct {
	Name   string
	Offset int
	Turns  int
	Day    int
}

func (s Space) IsSafe() bool {
	return s.Offset == 0 && s.Turns == 0
}

// Board is a game board
type Board []Space

var nellyboard = Board{
	Space{Name: "Start", Offset: 0, Turns: 0, Day: 1},
	Space{Name: "Clear", Offset: 3, Turns: 0, Day: 2},
	Space{Name: "Rain", Offset: -1, Turns: 0, Day: 3},
	Space{Name: "Storm", Offset: -2, Turns: 0, Day: 4},
	Space{Name: "Icebergs", Offset: -4, Turns: 0, Day: 5},
	Space{Name: "Clearing", Offset: 2, Turns: 0, Day: 6},
	Space{Name: "Fair Sailing", Offset: 0, Turns: 0, Day: 7},
	Space{Name: "Southampton", Offset: 1, Turns: 0, Day: 8},
	Space{Name: "Amiens", Offset: 0, Turns: 1, Day: 9},
	Space{Name: "Indian Mail Accident", Offset: -5, Turns: 0, Day: 10},
	Space{Name: "Brindisi Brigands", Offset: -2, Turns: 0, Day: 11},
	Space{Name: "Mediterranean", Offset: 0, Turns: 0, Day: 12},
	Space{Name: "Suez Canal", Offset: 0, Turns: -1, Day: 13},
	Space{Name: "Thanksgiving", Offset: 4, Turns: 0, Day: 14},
	Space{Name: "Ismailia", Offset: 0, Turns: 0, Day: 15},
	Space{Name: "Red Sea", Offset: 0, Turns: 0, Day: 16},
	Space{Name: "Stormy", Offset: -5, Turns: 0, Day: 17},
	Space{Name: "Strikes a Rock", Offset: -5, Turns: 0, Day: 18},
	Space{Name: "Aden", Offset: +5, Turns: 0, Day: 19},
	Space{Name: "Arabian Sea", Offset: 0, Turns: 1, Day: 20},
	Space{Name: "Stuck on Sand Bar", Offset: 0, Turns: -2, Day: 21},
	Space{Name: "Indian Ocean", Offset: -1, Turns: 0, Day: 22},
	Space{Name: "Indian Ocean", Offset: 0, Turns: 0, Day: 23},
	Space{Name: "Indian Ocean", Offset: 0, Turns: 0, Day: 24},
	Space{Name: "Indian Ocean - Out of Coal", Offset: 0, Turns: -1, Day: 25},
	Space{Name: "Colombo", Offset: 0, Turns: 0, Day: 26},
	Space{Name: "Ceylon", Offset: 0, Turns: 0, Day: 27},
	Space{Name: "Bay of Bengal", Offset: 6, Turns: 0, Day: 28},
	Space{Name: "Bay of Bengal", Offset: 0, Turns: 0, Day: 29},
	Space{Name: "Malacca Straits - Pirate Ship", Offset: -3, Turns: 0, Day: 30},
	Space{Name: "Off Sumtra", Offset: 0, Turns: 0, Day: 31},
	Space{Name: "Malacca Straits", Offset: 1, Turns: 0, Day: 32},
	Space{Name: "Singapore", Offset: 0, Turns: 0, Day: 33},
	Space{Name: "Siam", Offset: 2, Turns: 0, Day: 34},
	Space{Name: "China Sea", Offset: 0, Turns: 0, Day: 35},
	Space{Name: "Simoon", Offset: -10, Turns: 0, Day: 36},
	Space{Name: "Borneo", Offset: 2, Turns: 0, Day: 37},
	Space{Name: "China Sea", Offset: 0, Turns: 0, Day: 38},
	Space{Name: "China Sea", Offset: -3, Turns: 0, Day: 39},
	Space{Name: "Hong Kong", Offset: 0, Turns: 0, Day: 40},
	Space{Name: "Christmas", Offset: 0, Turns: 1, Day: 41},
	Space{Name: "Joss China", Offset: 0, Turns: 0, Day: 42},
	Space{Name: "Canton", Offset: 0, Turns: -1, Day: 43},
	Space{Name: "Hong Kong", Offset: 0, Turns: 0, Day: 44},
	Space{Name: "China Sea", Offset: 3, Turns: 0, Day: 45},
	Space{Name: "China Sea", Offset: 0, Turns: 0, Day: 46},
	Space{Name: "Off Formosa", Offset: -5, Turns: 0, Day: 47},
	Space{Name: "New Year's Day", Offset: 5, Turns: 0, Day: 48},
	Space{Name: "Yokohama", Offset: 1, Turns: 0, Day: 49},
	Space{Name: "Yokohama", Offset: 0, Turns: 0, Day: 50},
	Space{Name: "Yeddo", Offset: 0, Turns: 1, Day: 51},
	Space{Name: "Yokohama", Offset: 0, Turns: 0, Day: 52},
	Space{Name: "Yokoham - Delay", Offset: -5, Turns: 0, Day: 53},
	Space{Name: "Yokohama", Offset: 0, Turns: 0, Day: 54},
	Space{Name: "On the Pacific", Offset: -2, Turns: 0, Day: 55},
	Space{Name: "Stormy", Offset: -10, Turns: 0, Day: 56},
	Space{Name: "Clear", Offset: 1, Turns: 0, Day: 57},
	Space{Name: "Break in Machinery", Offset: -3, Turns: 0, Day: 58},
	Space{Name: "Clear", Offset: 0, Turns: 0, Day: 59},
	Space{Name: "Fair", Offset: +2, Turns: 0, Day: 60},
	Space{Name: "Clear", Offset: 0, Turns: 0, Day: 61},
	Space{Name: "Storm", Offset: 0, Turns: -3, Day: 62},
	Space{Name: "Collision", Offset: -15, Turns: 0, Day: 63},
	Space{Name: "On Raft", Offset: 0, Turns: -2, Day: 64},
	Space{Name: "Rescued", Offset: 1, Turns: 0, Day: 65},
	Space{Name: "Clear", Offset: 0, Turns: 0, Day: 66},
	Space{Name: "Pacific Ocean", Offset: 0, Turns: -1, Day: 67},
	Space{Name: "Golden Gate", Offset: 0, Turns: 0, Day: 68},
	Space{Name: "Sierra Mountains - Snow Bound", Offset: 0, Turns: -5, Day: 69},
	Space{Name: "Cheyenne Indians", Offset: -2, Turns: 0, Day: 70},
	Space{Name: "Omaha", Offset: 3, Turns: 0, Day: 71},
	Space{Name: "Leaving Chicago", Offset: 0, Turns: -1, Day: 72},
	Space{Name: "First Part of Day 73", Offset: 0, Turns: 0, Day: 73},
	Space{Name: "Center", Offset: 0, Turns: 0, Day: 74},
}
