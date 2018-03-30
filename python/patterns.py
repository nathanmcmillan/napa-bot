import math


def doji(candle):
    return True


def hammer(candle):
    if math.isclose(candle.closing, candle.high):
        body = abs(candle.closing - candle.open)
        wick = candle.closing - candle.low
        if wick > body * 2.0:
            return (True, 'white')
    return (False, '')


def invert_hammer(candle):
    return True


def marubozu(candle):
    if math.isclose(candle.open, candle.low) and math.isclose(candle.closing, candle.high):
        return (True, 'white')
    if math.isclose(candle.open, candle.high) and math.isclose(candle.closing, candle.low):
        return (True, 'black')
    return (False, '')


def shooting_star(candle):
    return True


def spin_top(candle):
    return True
