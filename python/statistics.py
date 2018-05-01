import sys
import signal
import time
import json
import os.path
import patterns
import genetics
import random
import simulation
from genetics import GetTrend
from gdax import Candle
from momentum import RelativeStrength
from trends import ConvergeDiverge
from genetics import Genetics
from operator import itemgetter

print('----------------------------------------')
print('|           napa statistics            |')
print('----------------------------------------')

BUY = 'buy'
SELL = 'sell'

UP_ONE = 1.001
DOWN_ONE = 1.0 - 0.001

GREEN_MARU_UP = 'green maru up'
RED_MARU_UP = 'red maru up'
GREEN_HAMMER_UP = 'green hammer up'
RED_HAMMER_DOWN = 'red hammer down'
GREEN_STAR_UP = 'green star up'
RED_STAR_DOWN = 'red star down'
CONTINUOUS_UP = 'continuous green'
CONTINUOUS_DOWN = 'continuous red'
DOUBLE_CONTINUOUS_UP = 'double continuous green'
DOUBLE_CONTINUOUS_DOWN = 'double continuous red'
RSI_LOW_UP = 'rsi low up'
RSI_HIGH_DOWN = 'rsi high down'
MACD_BUY_UP = 'macd buy up'
MACD_SELL_DOWN = 'macd sell down'


def action(candles, index):
    candle_len = len(candles)
    price = candles[index].closing
    up = price * UP_ONE
    down = price * DOWN_ONE
    index += 1
    while index < candle_len:
        closing = candles[index].closing
        if closing >= up:
            return BUY
        elif closing <= down:
            return SELL
        index += 1
    return ''


def stats(candles):
    ideas = {}
    ideas[GREEN_MARU_UP] = [0, 0]
    ideas[RED_MARU_UP] = [0, 0]
    ideas[GREEN_HAMMER_UP] = [0, 0]
    ideas[RED_HAMMER_DOWN] = [0, 0]
    ideas[GREEN_STAR_UP] = [0, 0]
    ideas[RED_STAR_DOWN] = [0, 0]
    ideas[CONTINUOUS_UP] = [0, 0]
    ideas[CONTINUOUS_DOWN] = [0, 0]
    ideas[DOUBLE_CONTINUOUS_UP] = [0, 0]
    ideas[DOUBLE_CONTINUOUS_DOWN] = [0, 0]
    ideas[MACD_BUY_UP] = [0, 0]
    ideas[MACD_SELL_DOWN] = [0, 0]
    #ideas[RSI_LOW_UP] = [0, 0]
    #ideas[RSI_HIGH_DOWN] = [0, 0]
    candle_len = len(candles) - 7
    index = 26
    while index < candle_len:
        prev_candle = candles[index - 1]
        candle = candles[index]
        next_candle = candles[index + 1]
        next_next_candle = candles[index + 2]

        signal = action(candles, index)

        # maru
        maru = patterns.marubozu(candle)
        if maru == 'green':
            ideas[GREEN_MARU_UP][1] += 1
            if next_candle.closing > candle.closing:
                ideas[GREEN_MARU_UP][0] += 1
        elif maru == 'red':
            ideas[RED_MARU_UP][1] += 1
            if next_candle.closing > candle.closing:
                ideas[RED_MARU_UP][0] += 1

        # hammer
        hammer = patterns.hammer(candle)
        if hammer == 'green':
            ideas[GREEN_HAMMER_UP][1] += 1
            if next_candle.closing > candle.closing:
                ideas[GREEN_HAMMER_UP][0] += 1
        elif hammer == 'red':
            ideas[RED_HAMMER_DOWN][1] += 1
            if next_candle.closing < candle.closing:
                ideas[RED_HAMMER_DOWN][0] += 1

        # star
        star = patterns.shooting_star(candle)
        if star == 'green':
            ideas[GREEN_STAR_UP][1] += 1
            if next_candle.closing > candle.closing:
                ideas[GREEN_STAR_UP][0] += 1
        elif star == 'red':
            ideas[RED_STAR_DOWN][1] += 1
            if next_candle.closing < candle.closing:
                ideas[RED_STAR_DOWN][0] += 1

        # continous
        if prev_candle.closing < candle.closing:
            ideas[CONTINUOUS_UP][1] += 1
            if candle.closing < next_candle.closing:
                ideas[CONTINUOUS_UP][0] += 1
                ideas[DOUBLE_CONTINUOUS_UP][1] += 1
                if next_candle.closing < next_next_candle.closing:
                    ideas[DOUBLE_CONTINUOUS_UP][0] += 1
        elif prev_candle.closing > candle.closing:
            ideas[CONTINUOUS_DOWN][1] += 1
            if candle.closing > next_candle.closing:
                ideas[CONTINUOUS_DOWN][0] += 1
                ideas[DOUBLE_CONTINUOUS_DOWN][1] += 1
                if next_candle.closing > next_next_candle.closing:
                    ideas[DOUBLE_CONTINUOUS_DOWN][0] += 1

        # macd
        macd = ConvergeDiverge(12, 26, candles[index - 26].closing)
        for jindex in range(index - 25, index):
            macd.update(candles[jindex].closing)
        if macd.signal == 'buy':
            ideas[MACD_BUY_UP][1] += 1
            if candles[index + 6].closing > candle.closing * UP_ONE:
                ideas[MACD_BUY_UP][0] += 1
        elif macd.signal == 'sell':
            ideas[MACD_SELL_DOWN][1] += 1
            if candles[index + 6].closing < candle.closing * DOWN_ONE:
                ideas[MACD_SELL_DOWN][0] += 1

        # relative strength index
        '''rsi = RelativeStrength(14)
        rsi.update(candles[index - 16:index])
        if rsi == 'buy':
            ideas[RSI_LOW_UP][1] += 1
            if candle.closing < next_candle.closing:
                ideas[RSI_LOW_UP][0] += 1
        elif rsi == 'sell':
            ideas[RSI_HIGH_DOWN][1] += 1
            if candle.closing > next_candle.closing:
                ideas[RSI_HIGH_DOWN][0] += 1'''

        index += 1

    print()
    for key, value in ideas.items():
        if value[1] == 0:
            print('{0:30} true {1} false {2}'.format(key, '-', '-'))
            continue
        percent = float(value[0]) / float(value[1]) * 100.0
        print('{0:30} true {1:.2f}% false {2:.2f}%'.format(key, percent, 100.0 - percent))


file_in = '../BTC-USD-300.txt'
candles = {}
candles['5 minute'] = []
candles['30 minute'] = []
candles['1 hour'] = []
candles['6 hour'] = []
candles['1 day'] = []
candles['7 day'] = []
candles['30 day'] = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        candles['5 minute'].append(candle)
        if candle.time % 1800 == 0:
            candles['30 minute'].append(candle)
        if candle.time % 3600 == 0:
            candles['1 hour'].append(candle)
        if candle.time % 21600 == 0:
            candles['6 hour'].append(candle)
        if candle.time % 86400 == 0:
            candles['1 day'].append(candle)
        if candle.time % 604800 == 0:
            candles['7 day'].append(candle)
        if candle.time % 2592000 == 0:
            candles['30 day'].append(candle)

for interval, data in candles.items():
    print(interval, 'candles ({})'.format(len(data)), end='', flush=True)
    stats(data)
    print('----------------------------------------')