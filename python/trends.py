class MovingAverage:
    def __init__(self, periods, initial):
        self.periods = periods
        self.weight = 2.0 / (float(periods) + 1.0)
        self.current = initial

    def update(self, value):
        self.current = (value - self.current) * self.weight + self.current


class ConvergeDiverge:
    def __init__(self, ema_short, ema_long, closing):
        self.ema_short = MovingAverage(ema_short, closing)
        self.ema_long = MovingAverage(ema_long, closing)
        self.current = 0
        self.signal = 'wait'

    def update(self, closing):
        self.ema_short.update(closing)
        self.ema_long.update(closing)
        before = self.current
        self.current = self.ema_short.current - self.ema_long.current
        if before < 0 and self.current > 0:
            self.signal = "buy"
        elif before > 0 and self.current < 0:
            self.signal = "sell"
        else:
            self.signal = "wait"


class AverageDirectional:
    def __init__(self, periods):
        self.current = None
        self.periods = periods

    def update(self, candles):
        end = len(candles)
        start = end - self.periods
        positive_dm = 0.0
        negative_dm = 0.0
        average_range = 0.0
        for index in range(start, end):
            today = candles[index]
            yesterday = candles[index - 1]
            up_move = today.high - yesterday.high
            down_move = yesterday.low - today.low
            if up_move > down_move and up_move > 0.0:
                positive_dm += up_move
            if down_move > up_move and down_move > 0.0:
                negative_dm += down_move
            average_range += true_range(today)
        positive_di = positive_dm / average_range
        negative_di = negative_dm / average_range
        direction_movement = abs(positive_di - negative_di) / (positive_di + negative_di)
        self.current = direction_movement


def true_range(candle):
    return max(candle.high - candle.low, abs(candle.high - candle.open), abs(candle.low - candle.open))


def support(candles, start, end):
    MIN_DIFFERENCE = 0.01
    low = candles[start].closing
    count = 1
    for index in range(start, end):
        candle = candles[index]
        diff = abs(low - candle.closing) / candle.closing
        if diff <= MIN_DIFFERENCE:
            count += 1
        elif candle.closing < low:
            low = candle.closing
            count = 1
    if count > 2:
        return low
    return None


def resistance(candles, start, end):
    MIN_DIFFERENCE = 0.01
    high = candles[start].closing
    count = 1
    for index in range(start, end):
        candle = candles[index]
        diff = abs(high - candle.closing) / candle.closing
        if diff <= MIN_DIFFERENCE:
            count += 1
        elif candle.closing > high:
            high = candle.closing
            count = 1
    if count > 2:
        return high
    return None


def liquidation(candles, start, middle, end):
    first = candles[start].closing
    second = candles[middle].closing
    third = candles[end].closing
    if first < third:
        return False
    diff = abs(second - first) / first
    if diff > 0.01:
        return False
    diff = abs(third - second) / second
    return diff > 0.01


def calculate_support_resistance(candles):
    return []


def derive(candles, start, end):
    delta = 0.0
    success = True
    for index in range(start + 1, end):
        current = candles[index].closing - candles[index - 1].closing
        if current >= delta:
            delta = current
        else:
            success = False
            break
    if success:
        return 'green'
    delta = 0.0
    for index in range(start + 1, end):
        current = candles[index].closing - candles[index - 1].closing
        if current <= delta:
            delta = current
        else:
            return ''
    return 'red'