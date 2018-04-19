import random
import patterns
from trends import ConvergeDiverge

ema_short = 12
ema_long = 26
limit = 22


class GetMacd:
    def __init__(self):
        self.signal = None
        self.items = None

    def random(self):
        self.signal = random.choice(['buy', 'sell'])
        return self

    def get(self):
        macd = ConvergeDiverge(ema_short, ema_long, global_candles[0].closing)
        candle_count = len(global_candles)
        for index in range(1, candle_count):
            current_candle = global_candles[index]
            macd.update(current_candle.closing)
        return self.signal == macd.signal

    def to_string(self):
        return '{macd, signal: ' + self.signal + '}'


class GetTrend:
    def __init__(self):
        self.periods = 0
        self.pattern = None
        self.items = None

    def random(self):
        self.periods = random.randint(2, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.trend(global_candles, self.periods)

    def to_string(self):
        return '{trend, periods: ' + str(self.periods) + ', signal: ' + self.pattern + '}'


class GetColor:
    def __init__(self):
        self.period = 0
        self.pattern = None
        self.items = None

    def random(self):
        self.period = random.randint(2, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.color(global_candles[-self.period])

    def to_string(self):
        return '{color, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'


class GetMaru:
    def __init__(self):
        self.period = 0
        self.pattern = None
        self.items = None

    def random(self):
        self.period = random.randint(1, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.marubozu(global_candles[-self.period])

    def to_string(self):
        return '{maru, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'


class GetHammer:
    def __init__(self):
        self.period = 0
        self.pattern = None
        self.items = None

    def random(self):
        self.period = random.randint(1, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.hammer(global_candles[-self.period])

    def to_string(self):
        return '{hammer, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'


class GetStar:
    def __init__(self):
        self.period = 0
        self.pattern = None
        self.items = None

    def random(self):
        self.period = random.randint(1, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.shooting_star(global_candles[-self.period])

    def to_string(self):
        return '{star, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'


class GateAnd:
    def __init__(self):
        self.items = []

    def get(self):
        for item in self.items:
            if not item.get():
                return False
        return True

    def to_string(self):
        out = '{and:'
        for item in self.items:
            out += item.to_string()
        return out + '}'


class GateOr:
    def __init__(self):
        self.items = []

    def get(self):
        for item in self.items:
            if item.get():
                return True
        return False

    def to_string(self):
        out = '{or:'
        for item in self.items:
            out += item.to_string()
        return out + '}'


class GateXor:
    def __init__(self):
        self.items = []

    def get(self):
        result = False
        for item in self.items:
            if item.get():
                if result:
                    return False
                result = True
        return result

    def to_string(self):
        out = '{xor:'
        for item in self.items:
            out += item.to_string()
        return out + '}'


def random_gate():
    number = random.randint(0, 2)
    if number == 0:
        return GateAnd()
    if number == 1:
        return GateOr()
    if number == 2:
        return GateXor()


def random_signal():
    number = random.randint(0, 5)
    if number == 0:
        return GetMacd().random()
    if number == 1:
        return GetTrend().random()
    if number == 2:
        return GetColor().random()
    if number == 3:
        return GetMaru().random()
    if number == 4:
        return GetHammer().random()
    if number == 5:
        return GetStar().random()


def random_pattern():
    return random.choice(['red', 'green'])


def random_criteria():
    gate = random_gate()
    gate.items.append(random_signal())
    while bool(random.getrandbits(1)):
        gate.items.append(random_signal())
    while bool(random.getrandbits(1)):
        gate.items.append(random_criteria())
    return gate


def mix_bool(a, b):
    if a and b:
        return True
    if not a and not b:
        return False
    return bool(random.getrandbits(1))


def mix_criteria(a, b):
    if isinstance(a, GateAnd) and isinstance(b, GateAnd):
        gate = GateAnd()
    elif isinstance(a, GateOr) and isinstance(b, GateOr):
        gate = GateOr()
    elif isinstance(a, GateXor) and isinstance(b, GateXor):
        gate = GateXor()
    else:
        gate = random_gate()
    if a.items and b.items:
        size_a = len(a.items)
        size_b = len(b.items)
        for index in range(max(size_a, size_b)):
            if index < size_a and index < size_b:
                gate.items.append(mix_criteria(a.items[index], b.items[index]))
            elif index < size_a and bool(random.getrandbits(1)):
                gate.items.append(a.items[index])
            elif index < size_b and bool(random.getrandbits(1)):
                gate.items.append(b.items[index])
    elif a.items:
        for item in a.items:
            if bool(random.getrandbits(1)):
                gate.items.append(item)
    elif b.items:
        for item in b.items:
            if bool(random.getrandbits(1)):
                gate.items.append(item)
    return gate


def mix(a, b):
    now = Genetics()
    now.buy = mix_criteria(a.buy, b.buy)
    now.sell = mix_criteria(a.sell, b.sell)
    now.conditions['prevent_similar'] = mix_bool(a.conditions['prevent_similar'], b.conditions['prevent_similar'])
    now.conditions['buy_percent'] = (a.conditions['buy_percent'] + b.conditions['buy_percent']) * 0.5
    now.conditions['sell_percent'] = (a.conditions['sell_percent'] + b.conditions['sell_percent']) * 0.5
    return now


class Genetics:
    def __init__(self):
        self.buy = None
        self.sell = None
        self.conditions = {}
        self.conditions['similarity'] = 0.05
        self.conditions['ema_short'] = 12
        self.conditions['ema_long'] = 26

    def randomize(self):
        self.buy = random_criteria()
        self.sell = random_criteria()
        self.conditions['prevent_similar'] = bool(random.getrandbits(1))
        self.conditions['buy_percent'] = 0.1 + random.random() * 0.9
        self.conditions['sell_percent'] = random.random() * 1.1

    def signal(self, candles):
        global global_candles
        global_candles = candles
        if self.buy.get():
            return 'buy'
        if self.sell.get():
            return 'sell'
        return ''
