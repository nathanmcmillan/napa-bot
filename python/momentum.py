from trends import MovingAverage


class MoneyFlow:
    def __init__(self, periods):
        self.signal = "wait"
        self.periods = periods
        self.current = None

    def update(self, candles):
        positive = 0.0
        negative = 0.0
        end = len(candles)
        start = end - self.periods
        previous = candles[start].typical_price()
        start += 1
        for index in range(start, end):
            candle = candles[index]
            typical_price = candle.typical_price()
            money_flow = typical_price * candle.volume
            if typical_price > previous:
                positive += money_flow
            elif typical_price < previous:
                negative += money_flow
            previous = typical_price
        self.current = positive / (positive + negative)
        if self.current >= 0.8:
            self.signal = "sell"
        elif self.current <= 0.2:
            self.signal = "buy"
        else:
            self.signal = "wait"


class RelativeStrength:
    def __init__(self, periods):
        self.signal = "wait"
        self.periods = periods
        self.current = None

    def update(self, candles):
        positive = MovingAverage(self.periods - 1, 0.0)
        negative = MovingAverage(self.periods - 1, 0.0)
        end = len(candles)
        start = end - self.periods + 1
        for index in range(start, end):
            prev = candles[index - 1].closing
            now = candles[index].closing
            if now > prev:
                positive.update(now - prev)
            else:
                negative.update(prev - now)
        self.current = positive.current / (positive.current + negative.current)
        if self.current >= 0.8:
            self.signal = "sell"
        elif self.current <= 0.2:
            self.signal = "buy"
        else:
            self.signal = "wait"