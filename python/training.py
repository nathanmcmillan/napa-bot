import sys
import signal
import time
import json
import os.path
import gdax
from datetime import datetime
from datetime import timedelta


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


print('----------------------------------------')
print('|            napa training             |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

run = True
file_in = '../candles-btc-usd.txt'
file_out = '../training-btc-usd.txt'

neural_network = []

with open(file_in, "r") as f:
    for line in f:
        if not run:
            break
        candle = gdax.Candle(line.split())

print('writing to file')
with open(file_out, "w+") as f:
    f.write('hammer 0.12\n')
print('finished')
