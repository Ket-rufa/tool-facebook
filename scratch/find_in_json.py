import json
import sys

def find_string(data, target, path=""):
    if isinstance(data, dict):
        for k, v in data.items():
            find_string(v, target, f"{path}['{k}']")
    elif isinstance(data, list):
        for i, v in enumerate(data):
            find_string(v, target, f"{path}[{i}]")
    elif isinstance(data, str):
        if target.lower() in data.lower():
            try:
                print(f"FOUND at {path}: {data}")
            except UnicodeEncodeError:
                print(f"FOUND at {path}: (non-ascii string found)")

target_strings = ["Hồ Chí Minh", "Đắk Lắk", "Thompson", "Victoria"]

if len(sys.argv) > 1:
    filename = sys.argv[1]
else:
    filename = "debug_about_61572157625165.json"

with open(filename, "r", encoding="utf-8") as f:
    for line_num, line in enumerate(f, 1):
        line = line.strip()
        if not line: continue
        try:
            # Handle potential "for (;;);" prefix
            if line.startswith("for (;;);"):
                line = line[len("for (;;);"):]
            data = json.loads(line)
            print(f"--- Line {line_num} ---")
            for ts in target_strings:
                find_string(data, ts)
        except Exception as e:
            print(f"Error parsing line {line_num}: {e}")
