class Strategy:
    def __init__(self, buy, sell, percent):
        self.buy = buy
        self.sell = sell
        self.percent = percent
        
    def buy(self, candles, index):
        for _, criteria in self.buy.items():
            if not criteria.met(candles, index):
                return False
        return True

    def stop_limit(self, order, ticker):
        return 0.0