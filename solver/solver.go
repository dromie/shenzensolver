package solver

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	"github.com/dromie/shenzensolver/util"
	"github.com/mohae/deepcopy"
)

type Places int

const (
	H1 Places = iota
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

func table_place(i int) Places {
	if 0 <= i && i < 8 {
		return Places(int(T1) + i)
	}
	panic("Invalid table place " + fmt.Sprintf("%d", i))
}

func hold_place(i int) Places {
	if 0 <= i && i < 3 {
		return Places(int(H1) + i)
	}
	panic("Invalid hold place " + fmt.Sprintf("%d", i))
}

func solved_place(i int) Places {
	if 0 <= i && i < 3 {
		return Places(int(S1) + i)
	}
	panic("Invalid solved place " + fmt.Sprintf("%d", i))
}

func block_place(color CardSuit) Places {
	return Places(int(RBLOCK) - 1 + int(color))
}

func (p Places) block_color() CardSuit {
	return CardSuit(int(p) - int(RBLOCK) + 1)
}

type CardCoordinate struct {
	place Places
	index int
}

type Move struct {
	from  CardCoordinate
	to    CardCoordinate
	depth int
}

func get_valid_moves(table *Table) []Move {
	pq := util.Pqueue_init[Move]()
	for i, column := range table.Table {
		column_height := len(column)
		for j, hold := range table.Hold {
			if hold == (Card{}) {
				if len(column) > 0 && column[len(column)-1] == constructCard("O") {
					pq.Push(
						&util.Item[Move]{
							Value: Move{
								CardCoordinate{table_place(i), 0},
								CardCoordinate{hold_place(j), 0},
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
					pq.Push(&util.Item[Move]{Value: Move{CardCoordinate{hold_place(j), 0}, CardCoordinate{table_place(i), 0}, 0}, Priority: -40})
				}
			}
		}
		if len(column) == 0 {
			continue
		}
		if column[len(column)-1] == constructCard("O") {
			return []Move{{
				CardCoordinate{table_place(i), column_height},
				CardCoordinate{OBLOCK, 0},
				0}}
		}
		for j, solved := range table.Solved {
			if column[len(column)-1].is_solution(solved) {
				pq.Push(
					&util.Item[Move]{
						Value: Move{
							CardCoordinate{table_place(i), column_height},
							CardCoordinate{solved_place(j), 0},
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
						CardCoordinate{table_place(i), column_height},
						CardCoordinate{table_place(j), len(column2)},
						0,
					},
					Priority: -10})
			}
			k := len(column) - 1
			for ; k > 0 && column[k].can_be_put_over(column[k-1]); k-- {
				if len(column2) > 0 && column[k-1].can_be_put_over(column2[len(column2)-1]) {
					pq.Push(&util.Item[Move]{
						Value: Move{
							CardCoordinate{table_place(i), k},
							CardCoordinate{table_place(j), len(column2)},
							k,
						},
						Priority: -(11 + k)})
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
					CardCoordinate{block_place(color), 0},
					CardCoordinate{block_place(color), 0},
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
							CardCoordinate{hold_place(i), 0},
							CardCoordinate{solved_place(j), 0},
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
	card := Card{move.from.place.block_color(), BLOCK}
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
	if move.from.place == RBLOCK || move.from.place == GBLOCK || move.from.place == BBLOCK {
		return make_block_move(table, move)
	}
	if move.from.place == OBLOCK {
		new_table := deepcopy.Copy(*table).(Table)
		new_table.Table[move.from.place] = new_table.Table[move.from.place][:len(new_table.Table[move.from.place])-1]
		return new_table
	}
	new_table := deepcopy.Copy(*table).(Table)
	if move.depth == 0 {
		card := new_table.PopCard(move.from.place)
		new_table.PushCard(move.to.place, card)
	} else {
		column := new_table.Table[move.from.place]
		cards := column[len(column)-move.depth:]
		new_table.Table[move.from.place] = column[:len(column)-move.depth]
		new_table.Table[move.to.place] = append(new_table.Table[move.to.place], cards...)
	}
	return new_table
}
