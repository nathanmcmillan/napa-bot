import math


def hammer(candle):
    body = abs(candle.open - candle.closing)
    wick = abs(min(candle.open, candle.closing) - candle.low)
    if wick > body * 2.0:
        if is_close(candle.closing, candle.high):
            return 'green'
        elif is_close(candle.open, candle.high):
            return 'red'
    return ''


def shooting_star(candle):
    body = abs(candle.open - candle.closing)
    wick = abs(max(candle.open, candle.closing) - candle.high)
    if wick > body * 2.0:
        if is_close(candle.open, candle.low):
            return 'green'
        elif is_close(candle.closing, candle.low):
            return 'red'
    return ''


def marubozu(candle):
    if is_close(candle.open, candle.low) and is_close(candle.closing, candle.high):
        return 'green'
    if is_close(candle.open, candle.high) and is_close(candle.closing, candle.low):
        return 'red'
    return ''


def trend(candles, periods):
    if candles[-periods].closing < candles[-1].closing:
        return 'green'
    return 'red'


def color(candle):
    if candle.closing > candle.open:
        return 'green'
    return 'red'


def difference(candle):
    return candle.low / candle.high


def is_close(a, b):
    relative = 1e-09
    absolute = 0.0
    return abs(a - b) <= max(relative * max(abs(a), abs(b)), absolute)
