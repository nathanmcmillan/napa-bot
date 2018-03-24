
class MovingAverage:
    
    def __init__(self, periods, initial):
        self.periods = periods
        self.weight = 2.0 / (float(periods) + 1.0)
        self.current = initial
        
        
    def update(self, value):
        self.current = (value - self.current) * self.weight + self.current