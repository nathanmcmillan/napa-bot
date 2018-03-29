import math


def doji(candle):
    return True


def hammer(candle):
    return True


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