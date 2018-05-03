import signal
import time
import json
import os.path
import patterns
import trends
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

MIN_UP = 1.0 + 0.01
MIN_DOWN = 1.0 - 0.01

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
RSI_BUY_UP = 'rsi buy up'
RSI_SELL_DOWN = 'rsi sell down'
MACD_BUY_UP = 'macd buy up'
MACD_SELL_DOWN = 'macd sell down'
BREAK_RESISTANCE_UP = 'break resistance up'
BOUNCE_RESISTANCE_DOWN = 'bounce reistance down'
BREAK_SUPPORT_DOWN = 'break support down'
BOUNCE_SUPPORT_UP = 'bounce support up'
LIQUIDATION_UP = 'liquidation up'


def action(candles, index):
    candle_len = len(candles)
    price = candles[index].closing
    up = price * MIN_UP
    down = price * MIN_DOWN
    index += 1
    while index < candle_len:
        closing = candles[index].closing
        if closing >= up:
            return BUY
        elif closing <= down:
            return SELL
        index += 1
    return ''


def signals(candles):
    return False


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
    ideas[RSI_BUY_UP] = [0, 0]
    ideas[RSI_SELL_DOWN] = [0, 0]
    ideas[BREAK_RESISTANCE_UP] = [0, 0]
    ideas[BOUNCE_RESISTANCE_DOWN] = [0, 0]  # support & resistance statistics should not be measured per candle
    ideas[BREAK_SUPPORT_DOWN] = [0, 0]
    ideas[BOUNCE_SUPPORT_UP] = [0, 0]
    ideas[LIQUIDATION_UP] = [0, 0]
    candle_len = len(candles)
    index = 26
    while index < candle_len:
        previous_candle = candles[index - 1]
        candle = candles[index]
        signal = action(candles, index)

        # maru
        maru = patterns.marubozu(candle)
        if maru == 'green':
            ideas[GREEN_MARU_UP][1] += 1
            if signal == BUY:
                ideas[GREEN_MARU_UP][0] += 1
        elif maru == 'red':
            ideas[RED_MARU_UP][1] += 1
            if signal == BUY:
                ideas[RED_MARU_UP][0] += 1

        # hammer
        hammer = patterns.hammer(candle)
        if hammer == 'green':
            ideas[GREEN_HAMMER_UP][1] += 1
            if signal == BUY:
                ideas[GREEN_HAMMER_UP][0] += 1
        elif hammer == 'red':
            ideas[RED_HAMMER_DOWN][1] += 1
            if signal == SELL:
                ideas[RED_HAMMER_DOWN][0] += 1

        # star
        star = patterns.shooting_star(candle)
        if star == 'green':
            ideas[GREEN_STAR_UP][1] += 1
            if signal == BUY:
                ideas[GREEN_STAR_UP][0] += 1
        elif star == 'red':
            ideas[RED_STAR_DOWN][1] += 1
            if signal == SELL:
                ideas[RED_STAR_DOWN][0] += 1

        # continous
        if candle.open < candle.closing:
            ideas[CONTINUOUS_UP][1] += 1
            if signal == BUY:
                ideas[CONTINUOUS_UP][0] += 1
            if previous_candle.open < previous_candle.closing:
                ideas[DOUBLE_CONTINUOUS_UP][1] += 1
                if signal == BUY:
                    ideas[DOUBLE_CONTINUOUS_UP][0] += 1
        elif candle.open > candle.closing:
            ideas[CONTINUOUS_DOWN][1] += 1
            if signal == SELL:
                ideas[CONTINUOUS_DOWN][0] += 1
            if previous_candle.open > previous_candle.closing:
                ideas[DOUBLE_CONTINUOUS_DOWN][1] += 1
                if signal == SELL:
                    ideas[DOUBLE_CONTINUOUS_DOWN][0] += 1

        # macd
        macd = ConvergeDiverge(12, 26, candles[index - 26].closing)
        for jindex in range(index - 25, index):
            macd.update(candles[jindex].closing)
        if macd.signal == 'buy':
            ideas[MACD_BUY_UP][1] += 1
            if signal == BUY:
                ideas[MACD_BUY_UP][0] += 1
        elif macd.signal == 'sell':
            ideas[MACD_SELL_DOWN][1] += 1
            if signal == SELL:
                ideas[MACD_SELL_DOWN][0] += 1

        # rsi
        rsi = RelativeStrength(14)
        rsi.update(candles, index)
        if rsi.signal == 'buy':
            ideas[RSI_BUY_UP][1] += 1
            if signal == BUY:
                ideas[RSI_BUY_UP][0] += 1
        elif rsi.signal == 'sell':
            ideas[RSI_SELL_DOWN][1] += 1
            if signal == SELL:
                ideas[RSI_SELL_DOWN][0] += 1

        # support
        support = trends.support(candles, index - 26, index - 1)
        if support and candle.closing < support:
            ideas[BREAK_SUPPORT_DOWN][1] += 1
            if signal == SELL:
                ideas[BREAK_SUPPORT_DOWN][0] += 1

        # resistance
        resistance = trends.resistance(candles, index - 26, index - 1)
        if resistance and candle.closing > resistance:
            ideas[BREAK_RESISTANCE_UP][1] += 1
            if signal == BUY:
                ideas[BREAK_RESISTANCE_UP][0] += 1

        # liquidation
        liquid = trends.liquidation(candles, index - 6, index - 3, index)
        if liquid:
            ideas[LIQUIDATION_UP][1] += 1
            if signal == BUY:
                ideas[LIQUIDATION_UP][0] += 1

        index += 1

    print()
    for key, value in ideas.items():
        if value[1] == 0:
            print('{0:30} -'.format(key))
            continue
        percent = float(value[0]) / float(value[1]) * 100.0
        print('{0:30} true {1:.2f}% false {2:.2f}% ({3:,})'.format(key, percent, 100.0 - percent, value[1]))


bear = False
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
        if candle.time < 1513515600 and bear:
            continue
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
    print(interval, 'candles ({:,})'.format(len(data)), end='', flush=True)
    stats(data)
    print('----------------------------------------')