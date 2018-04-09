import math


def hammer(candle):
    body = abs(candle.open - candle.closing)
    wick = abs(min(candle.open, candle.closing) - candle.low)
    if wick > body * 2.0:
        if is_close(candle.closing, candle.high):
            return 'buy'
        elif is_close(candle.open, candle.high):
            return 'sell'
    return ''


def shooting_star(candle):
    body = abs(candle.open - candle.closing)
    wick = abs(max(candle.open, candle.closing) - candle.high)
    if wick > body * 2.0:
        if is_close(candle.open, candle.low):
            return 'buy'
        elif is_close(candle.closing, candle.low):
            return 'sell'
    return ''


def marubozu(candle):
    if is_close(candle.open, candle.low) and is_close(candle.closing, candle.high):
        return 'buy'
    if is_close(candle.open, candle.high) and is_close(candle.closing, candle.low):
        return 'sell'
    return ''


def trend(candles):
    if candles[0].closing < candles[-1].closing:
        return 'uptrend'
    else:
        return 'downtrend'


def is_close(a, b):
    relative = 1e-09
    absolute = 0.0
    return abs(a - b) <= max(relative * max(abs(a), abs(b)), absolute)
