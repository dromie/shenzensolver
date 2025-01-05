#!/usr/bin/env python3
from typing import List, Tuple, Set, Any

import cv2
import sys
import numpy as np
import json
from solver import State, Card, solve

def analyze(screen) -> Tuple[State, Tuple[List[int], List[int], int | None]]:
    mask = {}
    c='o'
    mask[f"{c.upper()}"] = cv2.imread(f"masks/{c}BLOCK.png",1)
    for c in ["r","g","b"]:
        mask[f"{c.upper()}"] = cv2.imread(f"masks/{c}BLOCK.png",1)
        for i in range(9):
            mask[f"{c}{i+1}"]=cv2.imread(f"masks/{c}{i+1}.png",1)

    cardmap = {}
    c, w, h = mask["O"].shape[::-1]
    columns = set()
    rows = set()
    for k, m in mask.items():
        result = cv2.matchTemplate(screen, m, cv2.TM_CCOEFF_NORMED)
        threshold = 0.999
        loc = np.where( result >= threshold)
        print(f"{k}: {len(loc[0])}")
        cardmap[k] = list(map(lambda x: (int(x[0]),int(x[1])),zip(*loc[::-1])))
        for a in cardmap[k]:
            columns.add(a[0])
            rows.add(a[1])


    rows = sorted(list(rows))
    columns = sorted(list(columns))

    def find_card(x,y):
        for k,v in cardmap.items():
            for a in v:
                if a[0]==x and a[1]==y:
                    return k
        return None
    holdrows = []
    while len(rows)>5:
        if len(holdrows)==0 or abs(holdrows[-1] - rows[0]) <= 5:
            holdrows.append(rows[0])
            rows = rows[1:]
        else:
            print(f"holdrow: {holdrows}, rows: {rows}")
            break

    table = [[None for _ in range(len(rows))] for _ in range(len(columns))]
    for i,x in enumerate(rows):
        for j,y in enumerate(columns):
            table[j][i] = find_card(y,x)
    print(table)
    table_str = []
    for i in range(len(table)):
        if not all([x is None for x in table[i]]):
            table_str.append((" ".join([x if x is not None else " " for x in table[i]])).strip())
    state = State()
    for holdrow in holdrows:
        for j,y in enumerate(columns):
            c = find_card(y, holdrow)
            if c is not None and c != 'O':
                print(f"Hold: {c}, {j}, {y}, {holdrows}")
                if j in [0,1,2]:
                    state.hold[j] = Card.construct(c)
                elif j in [5,6,7]:
                    state.solved[j-5] = Card.construct(c)
    state.load_table(table_str)
    if len(columns) == 9:
        columns = columns[:5] + columns[6:]
    return (state, (rows, columns, holdrows))


if __name__ == "__main__":
    file = "test/test2.png"
    if sys.argv[1:]:
        file = sys.argv[1]
    screen = cv2.imread(file,1)
    analisis = analyze(screen)
    print(analisis[0])

