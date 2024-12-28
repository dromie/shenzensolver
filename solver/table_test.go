package solver

import (
	"reflect"
	"testing"

	"github.com/mohae/deepcopy"
)

func Test_loadtable(t *testing.T) {
	table := Table{}
	table.load_table([]string{"r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1", "b9 b8 b7 b6 b5 b4 b3 b2 b1"})
	if len(table.Table) != 8 {
		t.Errorf("Table length is not 3")
	}
	if len(table.Table[0]) != 9 {
		t.Errorf("Table row length is not 9 it is %d", len(table.Table[0]))
	}
	if table.Table[0][0].Suit != RED {
		t.Errorf("Card suit is not RED")
	}
}

func Test_is_solved(t *testing.T) {
	table := Table{}
	table.load_table([]string{"r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1", "b9 b8 b7 b6 b5 b4 b3 b2 b1"})
	if table.is_solved() {
		t.Errorf("Table is not solved")
	}
	table.Solved = []Card{{RED, NINE}, {GREEN, NINE}, {BLACK, NINE}}
	table.Table = [][]Card{}
	if !table.is_solved() {
		t.Errorf("Table is not solved")
	}
}

func Test_heuristic(t *testing.T) {
	table := Table{}
	table.load_table([]string{"r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1", "b9 b8 b7 b6 b5 b4 b3 b2 b1"})
	if table.heuristic() != 8 { // All rows stars with 9
		t.Errorf("Heuristic is not 3 it is %d", table.heuristic())
	}
	table.Hold = []Card{{RED, BLOCK}, {GREEN, BLOCK}, {BLACK, BLOCK}}
	table.Table = nil
	if table.heuristic() != -3 {
		t.Errorf("Heuristic is not -3")
	}
	table.Hold = nil
	table.Table = [][]Card{{}, {}, {}, {}}
	if table.heuristic() != 4 {
		t.Errorf("Heuristic is not 4")
	}
	table.Table = nil
	table.Hold = nil
	table.Solved = []Card{{RED, NINE}, {GREEN, NINE}, {BLACK, NINE}}
	if table.heuristic() != 9*6*3 {
		t.Errorf("Heuristic is not 27")
	}

}

func Test_Table_DeepCopy(t *testing.T) {
	table := Table{}
	table.init()
	table.load_table([]string{"r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1", "b9 b8 b7 b6 b5 b4 b3 b2 b1"})
	copy := deepcopy.Copy(table).(Table)
	if !reflect.DeepEqual(table, copy) {
		t.Errorf("Deep copy failed")
	}
}
