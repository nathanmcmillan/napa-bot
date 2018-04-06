import sys
import signal
import time
import json
import os.path
import gdax
import math
import neural
import patterns
from trends import MovingAverage, ConvergeDiverge
from momentum import MoneyFlow, RelativeStrength, OnBalanceVolume
from datetime import datetime
from datetime import timedelta


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


def prepare_input(candles):
    end = len(candles)
    '''
    low = candles[0].closing
    high = candles[0].closing
    for index in range(1, end):
        candle = candles[index]
        if candle.closing < low:
            low = candle.closing
        elif candle.closing > high:
            high = candle.closing
    price_range = high - low
    layer_in = []
    for candle in candles:
        percent = (candle.closing - low) / price_range
        layer_in.append(percent)
    return low, high, layer_in
    '''
    ema_short = 12
    ema_long = 26
    money_flow_index = MoneyFlow(ema_short)
    relative_strength_index = RelativeStrength(ema_short)
    balance_volume = OnBalanceVolume()
    macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
    money_flow_index.update(candles)
    relative_strength_index.update(candles)
    balance_volume.update(candles)
    for index in range(1, end):
        current_candle = candles[index]
        macd.update(current_candle.closing)
    # return 0.0, 0.0, [macd.current < 0.0, macd.current > 0.0, relative_strength_index.current < 0.3, relative_strength_index.current > 0.7]
    hammer = patterns.hammer(candles[-1])
    star = patterns.shooting_star(candles[-1])
    maru = patterns.marubozu(candles[-1])
    return 0.0, 0.0, [
        float(hammer[0] and hammer[1] == 'buy'),
        float(hammer[0] and hammer[1] == 'sell'),
        float(star[0] and star[1] == 'buy'),
        float(star[0] and star[1] == 'sell'),
        float(maru[0] and maru[1] == 'buy'),
        float(maru[0] and maru[1] == 'sell'),
        float(macd.signal == 'buy'),
        float(macd.signal == 'sell'),
        float(relative_strength_index.current < 0.2),
        float(relative_strength_index.current > 0.8)
    ]


def prepare_output(low, high, candle_end, candle_ahead):
    # return [(candle_ahead.closing - low) / (high - low)]
    return [float(candle_ahead.closing * 1.01 > candle_end.closing)]


print('----------------------------------------')
print('|            napa training             |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

run = True
file_in = '../candles-btc-usd.txt'
file_out = '../training-btc-usd.txt'

look_ahead = 1
batch = 28
candles = []

with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = gdax.Candle(line.split())
        candles.append(candle)

candle_count = len(candles)
network = neural.Network(30, [30, 20, 10], 1)
epochs = 5

for _ in range(epochs):
    error = 0.0
    start = 0
    end = batch
    print('training...', end=' ', flush=True)
    while run and end + look_ahead < candle_count:
        low, high, layer_in = prepare_input(candles[start:end])
        actual = prepare_output(low, high, candles[end], candles[end + look_ahead])
        network.set_input(layer_in)
        network.feed_forward()
        network.back_propagate(actual)
        error += network.get_error(actual)
        start += 1
        end += 1
    print('error:', error)

start = 0
end = batch
low, high, layer_in = prepare_input(candles[start:end])
actual = prepare_output(low, high, candles[end], candles[end + look_ahead])

prediction = network.predict(layer_in)
# print('predict', prediction[0] * (high - low) + low, 'actual', actual[0] * (high - low) + low)
print('predict', prediction[0], 'actual', actual[0], 'is', candles[end + look_ahead].closing * 1.01, '>', candles[end].closing)
'''
print('writing to file')
with open(file_out, "w+") as f:
    f.write('hammer 0.12\n')
print('finished')
'''