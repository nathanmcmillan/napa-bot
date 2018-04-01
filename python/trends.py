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