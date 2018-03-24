from ema import MovingAverage

class ConvergeDiverge:
    
    def __init__(self, short, long, initial):
        self.short = MovingAverage(short, initial)
        self.long = MovingAverage(long, initial)
        self.current = 0
        self.signal = 'wait'
        
        
    def update(self, value):
        self.short.update(value)
        self.long.update(value)
        before = self.current
        self.current = self.short.current - self.long.current
        if before < 0 and self.current > 0:
            self.signal = "buy"
        elif before > 0 and self.current < 0:
            self.signal = "sell"
        else:
            self.signal = "wait"