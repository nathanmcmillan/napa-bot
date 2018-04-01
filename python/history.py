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
print('|             napa candles             |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

run = True
file_out = '../candles-btc-usd.txt'
product = 'BTC-USD'
granularity = '3600'
time_interval = float(granularity) * 200.0
time_format = '%Y-%m-%d %I:%M:%S %p'
candle_dictionary = {}
start = datetime.utcnow() - timedelta(days=365.0) 

if os.path.exists(file_out):
    with open(file_out, "r") as f:
        for line in f:
            candle = gdax.Candle(line.split())
            candle_dictionary[candle.time] = candle

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
    for key, candle in sorted(candle_dictionary.items()):
        f.write('{} {:.2f} {:.2f} {:.2f} {:.2f} {:.2f}\n'.format(candle.time, candle.low, candle.high, candle.open, candle.closing, candle.volume))
print('finished')
