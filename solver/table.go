package solver

import (
	"fmt"
	"strings"

	A "github.com/IBM/fp-go/array"
	. "github.com/dromie/shenzensolver/util"
)

type Table struct {
	Hold   []Card
	Solved []Card
	Table  [][]Card
}

func (t *Table) init() {
	if t.Table == nil {
		t.Table = make([][]Card, 8)
		for i := range 8 {
			t.Table[i] = make([]Card, 0)
		}
	}
	if t.Hold == nil {
		t.Hold = []Card{{}, {}, {}}
	}
	if t.Solved == nil {
		t.Solved = []Card{{}, {}, {}}
	}
}

func (t *Table) load_table(tablestr []string) {
	t.init()
	for i, line := range tablestr {
		t.Table[i] = []Card{}
		for _, card := range strings.Split(line, " ") {
			t.Table[i] = append(t.Table[i], constructCard(card))
		}
	}
}

func rowStr(row []Card) string {
	var rowStr []string
	for _, card := range row {
		rowStr = append(rowStr, card.String())
	}
	return strings.Join(rowStr, ", ")
}

func (t *Table) String() string {
	var str string
	str += "Hold: " + rowStr(t.Hold) + "\n"
	str += "Solved: " + rowStr(t.Solved) + "\n"
	str += "Table:\n"
	for _, row := range t.Table {
		str += rowStr(row)
		str += "\n"
	}
	return str
}

func (t *Table) is_solved() bool {
	result := len(t.Solved) == 3 && !A.Any(func(card Card) bool { return card.Value != 9 })(t.Solved)
	return result
}

func (t *Table) heuristic() int {
	h := 0
	h += Sum()(A.Map(func(card Card) int { return int(card.Value) * 6 })(t.Solved))
	h -= Count[[]Card](func(card Card) bool { return card != Card{} && card.Value != BLOCK_HOLD })(t.Hold)
	h += Count[[][]Card](func(row []Card) bool { return len(row) == 0 })(t.Table)
	h += Count[[][]Card](func(row []Card) bool { return len(row) > 0 && row[0].Value == 9 })(t.Table)
	for _, color := range COLORDICT {
		b := Count[[][]Card](func(row []Card) bool {
			return len(row) > 0 && row[len(row)-1].Value == BLOCK && row[len(row)-1].Suit == color
		})(t.Table)
		if b == 4 {
			h += 4
		}
	}
	return h
}

func (t *Table) PopCard(place Places) Card {
	if place >= H1 && place <= H3 {
		card := t.Hold[place-H1]
		t.Hold[place-H1] = Card{}
		return card
	}
	if place >= T1 && place <= T8 {
		card := t.Table[place-T1][len(t.Table[place-T1])-1]
		t.Table[place-T1] = t.Table[place-T1][:len(t.Table[place-T1])-1]
		return card
	}
	panic(fmt.Sprintf("GetCard: invalid place %v", place))
}

func (t *Table) PushCard(place Places, card Card) {
	if place >= H1 && place <= H3 {
		t.Hold[place-H1] = card
		return
	}
	if place >= S1 && place <= S3 {
		t.Solved[place-S1] = card
		return
	}
	if place >= T1 && place <= T8 {
		t.Table[place-T1] = append(t.Table[place-T1], card)
		return
	}
	panic(fmt.Sprintf("PushCard: invalid place %v", place))
}
