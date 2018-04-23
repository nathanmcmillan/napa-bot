import random
import patterns
from trends import ConvergeDiverge

ema_short = 12
ema_long = 26
limit = 22


class GetMacd:
    def __init__(self):
        self.name = 'macd'
        self.signal = None

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

    def copy(self):
        dna = GetMacd()
        dna.signal = self.signal
        return dna


class GetTrend:
    def __init__(self):
        self.name = 'trend'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(2, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.trend(global_candles, 0, self.period)

    def to_string(self):
        return '{trend, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def copy(self):
        dna = GetTrend()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetColor:
    def __init__(self):
        self.name = 'color'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(2, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.color(global_candles[-self.period])

    def to_string(self):
        return '{color, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def copy(self):
        dna = GetColor()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetMaru:
    def __init__(self):
        self.name = 'maru'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(1, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.marubozu(global_candles[-self.period])

    def to_string(self):
        return '{maru, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def copy(self):
        dna = GetMaru()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetHammer:
    def __init__(self):
        self.name = 'hammer'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(1, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.hammer(global_candles[-self.period])

    def to_string(self):
        return '{hammer, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def copy(self):
        dna = GetHammer()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetStar:
    def __init__(self):
        self.name = 'star'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(1, limit)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.shooting_star(global_candles[-self.period])

    def to_string(self):
        return '{star, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def copy(self):
        dna = GetStar()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


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


def mix_bool(a, b):
    if a and b:
        return True
    if not a and not b:
        return False
    return bool(random.getrandbits(1))


def random_criteria(criteria):
    signal = random_signal()
    criteria[signal.name] = signal


def union(criteria, a, b):
    for key, value in a.items():
        criteria[key] = value
    for key, value in b.items():
        if key not in criteria:
            criteria[key] = value


def intersection(criteria, a, b):
    for key, value in a.items():
        if key in b:
            criteria[key] = value


def permutate(a, b):
    permutations = []

    gene = Genetics()
    intersection(gene.buy, a.buy, b.buy)
    intersection(gene.sell, a.sell, b.sell)
    gene.conditions['prevent_similar'] = mix_bool(a.conditions['prevent_similar'], b.conditions['prevent_similar'])
    gene.conditions['buy_percent'] = (a.conditions['buy_percent'] + b.conditions['buy_percent']) * 0.5
    gene.conditions['sell_percent'] = (a.conditions['sell_percent'] + b.conditions['sell_percent']) * 0.5
    permutations.append(gene)

    gene = Genetics()
    union(gene.buy, a.buy, b.buy)
    union(gene.sell, a.sell, b.sell)
    gene.conditions['prevent_similar'] = mix_bool(a.conditions['prevent_similar'], b.conditions['prevent_similar'])
    gene.conditions['buy_percent'] = (a.conditions['buy_percent'] + b.conditions['buy_percent']) * 0.5
    gene.conditions['sell_percent'] = (a.conditions['sell_percent'] + b.conditions['sell_percent']) * 0.5
    permutations.append(gene)

    return permutations


class Genetics:
    def __init__(self):
        self.buy = {}
        self.sell = {}
        self.conditions = {}

    def randomize(self):
        random_criteria(self.buy)
        random_criteria(self.sell)
        self.conditions['prevent_similar'] = bool(random.getrandbits(1))
        self.conditions['buy_percent'] = 0.1 + random.random() * 0.9
        self.conditions['sell_percent'] = random.random() * 1.1

    def signal(self, candles):
        global global_candles
        global_candles = candles
        success = True
        for _, criteria in self.buy.items():
            if not criteria.get():
                success = False
                break
        if success:
            return 'buy'
        for _, criteria in self.sell.items():
            if not criteria.get():
                return ''
        return 'sell'
