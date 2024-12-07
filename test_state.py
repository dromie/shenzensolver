#!/usr/bin/env pytest
import pytest
from copy import deepcopy

from main import Card, State, Places, Move, solve, BLOCK_CARD_HOLD, BLOCK_CARD

def test_OBLOCK():
    card = Card.construct("O")
    test_state = State()
    test_state.table = [[card]]
    moves = test_state.get_valid_moves()
    assert len(moves) == 1
    assert moves[0].to == Places.OBLOCK
    assert moves[0].from_ == Places.T1

def test_RBLOCK():
    card = Card.construct("R")
    test_state = State()
    test_state.table = [[card], [card], [card], [card]]
    moves = test_state.get_valid_moves()
    print(moves)
    assert Move(Places.RBLOCK,Places.RBLOCK) in moves

def test_solve_trivial():
    init_state = State()
    init_state.load_table(["r9 r8 r7 r6 r5 r4 r3 r2 r1", "g9 g8 g7 g6 g5 g4 g3 g2 g1", "b9 b8 b7 b6 b5 b4 b3 b2 b1"])
    moves = solve(init_state)
    assert len(moves) == 27

def test_solve_trivial_with_blockcards():
    init_state = State()
    init_state.load_table(["r9 r8 r7 r6 r5 r4 r3 r2 r1 R G B", "g9 g8 g7 g6 g5 g4 g3 g2 g1 R G B", "b9 b8 b7 b6 b5 b4 b3 b2 b1 R G B", "R", "G", "B"])
    moves = solve(init_state)
    print(moves)
    assert len(moves) < 50

def test_no_solution():
    init_state = State()
    init_state.load_table(["r9 r8 r7 r6 r5 r4 r3 r2 r1 R R R", "g9 g8 g7 g6 g5 g4 g3 g2 g1 G G G", "b9 b8 b7 b6 b5 b4 b3 b2 b1 B B B"])
    moves = solve(init_state)
    assert len(moves) == 0

def test_circle_detection():
    init_state1 = State()
    init_state1.load_table(["r9 r8 r7 r6 r5 r4 r3 r2 r1 R G B", "g9 g8 g7 g6 g5 g4 g3 g2 g1 R G B", "b9 b8 b7 b6 b5 b4 b3 b2 b1 R G B", "R", "G", "B"])
    init_state2 = deepcopy(init_state1)
    assert init_state1 == init_state2
    test_set1 = set()
    test_set1.add(init_state1)
    assert init_state1 in test_set1
    assert init_state2 in test_set1
    assert len(test_set1) == 1
    test_set1.add(init_state2)
    assert len(test_set1) == 1


def test_circle2():
    state1 = State()
    state1.solved = [None, Card.construct("b1"), None]
    state1.load_table(["r1 b2 G G r7 g6", "g5 B r8 B R", "r2 b8 g1 r4 g3", "r9 G g8 r3", "b7 R r6", "b4 b5 G g7 g9", "R B B b9 b3 g2", "r5 g4 b6 R"])
    test_set = set([state1])
    move = state1.get_valid_moves()[0]
    print("---BEGIN---")
    print(state1)
    print(move)
    state11 = state1.move(move)
    assert state11 not in test_set
    test_set.add(state11)
    move = state11.get_valid_moves()[0]
    print(move)
    state12 = state11.move(move)
    print("---END---")
    print(state12)
    assert str(state12) == str(state1)
    assert state12.__hash__() == state1.__hash__()
    assert state12 in test_set


def test_straight_move():
    state1 = State()
    state1.hold = [Card.construct("G"), Card.construct("G"), Card.construct("G")]
    state1.load_table(["r1 g4 b3 g2","b5"])
    expected = State()
    expected.hold = state1.hold
    expected.load_table(["r1", "b5 g4 b3 g2"])
    print(state1)
    moves = state1.get_valid_moves()
    print(moves[0])
    state11 = state1.move(moves[0])
    print(state11)
    assert state11 == expected

def test_why():
#b-10 b9 g-10
#b8 g9 r9
#1:
#2:
#3: r-5
#4: r-5 r-5
#5:
#6: r-5
#7:
#8:
    state1 = State()
    state1.hold = [Card(2,BLOCK_CARD_HOLD), Card.construct("b9"), Card(1,BLOCK_CARD_HOLD)]
    state1.solved = [Card.construct("b8"), Card.construct("g9"), Card.construct("r9")]
    state1.load_table(["", "", "R", "R R", "", "R", "", ""])
    print(state1)
    moves = state1.get_valid_moves()
    print(moves)
    assert any([move.from_ == Places.H2 for move in moves])


