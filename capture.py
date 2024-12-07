#!/usr/bin/env python3
import cv2
import numpy as np

screen = cv2.imread("test/test1.png")
card = cv2.imread("test/card.png")
card_edge_ignore_color = np.array([255, 0, 255])
card_value_ignore_color = np.array([0, 255, 255])

mask = np.ones(card.shape[:2], dtype=np.uint8)  # Start with all pixels relevant
mask &= np.all(card != card_edge_ignore_color, axis=-1).astype(np.uint8)
np.set_printoptions(threshold=np.inf)
mask &= np.all(card != card_value_ignore_color, axis=-1).astype(np.uint8)
c, w, h = card.shape[::-1]
print(screen.shape)
result = cv2.matchTemplate(screen, card, cv2.TM_CCOEFF_NORMED, mask=mask)
threshold = 0.93
loc = np.where( result >= threshold)

for pt in zip(*loc[::-1]):
    print(pt)
    if pt[1]>750:
        continue
    cv2.rectangle(screen, pt, (pt[0] + w, pt[1] + h), (0,0,255), 2)

min_val, max_val, min_loc, max_loc = cv2.minMaxLoc(result)

cv2.imshow("screen", screen)
cv2.waitKey(0)