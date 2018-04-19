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


class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price


def prepare_input(candles, start, end):
    volume_low = candles[start].volume
    volume_high = candles[start].volume
    low = candles[start].low
    high = candles[start].high
    for index in range(start, end):
        candle = candles[index]
        if candle.low < low:
            low = candle.low
        if candle.high > high:
            high = candle.high
        if candle.volume < volume_low:
            volume_low = candle.volume
        elif candle.volume > volume_high:
            volume_high = candle.volume
    price_range = high - low
    volume_range = volume_high - volume_low
    parameters = []
    for index in range(start, end):
        candle = candles[index]
        parameters.append((candle.low - low) / price_range)
        parameters.append((candle.open - low) / price_range)
        parameters.append((candle.closing - low) / price_range)
        parameters.append((candle.high - low) / price_range)
        parameters.append((candle.volume - volume_low) / volume_range)
    return parameters
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
    '''


print('----------------------------------------')
print('|            napa training             |')
print('----------------------------------------')

run = True
file_in = '../candles-btc-usd.txt'
file_out = '../training-btc-usd.txt'

candles = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = gdax.Candle(line.split())
        if candle.time < 1513515600:
            continue
        candles.append(candle)
candle_end = len(candles) - 1

signal_out = []
signal_out.append([0.0, 0.0, 1.0])
for index in range(1, candle_end):
    if candles[index + 1].closing > candles[index].closing and candles[index - 1].closing > candles[index].closing:
        signal_out.append([1.0, 0.0, 0.0])
    elif candles[index].closing > candles[index + 1].closing and candles[index].closing > candles[index - 1].closing:
        signal_out.append([0.0, 1.0, 0.0])
    else:
        signal_out.append([0.0, 0.0, 1.0])

parameters_per_period = 5
period_range = 24  # 672 (28 days of hours)
parameters = parameters_per_period * period_range
end_price = candles[-1].closing
network = neural.Network(parameters, [parameters], 3)
epochs = 10

for _ in range(epochs):
    start = 0
    end = period_range
    funds = 1000.0
    orders = []
    error = 0.0
    print('training...', end=' ', flush=True)
    while run and end < candle_end:
        signal_in = prepare_input(candles, start, end)
        network.set_input(signal_in)
        network.feed_forward()
        network.back_propagate(signal_out[end])
        error += network.get_error(signal_out[end])
        signal = network.get_results()
        if signal[0] > signal[1] and signal[0] > signal[2] and funds > 20.0:
            ticker = candles[end - 1]
            buy_size = funds * 0.6
            funds -= buy_size
            orders.append(SimOrder(ticker.closing, None, buy_size))
        elif signal[1] > signal[0] and signal[1] > signal[2]:
            ticker = candles[end - 1]
            for order_to_sell in orders[:]:
                if ticker.closing > order_to_sell.coin_price:
                    funds += ticker.closing * order_to_sell.size
                    orders.remove(order_to_sell)
        start += 1
        end += 1
    worth = 0.0
    for order in orders:
        worth += order.size * end_price
    worth += funds
    print('total ${:.2f} error {}'.format(worth, error))

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
