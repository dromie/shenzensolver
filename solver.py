#!/usr/bin/env python3
import sys
from dataclasses import dataclass, field
import os
import time
import queue
from typing import List, Tuple, Set, Any, Callable
from enum import Enum
from copy import deepcopy

BLOCK_CARD = -5
BLOCK_CARD_HOLD = -10

COLORDICT ={
            "r": 0,
            "g": 1,
            "b": 2,
            "o": 3
        }
REVCOLORDICT = {v:k for k,v in COLORDICT.items()} 

@dataclass
class Card:
    color: int
    number: int
    @classmethod
    def construct(cls, card: str) -> "Card":
        if card == "":
            return None
        color = COLORDICT[card[0].lower()]
        if len(card) == 2:
            number = int(card[1])
        else:
            number = BLOCK_CARD
        return cls(color, number)

    def __init__(self, color: int, number: int):
        self.color = color
        self.number = number

    def is_solution(self, other: "Card") -> bool:
        if self.number == BLOCK_CARD:
            return False
        if other is None:
            return self.number == 1
        return self.color == other.color and self.number == other.number +1

    def can_be_put_over(self, other: "Card") -> bool:
        if other is None:
            return True
        if self.number == BLOCK_CARD:
            return False
        if self.color != other.color and self.number == other.number-1:
            return True
        return False

    def can_be_put_over_column(self, col: List["Card"]) -> bool:
        if len(col) == 0:
            return True
        return self.can_be_put_over(col[-1])

    def __str__(self):
        return f"{REVCOLORDICT[self.color]}{self.number}"
    
    def __repr__(self):
        return self.__str__()

@dataclass(order=True)
class SinglePrioritizedItem:
    priority: int
    item: Any=field(compare=False)

@dataclass(order=True)
class DualPrioritizedItem:
    priority1: int
    priority2: int
    item: Any=field(compare=False)

class Places(Enum):
    H1 = 0
    H2 = 1
    H3 = 2
    S1 = 3
    S2 = 4
    S3 = 5
    T1 = 6
    T2 = 7
    T3 = 8
    T4 = 9
    T5 = 10
    T6 = 11
    T7 = 12
    T8 = 13
    RBLOCK = 14
    GBLOCK = 15
    BBLOCK = 16
    OBLOCK = 17

    @classmethod
    def table_place(cls, i: int) -> "Places":
        return cls(i+6)

    @classmethod
    def hold_place(cls, i: int) -> "Places":
        return cls(i)
    
    @classmethod
    def solved_place(cls, i: int) -> "Places":
        return cls(i+3)
    
    @classmethod
    def block_place(cls, color: int) -> "Places":
        return cls(cls.RBLOCK.value + color)
    
    def color(self) -> int:
        return self.value - Places.RBLOCK.value;

@dataclass
class CardCoordinates:
    place: Places
    position: int = 0

@dataclass
class Move:
    from_: CardCoordinates
    to: CardCoordinates
    depth: int

    def __init__(self, from_: CardCoordinates, to: CardCoordinates, depth = 0):
        self.from_ = from_
        self.to = to
        self.depth = depth

    def __str__(self):
        return f"{self.from_}(depth: {self.depth}) -> {self.to}"

    def __repr__(self):
        return self.__str__()

@dataclass
class State:
    hold: List[Card]
    solved: List[Card]
    table: List[List[Card]]
    prev_state: Tuple["State", Move, int]

    def __init__(self):
        self.hold = [None, None, None]
        self.solved = [None, None, None]
        self.table = []
        self.prev_state = (None, None, 0)

    def load_table(self, table:List[str]):
        self.table = []
        for line in table:
            self.table.append(list(filter(lambda x: x is not None, map(lambda x: Card.construct(x), line.strip().split(" ")))))

    def __hash__(self):
        return hash(self.__repr__())

    def __eq__(self, other: "State"):
        return str(self) == str(other)

    def __str__(self):
        #hold = " ".join(sorted(map(lambda x: str(x), self.hold)))
        #solved = " ".join(sorted(map(lambda x: str(x), self.solved)))
        hold = " ".join(map(lambda x: str(x), self.hold))
        solved = " ".join(map(lambda x: str(x), self.solved))
        table = "\n".join([" ".join([str(i+1)+":"] + list(map(lambda x: str(x), self.table[i]))) for i in range(len(self.table))])
        return f"{hold}\n{solved}\n{table}"
    
    def __repr__(self) -> str:
        hold = " ".join(map(lambda x: str(x), self.hold))
        solved = " ".join(map(lambda x: str(x), self.solved))
        table = "\n".join(reversed(list(map(lambda x: " ".join(map(lambda y: str(y), x)), self.table))))
        return f"{hold}\n{solved}\n{table}"

    def is_solved(self) -> bool:
        return all(map(lambda x: x is not None and x.number == 9, self.solved))

    def heuristic(self) -> int:
        h = 0
        h += sum(map(lambda x: 0 if x is None else x.number, self.solved)) * 6 
        h -= sum(map(lambda x: 0 if x is None else 1, self.hold)) 
        h += sum(map(lambda x: 1 if len(x) == 0 else 0, self.table)) 
        h += sum(map(lambda x: 1 if x is not None and x.number == BLOCK_CARD_HOLD else 0, self.hold))
        h += sum(map(lambda x: 1 if len(x) > 0 and x[0].number == 9 else 0, self.table))
        for color in ["r", "g", "b"]:
            b = sum(map(lambda x: 1 if len(x) > 0 and x[-1].number == BLOCK_CARD and x[-1].color==COLORDICT[color] else 0, self.table))
            if b==4:
                h+=4
        return h

    
    def get_valid_moves(self) -> List[Move]:
        moves = queue.PriorityQueue()
        for i in range(len(self.table)):
            column_height = len(self.table[i])
            for j in range(len(self.hold)):
                if self.hold[j] is None:
                    if len(self.table[i]) > 0 and not (self.table[i][-1].number == BLOCK_CARD and self.table[i][-1].color == COLORDICT["o"]):
                        moves.put(SinglePrioritizedItem(50, Move(CardCoordinates(Places.table_place(i), column_height), CardCoordinates(Places.hold_place(j)) )))
                else:
                    if self.hold[j].number == BLOCK_CARD_HOLD:
                        continue
                    if len(self.table[i]) == 0 or self.hold[j].can_be_put_over(self.table[i][-1]):
                        moves.put(SinglePrioritizedItem(40, Move(CardCoordinates(Places.hold_place(j)), CardCoordinates(Places.table_place(i), column_height))))
            if len(self.table[i]) == 0:
                continue
            if self.table[i][-1].number == BLOCK_CARD and self.table[i][-1].color == COLORDICT["o"]:
                moves.put(SinglePrioritizedItem(-9, Move(CardCoordinates(Places.table_place(i), column_height), CardCoordinates(Places.OBLOCK))))
                continue
            for j in range(len(self.solved)):
                    if self.table[i][-1].is_solution(self.solved[j]):
                        moves.put(SinglePrioritizedItem(-5, Move(CardCoordinates(Places.table_place(i), column_height), CardCoordinates(Places.solved_place(j)))))
                        if self.table[i][-1].number == 1:
                            break
            for j in range(len(self.table)):
                if len(self.table[j])==0 or self.table[i][-1].can_be_put_over(self.table[j][-1]):
                    moves.put(SinglePrioritizedItem(10, Move(CardCoordinates(Places.table_place(i), column_height), CardCoordinates(Places.table_place(j), len(self.table[j])))))
                k = -1
                while k > -len(self.table[i]) and self.table[i][k].can_be_put_over(self.table[i][k-1]):
                    if len(self.table[j])>0 and self.table[i][k-1].can_be_put_over(self.table[j][-1]):
                        moves.put(SinglePrioritizedItem(11+k, Move(CardCoordinates(Places.table_place(i), column_height+k), CardCoordinates(Places.table_place(j), len(self.table[j])), -k)))
                    k -= 1
        for color in range(3):
            free_hold = sum(map(lambda x: x is None, self.hold)) > 0
            count = 0
            for i in range(len(self.table)):
                if len(self.table[i]) == 0:
                    continue
                if self.table[i][-1].color == color and self.table[i][-1].number == BLOCK_CARD:
                    count += 1
            for i in range(len(self.hold)):
                if self.hold[i] is not None:
                    if self.hold[i].color == color and self.hold[i].number == BLOCK_CARD:
                        free_hold = True
                        count += 1
            if count == 4 and free_hold:
                moves.put(SinglePrioritizedItem(-2, Move(CardCoordinates(Places.block_place(color)), CardCoordinates(Places.block_place(color)))))
        for i in range(len(self.hold)):
            if self.hold[i] is None:
                continue
            for j in range(len(self.solved)):
                if self.hold[i].is_solution(self.solved[j]):
                    moves.put(SinglePrioritizedItem(-5, Move(CardCoordinates(Places.hold_place(i)), CardCoordinates(Places.solved_place(j)))))
                
        return [x.item for x in moves.queue]

    def move(self, move: Move) -> "State":
        new_state = State()
        new_state.hold = deepcopy(self.hold)
        new_state.solved = deepcopy(self.solved)
        new_state.table = deepcopy(self.table)
        new_state.prev_state = (self, move, self.prev_state[2]+1)
        if move.from_.place in [Places.RBLOCK, Places.GBLOCK, Places.BBLOCK]:
            card = Card(move.from_.place.color(), BLOCK_CARD)
            hold_index = None
            for i in range(len(new_state.hold)):
                if new_state.hold[i] == card:
                    new_state.hold[i] = None
                    if hold_index is None:
                        hold_index = i
                if new_state.hold[i] is None:
                    hold_index = i
            for i in range(len(new_state.table)):
                if len(new_state.table[i])>0 and new_state.table[i][-1] == card:
                    new_state.table[i].pop()
            if hold_index is not None:
                new_state.hold[hold_index] = card
                new_state.hold[hold_index].number = BLOCK_CARD_HOLD
            return new_state
        if move.to.place==Places.OBLOCK:
            new_state.table[move.from_.place.value-6].pop()
            return new_state
        elif move.from_.place in [Places.H1, Places.H2, Places.H3]:
            card = new_state.hold[move.from_.place.value]
            new_state.hold[move.from_.place.value] = None
        elif move.from_.place in [Places.S1, Places.S2, Places.S3]:
            card = new_state.solved[move.from_.place.value-3]
            new_state.solved[move.from_.place.value-3] = None
        else:
            if move.depth == 0:
                card = new_state.table[move.from_.place.value-6].pop()
            else:
                cards = new_state.table[move.from_.place.value-6][-move.depth-1:]
                new_state.table[move.from_.place.value-6] = new_state.table[move.from_.place.value-6][:-move.depth-1]
        if move.to.place in [Places.H1, Places.H2, Places.H3]:
            new_state.hold[move.to.place.value] = card
        elif move.to.place in [Places.S1, Places.S2, Places.S3]:
            new_state.solved[move.to.place.value-3] = card
        else:
            if move.depth == 0:
                new_state.table[move.to.place.value-6].append(card)
            else:
                new_state.table[move.to.place.value-6].extend(cards)
        return new_state

def get_moves(state: State) -> List[Move]:
    moves = []
    while state.prev_state[1] is not None:
        moves.append(state.prev_state[1])
        state = state.prev_state[0]
    moves.reverse()
    return moves

def a_star(P: queue.PriorityQueue, Pset: Set[State], Qset: Set[State], heuristic: Callable[[State], int]) -> List[Move]:
    size = [1, 1]
    score = [0, 0]
    while not P.empty():
        state = P.get().item
        Qset.add(state)
        Pset.remove(state)
        if heuristic(state) > score[0]:
            score[0] = heuristic(state)
            print("New best score: ", score[0])
            print(state)
            print(get_moves(state))
        if P.qsize() + len(Qset) > size[0]:
            size[0] = P.qsize() + len(Qset)
            if (size[0] - size[1]) > 10000:
                size[1] = size[0]
                print(f"{time.strftime("%H:%M:%S")}  max_size: {size[0]} = {len(Pset)} + {len(Qset)} best score: {heuristic(state)}, depth: {state.prev_state[2]}")
                print(state)
        for move in state.get_valid_moves():
            new_state = state.move(move)
            if new_state.is_solved():
                return get_moves(new_state)
            if new_state not in Pset and new_state not in Qset:
                P.put(DualPrioritizedItem(-heuristic(new_state), new_state.prev_state[2], new_state))
                Pset.add(new_state)
    return []

def solve(init_state: State) -> List[Move]:
    P = queue.PriorityQueue()
    P.put(DualPrioritizedItem(0, 0, init_state))
    Pset = set([init_state])
    Qset = set()
    moves = a_star(P, Pset, Qset, lambda x: x.heuristic())
#    P2 = queue.PriorityQueue()
#    for state in Pset:
#        if (state.prev_state[2]<len(moves)):
#            P2.put(PrioritizedItem(-state.prev_state[2], state))
#    moves2 = a_star(P2, Pset, Qset, lambda x: x.prev_state[2])
#    print("Found 1st solution, moving to find faster one.")
#    if len(moves2) < len(moves):
#        return moves2
    return moves

def main():
    init_state = State()
    init_state.load_table(sys.stdin.readlines())
    moves = solve(init_state)
    for move in moves:
        print(move)

if __name__ == "__main__":
    main()


