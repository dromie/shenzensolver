package solver

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mohae/deepcopy"
)

func Test_RBLOCK(t *testing.T) {
	test_table := Table{}
	test_table.init()
	test_table.load_table([]string{"R", "R", "R", "R"})
	moves := get_valid_moves(&test_table)
	fmt.Println(moves)
	if RBLOCK != moves[0].from || RBLOCK != moves[0].to {
		t.Errorf("Invalid move")
	}
}

func Test_straight_move(t *testing.T) {
	state1 := Table{}
	state1.init()
	state1.Hold = []Card{constructCard("G"), constructCard("G"), constructCard("G")}
	state1.load_table([]string{"r1 g4 b3 g2", "b5"})
	expected := Table{}
	expected.Hold = deepcopy.Copy(state1.Hold).([]Card)
	expected.load_table([]string{"r1", "b5 g4 b3 g2"})
	fmt.Println(state1.String())
	moves := get_valid_moves(&state1)
	fmt.Println(len(moves))
	fmt.Println(moves[0])
	if len(moves) < 1 {
		t.Errorf("No moves found")
	}
	state11 := make_move(&state1, moves[0])
	if !reflect.DeepEqual(state11, expected) {
		t.Errorf("Not expected state")
	}

}

func Test_block_color(t *testing.T) {
	redBlockPlace := RBLOCK
	if redBlockPlace.block_color() != RED {
		t.Errorf("Invalid block color")
	}
	blackBlockPlace := BBLOCK
	if blackBlockPlace.block_color() != BLACK {
		t.Errorf("Invalid block color")
	}
	greenBlockPlace := GBLOCK
	if greenBlockPlace.block_color() != GREEN {
		t.Errorf("Invalid block color")
	}
}

func Test_block_move(t *testing.T) {
	state1 := Table{}
	state1.init()
	state1.Hold = []Card{constructCard("G"), constructCard("G"), constructCard("G")}
	state1.load_table([]string{"r1 g3 b3 g2", "b5", "G"})
	expected := Table{}
	expected.init()
	expected.Hold[0] = Card{GREEN, BLOCK_HOLD}
	expected.load_table([]string{"r1 g3 b3 g2", "b5"})
	moves := get_valid_moves(&state1)
	if len(moves) == 0 {
		t.Errorf("No moves found")
	}
	state11 := make_move(&state1, moves[0])
	if !reflect.DeepEqual(state11, expected) {
		t.Errorf("Not expected state: \n%svs \n%s", state11.String(), expected.String())
	}

}

func Test_OBLOCK(t *testing.T) {
	test_table := Table{}
	test_table.init()
	test_table.load_table([]string{"O"})
	moves := get_valid_moves(&test_table)
	fmt.Println(moves)
	moves = get_valid_moves(&test_table)
	if len(moves) != 1 {
		t.Errorf("Invalid move OBLOCK not found")
	}
	if moves[0].from != T1 || moves[0].to != OBLOCK {
		t.Errorf("Invalid move, not an OBLOCK move")

	}
}

func Test_one_move_to_win(t *testing.T) {
	test_table := Table{}
	test_table.init()
	test_table.Solved = []Card{{RED, EIGHT}, {GREEN, NINE}, {BLACK, NINE}}
	test_table.load_table([]string{"r9"})
	moves := solve(&test_table)
	if len(moves) != 1 {
		t.Errorf("Invalid solution length %d", len(moves))
	}
	fmt.Println(moves)
}

func Test_more_move_to_win(t *testing.T) {
	test_table := Table{}
	test_table.init()
	test_table.Solved = []Card{{}, {}, {BLACK, NINE}}
	test_table.load_table([]string{"r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1"})
	moves := solve(&test_table)
	if len(moves) != 2 {
		t.Errorf("Invalid solution length %d", len(moves))
	}
	fmt.Println(moves)
}

func Test_solve_trivial(t *testing.T) {
	table := Table{}
	table.init()
	table.load_table([]string{"r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1", "b9 b8 b7 b6 b5 b4 b3 b2 b1"})
	moves := solve(&table)
	if len(moves) != 27 {
		t.Errorf("Invalid solution length %d", len(moves))
	}
}
