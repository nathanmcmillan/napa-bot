class MovingAverage:
    def __init__(self, periods, initial):
        self.periods = periods
        self.weight = 2.0 / (float(periods) + 1.0)
        self.current = initial

    def update(self, value):
        self.current = (value - self.current) * self.weight + self.current


class ConvergeDiverge:
    def __init__(self, short, long, closing):
        self.short = MovingAverage(short, closing)
        self.long = MovingAverage(long, closing)
        self.current = 0
        self.signal = 'wait'

    def update(self, closing):
        self.short.update(closing)
        self.long.update(closing)
        before = self.current
        self.current = self.short.current - self.long.current
        if before < 0 and self.current > 0:
            self.signal = "buy"
        elif before > 0 and self.current < 0:
            self.signal = "sell"
        else:
            self.signal = "wait"