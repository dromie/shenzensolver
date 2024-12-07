#!/usr/bin/env python3
import cv2
import numpy as np
import json
from main import State, Card, solve

def analyze(screen):
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
    holdrow = None
    if len(rows)>5:
        holdrow = rows[0]
        rows = rows[1:]

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
    i = 0
    if holdrow is not None:
        for j,y in enumerate(columns):
            c = find_card(y, holdrow)
            if c is not None and c != 'O':
                state.solved[i] = Card.construct(c)
                i+=1
                
                
    state.load_table(table_str)
    print(state)
    return state


if __name__ == "__main__":
    screen = cv2.imread("test/test2.png",1)
    state = analyze(screen)
    #print("\n".join(map(str,solve(state))))

