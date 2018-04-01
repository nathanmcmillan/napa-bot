import logging
import sys
import signal
import time
import json
import gdax
import trading
import patterns
import printing
from trends import MovingAverage, ConvergeDiverge
from momentum import MoneyFlow, RelativeStrength, OnBalanceVolume
from safefile import SafeFile
from auth import Auth
from datetime import datetime
from datetime import timedelta


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


print('----------------------------------------')
print('|             napa candles             |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

run = True
file_out = '../candles.txt'
product = 'BTC-USD'
granularity = '3600'
time_interval = float(granularity)
time_format = '%Y-%m-%d %I:%M:%S %p'
candle_dictionary = {}
start = datetime.utcnow() - timedelta(days=365.0) 

while run:
    end = start + timedelta(seconds=time_interval)
    print('{} - {}'.format(start.strftime(time_format), end.strftime(time_format)))
    candles, status = gdax.get_candles(product, start.isoformat(), end.isoformat(), granularity)
    for candle in candles:
        candle_dictionary[candle.time] = candle
    time.sleep(5.0)
    start = end
    if start > datetime.utcnow():
        break
print('writing to file')
with open(file_out, "w+") as f:
    for key, candle in candle_dictionary.items():
        f.write('{} {:.2f} {:.2f} {:.2f} {:.2f} {:.2f}\n'.format(candle.time, candle.low, candle.high, candle.open, candle.closing, candle.volume))
print('finished')
    