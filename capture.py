#!/usr/bin/env python3
import cv2
import numpy as np

screen = cv2.imread("test/test1.png")
card = cv2.imread("test/card.png")
card_edge_ignore_color = np.array([255, 0, 255])
card_value_ignore_color = np.array([255, 255, 0])
card_value_mask = np.all(card == card_value_ignore_color, axis=-1).astype(np.uint8)

y_indices, x_indices = np.where(card_value_mask)
if y_indices.size > 0 and x_indices.size > 0:
    x_min, x_max = x_indices.min(), x_indices.max()
    y_min, y_max = y_indices.min(), y_indices.max()
else:
    print("No cyan areas detected.")
    exit()


mask = np.ones(card.shape[:2], dtype=np.uint8)  # Start with all pixels relevant
mask &= np.all(card != card_edge_ignore_color, axis=-1).astype(np.uint8)
np.set_printoptions(threshold=np.inf)
mask &= ~card_value_mask
c, w, h = card.shape[::-1]
result = cv2.matchTemplate(screen, card, cv2.TM_CCOEFF_NORMED, mask=mask)
threshold = 0.93
loc = np.where( result >= threshold)
counter = 1
for pt in zip(*loc[::-1]):
    x, y = pt
    print(x,y)
    if y>750:
        continue
    cropped_area = screen[y + y_min:y + y_max + 1, x + x_min:x + x_max + 1]
    cropped_area_gray = cv2.cvtColor(cropped_area, cv2.COLOR_BGR2GRAY)
    _, cropped_area_bw = cv2.threshold(cropped_area_gray, 127, 255, cv2.THRESH_BINARY)
    cv2.imwrite(f"dump/cyan_area_{counter}.png", cropped_area)
    cv2.imwrite(f"dump/cyan_area_grey_{counter}.png", cropped_area_bw)
    counter += 1    

for pt in zip(*loc[::-1]):
    cv2.rectangle(screen, pt, (pt[0] + w, pt[1] + h), (0,0,255), 2)

min_val, max_val, min_loc, max_loc = cv2.minMaxLoc(result)

cv2.imshow("screen", screen)
cv2.waitKey(0)