import random
import patterns
from trends import ConvergeDiverge

candles = None
index = 0


class GetMacd:
    def __init__(self):
        self.name = 'macd'
        self.signal = None

    def random(self):
        self.signal = random.choice(['buy', 'sell'])
        return self

    def get(self):
        macd = ConvergeDiverge(12, 26, candles[index - 21].closing)
        for jindex in range(index - 20, index):
            macd.update(candles[jindex].closing)
        return self.signal == macd.signal

    def to_string(self):
        return '{macd, signal: ' + self.signal + '}'

    def key(self):
        return (self.name, self.signal)

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
        self.period = random.randint(1, 20)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.trend(candles, index - self.period, index)

    def to_string(self):
        return '{trend, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def key(self):
        return (self.name, self.period, self.pattern)

    def copy(self):
        dna = GetTrend()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetVolume:
    def __init__(self):
        self.name = 'volume'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(1, 20)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.volume_trend(candles, index - self.period, index)

    def to_string(self):
        return '{volume, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def key(self):
        return (self.name, self.period, self.pattern)

    def copy(self):
        dna = GetVolume()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetChange:
    def __init__(self):
        self.name = 'percent'
        self.period = 0
        self.percent = 0
        self.float_percent = 0.0

    def random(self):
        self.period = random.randint(1, 20)
        self.percent = random.randint(1, 500)
        self.float_percent = float(self.percent) / 1000.0
        return self

    def get(self):
        return patterns.change(candles, index - self.period, index) > self.float_percent

    def to_string(self):
        return '{change, period: ' + str(self.period) + ', percent: ' + str(self.percent) + '}'

    def key(self):
        return (self.name, self.period, self.percent)

    def copy(self):
        dna = GetChange()
        dna.period = self.period
        dna.percent = self.percent
        return dna


class GetColor:
    def __init__(self):
        self.name = 'color'
        self.period = 0
        self.pattern = None

    def random(self):
        self.period = random.randint(0, 6)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.color(candles[index - self.period])

    def to_string(self):
        return '{color, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def key(self):
        return (self.name, self.period, self.pattern)

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
        self.period = random.randint(0, 6)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.marubozu(candles[index - self.period])

    def to_string(self):
        return '{maru, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def key(self):
        return (self.name, self.period, self.pattern)

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
        self.period = random.randint(0, 6)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.hammer(candles[index - self.period])

    def to_string(self):
        return '{hammer, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def key(self):
        return (self.name, self.period, self.pattern)

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
        self.period = random.randint(0, 6)
        self.pattern = random_pattern()
        return self

    def get(self):
        return self.pattern == patterns.shooting_star(candles[index - self.period])

    def to_string(self):
        return '{star, period: ' + str(self.period) + ', signal: ' + self.pattern + '}'

    def key(self):
        return (self.name, self.period, self.pattern)

    def copy(self):
        dna = GetStar()
        dna.period = self.period
        dna.pattern = self.pattern
        return dna


class GetSupport:
    def __init__(self):
        self.name = 'support'
        self.period = 0

    def random(self):
        self.period = random.randint(0, 20)
        return self

    def get(self):
        low = candles[index].closing
        for jindex in range(index - self.period, index - 1):
            if candles[jindex].closing < low:
                low = candles[jindex].closing
        return candles[index].closing > low

    def to_string(self):
        return '{support, period: ' + str(self.period) + '}'

    def key(self):
        return (self.name, self.period)

    def copy(self):
        dna = GetSupport()
        dna.period = self.period
        return dna


class GetResistance:
    def __init__(self):
        self.name = 'resistance'
        self.period = 0

    def random(self):
        self.period = random.randint(0, 20)
        return self

    def get(self):
        high = candles[index].closing
        for jindex in range(index - self.period, index - 1):
            if candles[jindex].closing > high:
                high = candles[jindex].closing
        return candles[index].closing < high

    def to_string(self):
        return '{resistance, period: ' + str(self.period) + '}'

    def key(self):
        return (self.name, self.period)

    def copy(self):
        dna = GetResistance()
        dna.period = self.period
        return dna


def random_signal():
    number = random.randint(0, 9)
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
    if number == 6:
        return GetChange().random()
    if number == 7:
        return GetVolume().random()
    if number == 8:
        return GetSupport().random()
    if number == 9:
        return GetResistance().random()


def random_pattern():
    return random.choice(['red', 'green'])


def random_criteria(criteria):
    signal = random_signal()
    criteria[signal.key()] = signal
    while bool(random.getrandbits(1)):
        signal = random_signal()
        criteria[signal.key()] = signal


def equals(gene_a, gene_b):
    for key in gene_a.buy:
        if key not in gene_b.buy:
            return False
    for key in gene_b.buy:
        if key not in gene_a.buy:
            return False
    for key in gene_a.sell:
        if key not in gene_b.sell:
            return False
    for key in gene_b.sell:
        if key not in gene_a.sell:
            return False
    return True


def union(criteria, a, b):
    for key, value in a.items():
        criteria[key] = value
    for key, value in b.items():
        criteria[key] = value


def intersection(criteria, a, b):
    for key, value in a.items():
        if key in b:
            criteria[key] = value
        else:
            for key2, value2 in b.items():
                if key[0] == key2[0]:
                    if bool(random.getrandbits(1)):
                        criteria[key] = value
                    else:
                        criteria[key2] = value2


def copy_criteria(criteria, a):
    for key, value in a.items():
        criteria[key] = value


def permutate(a, b):
    permutations = []

    gene = Genetics()
    intersection(gene.buy, a.buy, b.buy)
    if gene.buy:
        intersection(gene.sell, a.sell, b.sell)
        gene.conditions['fund_percent'] = (a.conditions['fund_percent'] + b.conditions['fund_percent']) * 0.5
        # gene.conditions['min_sell'] = (a.conditions['min_sell'] + b.conditions['min_sell']) * 0.5
        permutations.append(gene)

    gene = Genetics()
    union(gene.buy, a.buy, b.buy)
    union(gene.sell, a.sell, b.sell)
    gene.conditions['fund_percent'] = (a.conditions['fund_percent'] + b.conditions['fund_percent']) * 0.5
    # gene.conditions['min_sell'] = (a.conditions['min_sell'] + b.conditions['min_sell']) * 0.5
    permutations.append(gene)

    return permutations


def mutate(a):
    mutations = []

    gene = Genetics()
    copy_criteria(gene.buy, a.buy)
    copy_criteria(gene.sell, a.sell)
    gene.conditions['fund_percent'] = min(1.0, a.conditions['fund_percent'] + 0.01)
    # gene.conditions['min_sell'] = a.conditions['min_sell']
    mutations.append(gene)

    gene = Genetics()
    copy_criteria(gene.buy, a.buy)
    copy_criteria(gene.sell, a.sell)
    gene.conditions['fund_percent'] = max(0.2, a.conditions['fund_percent'] - 0.01)
    # gene.conditions['min_sell'] = a.conditions['min_sell']
    mutations.append(gene)

    gene = Genetics()
    copy_criteria(gene.buy, a.buy)
    copy_criteria(gene.sell, a.sell)
    gene.conditions['fund_percent'] = a.conditions['fund_percent']
    # gene.conditions['min_sell'] = min(1.1, a.conditions['min_sell'] + 0.01)
    mutations.append(gene)

    gene = Genetics()
    copy_criteria(gene.buy, a.buy)
    copy_criteria(gene.sell, a.sell)
    gene.conditions['fund_percent'] = a.conditions['fund_percent']
    # gene.conditions['min_sell'] = max(-1.0, a.conditions['min_sell'] - 0.01)
    mutations.append(gene)

    return mutations


class Genetics:
    def __init__(self):
        self.buy = {}
        self.sell = {}
        self.conditions = {}
        self.conditions['min_sell'] = 1.006

    def randomize(self):
        random_criteria(self.buy)
        random_criteria(self.sell)
        self.conditions['fund_percent'] = 0.2 + random.random() * 0.8
        # self.conditions['min_sell'] = 1.1 * random.random() - 1.0

    def signal(self, in_candles, in_index):
        global candles
        global index
        candles = in_candles
        index = in_index
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
