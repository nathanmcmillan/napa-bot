import sys
import signal
import time
import json
import os.path
import binance
from datetime import datetime
from datetime import timedelta


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


print('----------------------------------------')
print('|         napa binance candles         |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

symbol = 'LTCBTC'
interval = '5m'
start = 1525670700000 - 360000
end = 1525670700000  # 1525670810
candles, status = binance.get_candles(symbol, interval, start, end)
print('candles...')
for candle in candles:
    print(candle.open_time, candle.closing, candle.number_of_trades)
'''
run = True
products = ['btcusdt', 'ethusdt', 'xlmeth']
granularity = '300'
file_out = '../' + product + '-' + granularity + '.txt'

candle_start_time = 2000000000
candle_end_time = -1
candle_dictionary = {}

if os.path.exists(file_out):
    with open(file_out, 'r') as f:
        for line in f:
            candle = gdax.Candle(line.split())
            candle_dictionary[candle.time] = candle
            if candle.time < candle_start_time:
                candle_start_time = candle.time
            if candle.time > candle_end_time:
                candle_end_time = candle.time

time_interval = float(granularity) * 200.0
time_format = '%Y-%m-%d %I:%M:%S %p'

# backwards
if candle_start_time > -1:
    end = datetime.utcfromtimestamp(candle_start_time)
    while run:
        start = end - timedelta(seconds=time_interval)
        print('{} - {}'.format(start.strftime(time_format), end.strftime(time_format)))
        candles, status = gdax.get_candles(product, start.isoformat(), end.isoformat(), granularity)
        if status != 200 or len(candles) == 0:
            break
        for candle in candles:
            candle_dictionary[candle.time] = candle
        time.sleep(1.0)
        end = start

# forwards
if candle_end_time == -1:
    start = datetime.utcnow() - timedelta(days=(365.0 * 3.0))
else:
    start = datetime.utcfromtimestamp(candle_end_time)
while run:
    end = start + timedelta(seconds=time_interval)
    print('{} - {}'.format(start.strftime(time_format), end.strftime(time_format)))
    candles, status = gdax.get_candles(product, start.isoformat(), end.isoformat(), granularity)
    if status != 200:
        print('something went wrong', status)
    for candle in candles:
        candle_dictionary[candle.time] = candle
    time.sleep(1.0)
    start = end
    if start > datetime.utcnow():
        break

print('writing to file')
with open(file_out, "w+") as f:
    for key, candle in sorted(candle_dictionary.items()):
        f.write('{} {:.2f} {:.2f} {:.2f} {:.2f} {:.2f}\n'.format(candle.time, candle.low, candle.high, candle.open, candle.closing, candle.volume))
print('finished')
'''