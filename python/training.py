import sys
import signal
import time
import json
import os.path
import gdax
import math
import neural
import patterns
from trends import MovingAverage, ConvergeDiverge, AverageDirectional
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
    relative_strength_index = RelativeStrength(ema_short)
    directional_index = AverageDirectional(ema_short)
    macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
    relative_strength_index.update(candles)
    directional_index.update(candles)
    for index in range(1, end):
        current_candle = candles[index]
        macd.update(current_candle.closing)
    hammer = patterns.hammer(candles[-1])
    star = patterns.shooting_star(candles[-1])
    maru = patterns.marubozu(candles[-1])
    return 0.0, 0.0, [
        float(hammer == 'buy'),
        float(hammer == 'sell'),
        float(star == 'buy'),
        float(star == 'sell'),
        float(maru == 'buy'),
        float(maru == 'sell'),
        float(macd.signal == 'buy'),
        float(macd.signal == 'sell'),
        float(relative_strength_index.current < 0.2),
        float(relative_strength_index.current > 0.8),
        float(directional_index.current < 0.2),
        float(directional_index.current > 0.4)
    ]


def prepare_output(low, high, candles, start, limit):
    # return [(candle_ahead.closing - low) / (high - low)]
    # return [float(candle_ahead.closing > candle_end.closing), float(candle_ahead.closing < candle_end.closing)]
    # return [candle_ahead.closing * 1.01 > candle_end.closing]
    price = candles[start].closing
    for index in range(start + 1, start + limit):
        candle = candles[index]
        if candle.closing * 1.01 > price:
            return [1.0, 0.0, 0.0]
        if candle.closing < price * 0.99:
            return [0.0, 1.0, 0.0]
    return [0.0, 0.0, 1.0]


print('----------------------------------------')
print('|            napa training             |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

run = True
file_in = '../candles-btc-usd.txt'
file_out = '../training-btc-usd.txt'

batch = 28
candles = []

with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = gdax.Candle(line.split())
        candles.append(candle)

candle_count = len(candles)
network = neural.Network(12, [12], 3)
epochs = 20
look_ahead = 6

for _ in range(epochs):
    error = 0.0
    start = 0
    end = batch
    print('training...', end=' ', flush=True)
    while run and end + look_ahead < candle_count:
        low, high, layer_in = prepare_input(candles[start:end])
        actual = prepare_output(low, high, candles, end, look_ahead)
        network.set_input(layer_in)
        network.feed_forward()
        network.back_propagate(actual)
        error += network.get_error(actual)
        start += 1
        end += 1
    print('error:', error)

start = 0
end = batch

for _ in range(200):
    low, high, layer_in = prepare_input(candles[start:end])
    actual = prepare_output(low, high, candles, end, look_ahead)
    prediction = network.predict(layer_in)
    # print('predict', prediction[0] * (high - low) + low, 'actual', actual[0] * (high - low) + low)
    print('predict', prediction, 'actual', actual, 'is', candles[end + look_ahead].closing, '--', candles[end].closing)
    start += 1
    end += 1

print('writing to file')
with open(file_out, "w+") as f:
    size = len(network.layers)
    for index in range(1, size):
        current_layer = network.layers[index]
        for neuron in current_layer:
            for synapse in neuron.synapses:
                f.write(str(synapse.weight) + ' ')
            f.write('\n')
        f.write('\n')
print('finished')
