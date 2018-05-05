import patterns
import trends
from momentum import RelativeStrength
from trends import ConvergeDiverge


class Strategy:
    def __init__(self, name):
        self.name = name
        self.buy = None
        self.sell = nothing
        self.stop_limit = None

    def amount(self, funds):
        probability = 0.55
        profit = 1.01
        percent = (probability * profit - (1.0 - probability)) / profit
        if percent < 0.001:
            return 0.0
        return funds * percent


def no_loss(candles, index, order):
    closing = candles[index].closing
    if closing < order.coin_price:
        order.stop_limit = 0.0
    else:
        low = closing * 0.80
        if low > order.stop_limit:
            order.stop_limit = low


def large_trail(candles, index, order):
    low = candles[index].closing * 0.80
    if low > order.stop_limit:
        order.stop_limit = low


def simple_trail(candles, index, order):
    low = candles[index].closing * 0.95
    if low > order.stop_limit:
        order.stop_limit = low


def chandelier(candles, index, order):
    highest = candles[index - 22].high
    for jindex in range(index - 21, index):
        high = candles[jindex].high
        if high > highest:
            highest = high
    average_true_range = 0.0
    for jindex in range(index - 22, index):
        average_true_range += trends.true_range(candles[jindex])
    average_true_range /= 22.0
    order.stop_limit = highest - average_true_range * 3.0


def nothing(candles, index):
    return False


def continue_trend(candles, index):
    return candles[index - 7].open < candles[index].closing


def green_maru(candles, index):
    return patterns.marubozu(candles[index]) == 'green'


def green_hammer(candles, index):
    return patterns.hammer(candles[index]) == 'green'


def green_star(candles, index):
    return patterns.shooting_star(candles[index]) == 'green'


def macd_buy(candles, index):
    macd = ConvergeDiverge(12, 26, candles[index - 26].closing)
    for jindex in range(index - 25, index):
        macd.update(candles[jindex].closing)
    return macd.signal == 'buy'


def rsi_buy(candles, index):
    rsi = RelativeStrength(14)
    rsi.update(candles, index)
    return rsi.signal == 'buy'


def liquidation_drop(candles, index):
    return trends.liquidation(candles, index - 6, index - 3, index)


def resistance_breakout(candles, index):
    line = trends.resistance(candles, index - 26, index - 1)
    if line:
        return candles[index].closing > line
    else:
        return False


def trend_and_maru(candles, index):
    a = patterns.marubozu(candles[index]) == 'green'
    b = patterns.marubozu(candles[index - 1]) == 'green'
    return candles[index - 1].open < candles[index].closing and (a or b)
