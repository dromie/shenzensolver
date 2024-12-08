#!/usr/bin/env python3
import subprocess
from PIL import ImageGrab
import cv2
import numpy as np

from solver import State, solve
from analyze import analyze

# Replace with the title of your window
window_title = "SHENZHEN I/O"

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

    # Analyze the screen
    state = analyze(screen_bgr)
    print(state)
    print("\n".join(map(str,solve(state))))

except subprocess.CalledProcessError:
    print("Window not found. Please check the window title and try again.")
except Exception as e:
    print(f"An error occurred: {e}")
