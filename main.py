#!/usr/bin/env python3
from typing import List, Tuple, Set, Any

import subprocess
import time
import datetime
import sys

from PIL import ImageGrab
import cv2
import numpy as np
import pyautogui

from solver import State, solve, Move, Places, CardCoordinates
from analyze import analyze

# Replace with the title of your window
window_title = "SHENZHEN I/O"

def capture_screen() -> Tuple[np.ndarray, int, int]:
    # Run xwininfo to get the window details
    try:
        result = subprocess.run(
            ["xwininfo", "-name", window_title],
            capture_output=True, text=True, check=True
        )
        output = result.stdout

        # Parse the output to get the window's position and size
        for line in output.split("\n"):
            if "Absolute upper-left X:" in line:
                x = int(line.split(":")[-1].strip())
            elif "Absolute upper-left Y:" in line:
                y = int(line.split(":")[-1].strip())
            elif "Width:" in line:
                width = int(line.split(":")[-1].strip())
            elif "Height:" in line:
                height = int(line.split(":")[-1].strip())

        # Capture the screen region using Pillow
        screenshot = ImageGrab.grab(bbox=(x, y, x + width, y + height))

        # Convert to NumPy array and OpenCV format
        screen_array = np.array(screenshot)
        screen_bgr = cv2.cvtColor(screen_array, cv2.COLOR_RGB2BGR)
        return (screen_bgr, x, y)

    except subprocess.CalledProcessError:
        print("Window not found. Please check the window title and try again.")
    except Exception as e:
        print(f"An error occurred: {e}")

OFFSET=10
def make_a_move(move:Move, screengeometry:Tuple[List[int], List[int], int], screenx: int, screeny: int) -> None:
    def calculate_position(place: Places, position: int) -> Tuple[int, int]:
        if place.value >= Places.T1.value and place.value <= Places.T8.value:
            if len(screengeometry[0]) > position:
                y = screengeometry[0][position - 1]
            else:
                y = screengeometry[0][4]+(screengeometry[0][1]-screengeometry[0][0])*(position-5)
            x = screengeometry[1][place.value - Places.T1.value]
            return x+screenx, y+screeny
        elif place.value >= Places.H1.value and place.value <= Places.H3.value:
            if len(screengeometry[2]) == 0:
                 y = screengeometry[0][0]-(screengeometry[0][4]-screengeometry[0][0])
            else:
                y = screengeometry[2][-1]
            x = screengeometry[1][place.value - Places.H1.value]
            return x+screenx, y+screeny
        elif place.value >= Places.S1.value and place.value <= Places.S3.value:
            if len(screengeometry[2]) == 0:
                 y = screengeometry[0][0]-(screengeometry[0][4]-screengeometry[0][0])
            else:
                y = screengeometry[2][-1]
            x = screengeometry[1][place.value - Places.S1.value + 5]
            return x+screenx, y+screeny
    if move.to.place in [Places.RBLOCK, Places.GBLOCK, Places.BBLOCK, Places.OBLOCK]:
        if move.to.place == Places.OBLOCK:
            pass
        else:
            x=screenx + screengeometry[1][3]+OFFSET
            y=screeny + 160 + (move.to.place.value - Places.RBLOCK.value) * 80
            pyautogui.moveTo(x, y, duration=0.1, tween=pyautogui.easeInOutQuad)
            time.sleep(0.1)
            pyautogui.click()
            pyautogui.doubleClick(x, y, duration=0.2)
            print(f"Click!")
            time.sleep(5)

    else:
        source = calculate_position(move.from_.place, move.from_.position)
        dest = calculate_position(move.to.place, move.to.position+1)
        print(f"Clicking at {source[0]+OFFSET}, {source[1]+OFFSET}")
        pyautogui.moveTo(source[0]+OFFSET, source[1]+OFFSET, duration=0.1, tween=pyautogui.easeInOutQuad)
        print(f"Dragging to {dest[0]+OFFSET}, {dest[1]+OFFSET}")
        pyautogui.dragTo(dest[0]+OFFSET, dest[1]+OFFSET, duration=0.5, button='left', tween=pyautogui.easeInOutQuad)


if __name__ == "__main__":
        filenamebase = f"captured/{datetime.datetime.now():%Y-%m-%d-%H:%M:%S}"
        screen = capture_screen()
        cv2.imwrite(filenamebase + ".png", screen[0])
        # Analyze the screen
        analisis = analyze(screen[0])
        print(analisis, screen[1:])
        with open(filenamebase + ".state", "w") as f:
            f.write(str(analisis[0].table))
        pyautogui.click(screen[1]+OFFSET, screen[2]+OFFSET, duration=0.1)

        #make_a_move(Move(CardCoordinates(Places.GBLOCK, 0), CardCoordinates(Places.GBLOCK, 0)), analisis[1], screen[1], screen[2])
        #sys.exit(0)
        moves = solve(analisis[0])
        for move in moves:
            print(move)
            make_a_move(move, analisis[1], screen[1], screen[2])
            time.sleep(1)
        #print("\n".join(map(str, moves)))
