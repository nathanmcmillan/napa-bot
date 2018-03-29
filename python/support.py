fib38 = 0.382
fib50 = 0.500
fib62 = 0.618


class Retracement:
    def __init__(self, candles):
        self.signal = "wait"
        if len(candles) == 0:
            return
        self.high = candles[0].high
        self.low = candles[0].low
        for candle in candles:
            if candle.high > self.high:
                self.high = candle.high
            if candle.low > self.low:
                self.low = candle.low
        self.line38 = self.high - (self.high - self.low) * fib38
        self.line50 = self.high - (self.high - self.low) * fib50
        self.line68 = self.high - (self.high - self.low) * fib62