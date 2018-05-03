import patterns


class Strategy:
    def __init__(self, name, percent):
        self.name = name
        self.buy = []
        self.percent = percent

    def algorithm(self, candles, index):
        for criteria in self.buy:
            if not criteria(candles, index):
                return False
        return True

    def update_stop_limit(self, order, ticker):
        low = ticker * 0.95
        if low > order.stop_limit:
            order.stop_limit = low


def green_maru(candles, index):
    maru = patterns.marubozu(candles[index])
    return maru == 'green'