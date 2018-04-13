import random
import patterns
from trends import ConvergeDiverge


def random_pattern():
    return random.choice(['red', 'green'])


def random_signal():
    return random.choice(['buy', 'sell'])


def mix_bool(a, b):
    if a and b:
        return True
    if not a and not b:
        return False
    return bool(random.getrandbits(1))


def mix_pattern(a, b):
    if a and b:
        return True
    if not a and not b:
        return False
    return bool(random.getrandbits(1))


def mix(a, b):
    now = Genetics()

    if 'maru' in a.buy and 'maru' in b.buy and a.buy['maru'] == b.buy['maru']:
        now.buy['maru'] = a.buy['maru']
    elif 'maru' in a.buy and bool(random.getrandbits(1)):
        now.buy['maru'] = a.buy['maru']
    elif 'maru' in b.buy and bool(random.getrandbits(1)):
        now.buy['maru'] = b.buy['maru']

    if 'maru' in a.sell and 'maru' in b.sell and a.sell['maru'] == b.sell['maru']:
        now.sell['maru'] = a.sell['maru']
    elif 'maru' in a.sell and bool(random.getrandbits(1)):
        now.sell['maru'] = a.sell['maru']
    elif 'maru' in b.sell and bool(random.getrandbits(1)):
        now.sell['maru'] = b.sell['maru']

    if 'trend' in a.buy and 'trend' in b.buy and a.buy['trend'] == b.buy['trend']:
        now.buy['trend'] = a.buy['trend']
        now.buy['trend_periods'] = int((a.buy['trend_periods'] + b.buy['trend_periods']) * 0.5)
    elif 'trend' in a.buy and bool(random.getrandbits(1)):
        now.buy['trend'] = a.buy['trend']
        now.buy['trend_periods'] = int((a.buy['trend_periods'] + random.randint(2, 21)) * 0.5)
    elif 'trend' in b.buy and bool(random.getrandbits(1)):
        now.buy['trend'] = b.buy['trend']
        now.buy['trend_periods'] = int((b.buy['trend_periods'] + random.randint(2, 21)) * 0.5)

    if 'trend' in a.sell and 'trend' in b.sell and a.sell['trend'] == b.sell['trend']:
        now.sell['trend'] = a.sell['trend']
        now.sell['trend_periods'] = int((a.sell['trend_periods'] + b.sell['trend_periods']) * 0.5)
    elif 'trend' in a.sell and bool(random.getrandbits(1)):
        now.sell['trend'] = a.sell['trend']
        now.sell['trend_periods'] = int((a.sell['trend_periods'] + random.randint(2, 21)) * 0.5)
    elif 'trend' in b.sell and bool(random.getrandbits(1)):
        now.sell['trend'] = b.sell['trend']
        now.sell['trend_periods'] = int((b.sell['trend_periods'] + random.randint(2, 21)) * 0.5)

    if 'macd' in a.buy and 'macd' in b.buy and a.buy['macd'] == b.buy['macd']:
        now.buy['macd'] = a.buy['macd']
    elif 'macd' in a.buy and bool(random.getrandbits(1)):
        now.buy['macd'] = a.buy['macd']
    elif 'macd' in b.buy and bool(random.getrandbits(1)):
        now.buy['macd'] = b.buy['macd']

    if 'macd' in a.sell and 'macd' in b.sell and a.sell['macd'] == b.sell['macd']:
        now.sell['macd'] = a.sell['macd']
    elif 'macd' in a.sell and bool(random.getrandbits(1)):
        now.sell['macd'] = a.sell['macd']
    elif 'macd' in b.sell and bool(random.getrandbits(1)):
        now.sell['macd'] = b.sell['macd']

    if 'color' in a.buy and 'color' in b.buy and a.buy['color'] == b.buy['color']:
        now.buy['color'] = a.buy['color']
    elif 'color' in a.buy and bool(random.getrandbits(1)):
        now.buy['color'] = a.buy['color']
    elif 'color' in b.buy and bool(random.getrandbits(1)):
        now.buy['color'] = b.buy['color']

    if 'color' in a.sell and 'color' in b.sell and a.sell['color'] == b.sell['color']:
        now.sell['color'] = a.sell['color']
    elif 'color' in a.sell and bool(random.getrandbits(1)):
        now.sell['color'] = a.sell['color']
    elif 'color' in b.sell and bool(random.getrandbits(1)):
        now.sell['color'] = b.sell['color']

    if 'star' in a.buy and 'star' in b.buy and a.buy['star'] == b.buy['star']:
        now.buy['star'] = a.buy['star']
    elif 'star' in a.buy and bool(random.getrandbits(1)):
        now.buy['star'] = a.buy['star']
    elif 'star' in b.buy and bool(random.getrandbits(1)):
        now.buy['star'] = b.buy['star']

    if 'star' in a.sell and 'star' in b.sell and a.sell['star'] == b.sell['star']:
        now.sell['star'] = a.sell['star']
    elif 'star' in a.sell and bool(random.getrandbits(1)):
        now.sell['star'] = a.sell['star']
    elif 'star' in b.sell and bool(random.getrandbits(1)):
        now.sell['star'] = b.sell['star']

    if 'hammer' in a.buy and 'hammer' in b.buy and a.buy['hammer'] == b.buy['hammer']:
        now.buy['hammer'] = a.buy['hammer']
    elif 'hammer' in a.buy and bool(random.getrandbits(1)):
        now.buy['hammer'] = a.buy['hammer']
    elif 'hammer' in b.buy and bool(random.getrandbits(1)):
        now.buy['hammer'] = b.buy['hammer']

    if 'hammer' in a.sell and 'hammer' in b.sell and a.sell['hammer'] == b.sell['hammer']:
        now.sell['hammer'] = a.sell['hammer']
    elif 'hammer' in a.sell and bool(random.getrandbits(1)):
        now.sell['hammer'] = a.sell['hammer']
    elif 'hammer' in b.sell and bool(random.getrandbits(1)):
        now.sell['hammer'] = b.sell['hammer']

    now.conditions['prevent_similar'] = mix_bool(a.conditions['prevent_similar'], b.conditions['prevent_similar'])
    now.conditions['buy_percent'] = (a.conditions['buy_percent'] + b.conditions['buy_percent']) * 0.5
    now.conditions['sell_percent'] = (a.conditions['sell_percent'] + b.conditions['sell_percent']) * 0.5

    return now


class Genetics:
    def __init__(self):
        self.buy = {}
        self.sell = {}
        self.conditions = {}
        self.conditions['similarity'] = 0.05
        self.conditions['ema_short'] = 12
        self.conditions['ema_long'] = 26

    def randomize(self):
        if bool(random.getrandbits(1)):
            self.buy['maru'] = random_pattern()
        if bool(random.getrandbits(1)):
            self.sell['maru'] = random_pattern()

        if bool(random.getrandbits(1)):
            self.buy['trend'] = random_pattern()
            self.buy['trend_periods'] = random.randint(2, 21)
        if bool(random.getrandbits(1)):
            self.sell['trend'] = random_pattern()
            self.sell['trend_periods'] = random.randint(2, 21)

        if bool(random.getrandbits(1)):
            self.buy['macd'] = random_signal()
        if bool(random.getrandbits(1)):
            self.sell['macd'] = random_signal()

        if bool(random.getrandbits(1)):
            self.buy['color'] = random_pattern()
        if bool(random.getrandbits(1)):
            self.sell['color'] = random_pattern()

        if bool(random.getrandbits(1)):
            self.buy['star'] = random_pattern()
        if bool(random.getrandbits(1)):
            self.sell['star'] = random_pattern()

        if bool(random.getrandbits(1)):
            self.buy['hammer'] = random_pattern()
        if bool(random.getrandbits(1)):
            self.sell['hammer'] = random_pattern()

        if bool(random.getrandbits(1)):
            self.buy['difference'] = random.random()
        if bool(random.getrandbits(1)):
            self.sell['difference'] = random.random()

        self.conditions['prevent_similar'] = bool(random.getrandbits(1))
        self.conditions['buy_percent'] = 0.1 + random.random() * 0.9
        self.conditions['sell_percent'] = random.random() * 1.1

    def signal(self, candles):
        buy = True
        sell = True

        if 'macd' in self.buy or 'macd' in self.sell:
            macd = ConvergeDiverge(self.conditions['ema_short'], self.conditions['ema_long'], candles[0].closing)
            candle_count = len(candles)
            for index in range(1, candle_count):
                current_candle = candles[index]
                macd.update(current_candle.closing)
            if 'macd' in self.buy and self.buy['macd'] != macd:
                buy = False
            if 'macd' in self.sell and self.sell['macd'] != macd:
                sell = False

        if 'maru' in self.buy or 'maru' in self.sell:
            maru = patterns.marubozu(candles[-1])
            if 'maru' in self.buy and self.buy['maru'] != maru:
                buy = False
            if 'maru' in self.sell and self.sell['maru'] != maru:
                sell = False

        if 'color' in self.buy or 'color' in self.sell:
            color = patterns.color(candles[-1])
            if 'color' in self.buy and self.buy['color'] != color:
                buy = False
            if 'color' in self.sell and self.sell['color'] != color:
                sell = False

        if 'hammer' in self.buy or 'hammer' in self.sell:
            hammer = patterns.hammer(candles[-1])
            if 'hammer' in self.buy and self.buy['hammer'] != hammer:
                buy = False
            if 'hammer' in self.sell and self.sell['hammer'] != hammer:
                sell = False

        if 'star' in self.buy or 'star' in self.sell:
            star = patterns.shooting_star(candles[-1])
            if 'star' in self.buy and self.buy['star'] != star:
                buy = False
            if 'star' in self.sell and self.sell['star'] != star:
                sell = False

        if 'trend' in self.buy:
            trend = patterns.trend(candles, self.buy['trend_periods'])
            if self.buy['trend'] != trend:
                buy = False

        if 'trend' in self.sell:
            trend = patterns.trend(candles, self.sell['trend_periods'])
            if self.sell['trend'] != trend:
                sell = False

        if 'difference' in self.buy or 'difference' in self.sell:
            difference = patterns.difference(candles[-1])
            if 'difference' in self.buy and self.buy['difference'] < difference:
                buy = False
            if 'difference' in self.sell and self.sell['difference'] < difference:
                sell = False

        if buy:
            return 'buy'
        if sell:
            return 'sell'
        return ''
