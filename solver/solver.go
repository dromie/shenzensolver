package solver

import (
	"fmt"
	"runtime"
	"sync"

	A "github.com/IBM/fp-go/array"
	"github.com/dromie/shenzensolver/util"
	"github.com/mohae/deepcopy"
)

type Place int

const (
	H1 Place = iota
	H2
	H3
	S1
	S2
	S3
	T1
	T2
	T3
	T4
	T5
	T6
	T7
	T8
	RBLOCK
	GBLOCK
	BBLOCK
	OBLOCK
)

func table_place(i int) Place {
	if 0 <= i && i < 8 {
		return Place(int(T1) + i)
	}
	panic("Invalid table place " + fmt.Sprintf("%d", i))
}

func hold_place(i int) Place {
	if 0 <= i && i < 3 {
		return Place(int(H1) + i)
	}
	panic("Invalid hold place " + fmt.Sprintf("%d", i))
}

func solved_place(i int) Place {
	if 0 <= i && i < 3 {
		return Place(int(S1) + i)
	}
	panic("Invalid solved place " + fmt.Sprintf("%d", i))
}

func block_place(color CardSuit) Place {
	return Place(int(RBLOCK) - 1 + int(color))
}

func (p Place) block_color() CardSuit {
	return CardSuit(int(p) - int(RBLOCK) + 1)
}

type Move struct {
	from  Place
	to    Place
	depth int
}

func get_valid_moves(table *Table) []Move {
	pq := util.Pqueue_init[Move]()
	for i, column := range table.Table {
		for j, hold := range table.Hold {
			if hold == (Card{}) {
				if len(column) > 0 && column[len(column)-1] == constructCard("O") {
					pq.Push(
						&util.Item[Move]{
							Value: Move{
								table_place(i),
								hold_place(j),
								0,
							},
							Priority: -50,
						},
					)
				}
			} else {
				if hold.Value == BLOCK_HOLD {
					continue
				}
				if len(column) == 0 || hold.can_be_put_over(column[len(column)-1]) {
					pq.Push(&util.Item[Move]{Value: Move{hold_place(j), table_place(i), 0}, Priority: -40})
				}
			}
		}
		if len(column) == 0 {
			continue
		}
		if column[len(column)-1] == constructCard("O") {
			return []Move{{
				table_place(i),
				OBLOCK,
				0}}
		}
		for j, solved := range table.Solved {
			if column[len(column)-1].is_solution(solved) {
				pq.Push(
					&util.Item[Move]{
						Value: Move{
							table_place(i),
							solved_place(j),
							0,
						},
						Priority: 5},
				)
				if column[len(column)-1].Value == 1 {
					break
				}
			}
		}
		for j, column2 := range table.Table {
			if len(column2) == 0 || column[len(column)-1].can_be_put_over(column2[len(column2)-1]) {
				pq.Push(&util.Item[Move]{
					Value: Move{
						table_place(i),
						table_place(j),
						0,
					},
					Priority: -10})
			}
			k := len(column) - 1
			for ; k > 0 && column[k].can_be_put_over(column[k-1]); k-- {
				if len(column2) > 0 && column[k-1].can_be_put_over(column2[len(column2)-1]) {
					pq.Push(&util.Item[Move]{
						Value: Move{
							table_place(i),
							table_place(j),
							len(column) - k,
						},
						Priority: -(11 - k)})
				}
			}
		}
	}
	for _, color := range COLORDICT {
		free_hold := A.Any(func(card Card) bool { return card == (Card{}) })(table.Hold)
		count := 0
		for _, column := range table.Table {
			if len(column) == 0 {
				continue
			}
			lastCard := column[len(column)-1]
			if lastCard.Suit == color && column[0].Value == BLOCK {
				count++
			}
		}
		for _, hold := range table.Hold {
			if hold != (Card{}) {
				if hold.Suit == color && hold.Value == BLOCK {
					free_hold = true
					count++
				}
			}
		}
		if count == 4 && free_hold {
			pq.Push(&util.Item[Move]{
				Value: Move{
					block_place(color),
					block_place(color),
					0,
				}, Priority: 2})
		}
	}
	for i, hold := range table.Hold {
		if hold != (Card{}) {
			for j, solved := range table.Solved {
				if hold.is_solution(solved) {
					pq.Push(&util.Item[Move]{
						Value: Move{
							hold_place(i),
							solved_place(j),
							0,
						},
						Priority: 5})
				}
			}
		}
	}
	moves := make([]Move, pq.Len())
	i := 0
	for pq.Len() > 0 {
		moves[i] = pq.Pop().Value
		i++
	}
	return moves
}

func make_block_move(table *Table, move Move) Table {
	new_table := deepcopy.Copy(*table).(Table)
	card := Card{move.from.block_color(), BLOCK}
	hold_index := -1
	for i, hold := range new_table.Hold {
		if hold == card {
			new_table.Hold[i] = (Card{})
			if hold_index == -1 {
				hold_index = i
			}
		}
		if hold == (Card{}) && hold_index == -1 {
			hold_index = i
		}
	}
	for i, column := range new_table.Table {
		if len(column) > 0 && column[len(column)-1] == card {
			new_table.Table[i] = column[:len(column)-1]
		}
	}
	if hold_index == -1 {
		panic("Invalid move")
	} else {
		new_table.Hold[hold_index] = card
		new_table.Hold[hold_index].Value = BLOCK_HOLD
	}
	return new_table
}

func make_move(table *Table, move Move) Table {
	if move.from == RBLOCK || move.from == GBLOCK || move.from == BBLOCK {
		return make_block_move(table, move)
	}
	if move.from == OBLOCK {
		new_table := deepcopy.Copy(*table).(Table)
		new_table.Table[move.from] = new_table.Table[move.from][:len(new_table.Table[move.from])-1]
		return new_table
	}
	new_table := deepcopy.Copy(*table).(Table)
	if move.depth == 0 {
		card := new_table.PopCard(move.from)
		new_table.PushCard(move.to, card)
	} else {
		column := new_table.Table[move.from-T1]
		cards := column[len(column)-move.depth-1:]
		new_table.Table[move.from-T1] = column[:len(column)-move.depth-1]
		new_table.Table[move.to-T1] = append(new_table.Table[move.to-T1], cards...)
	}
	return new_table
}

type State struct {
	table     *Table
	prevState *State
	move      Move
	depth     int
}

func solve_helper(in chan<- interface{}, out <-chan interface{}, solutions chan State, allStates *sync.Map) {
	for item := range out {
		state := item.(State)
		for _, move := range get_valid_moves(state.table) {
			new_table := make_move(state.table, move)
			new_state := State{&new_table, &state, move, state.depth + 1}
			if new_table.is_solved() {
				solutions <- new_state
				return
			}
			if oo, found := allStates.Load(new_state.table.String()); !found {
				allStates.Store(new_state.table.String(), new_state)
				in <- new_state
			} else {
				oldState := oo.(State)
				if oldState.depth > state.depth+1 {
					allStates.Store(state.table.String(), state)
				}
			}
		}
	}
}

func solve(table *Table) []Move {
	states := sync.Map{}
	in, out := util.MakeInfinite()
	solutions := make(chan State)
	init_state := State{table, nil, Move{}, 0}
	states.Store(init_state.table.String(), init_state)
	in <- init_state
	for i := 0; i < runtime.NumCPU(); i++ {
		go solve_helper(in, out, solutions, &states)
	}

	solution := <-solutions
	moves := []Move{}
	for solution.prevState != nil {
		fmt.Println(solution.move)
		moves = append(moves, solution.move)
		solution = *solution.prevState
		if solution2, found := states.Load(solution.table.String()); found {
			solution = solution2.(State)
		} else {
			panic("Failed to find modified state in state map")
		}
	}
	return moves
}
