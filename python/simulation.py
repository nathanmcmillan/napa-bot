class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price
        self.stop_limit = 0.0
        self.high = coin_price
        self.low = coin_price
        self.draw_down = 0.0

    def update(self, ticker):
        if ticker < self.low:
            self.low = ticker
        elif ticker > self.high:
            self.high = ticker
        self.draw_down = (self.high - self.low) / self.high


def run(candles, intervals, funds, fees, strategy, print_trades):
    candle_count = len(candles)
    orders = []
    low = funds
    high = funds
    coins = 0.0
    buys = 0
    sells = 0
    gains = 0
    losses = 0
    draw_down = 0.0
    index = intervals
    while index < candle_count:
        ticker_price = candles[index].closing

        for order in orders[:]:
            if strategy.sell(candles, index, order) or ticker_price < order.stop_limit:
                orders.remove(order)
                usd = (ticker_price * order.size) * (1.0 - fees)
                funds += usd
                coins -= order.size
                sells += 1
                total = funds + coins * ticker_price
                if total > high:
                    high = total
                profit = usd - order.usd * (1.0 + fees)
                if profit > 0:
                    gains += 1
                else:
                    losses += 1
                if print_trades:
                    print('time - {} - ticker ${:,.2f} - profit ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, profit, funds, coins))
            else:
                strategy.stop_limit(candles, index, order)
                order.update(ticker_price)
                if order.draw_down > draw_down:
                    draw_down = order.draw_down

        if strategy.buy(candles, index):
            usd = strategy.amount(funds)
            if usd > 10.0:
                order = SimOrder(ticker_price, None, usd)
                orders.append(order)
                strategy.stop_limit(candles, index, order)
                usd *= (1.0 + fees)
                funds -= usd
                coins += orders[-1].size
                buys += 1
                total = funds + coins * ticker_price
                if total < low:
                    low = total
                if print_trades:
                    print('time - {} - ticker ${:,.2f} - spent ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, usd, funds, coins))

        index += 1

    total = 0.0
    coins = 0.0
    end_price = candles[-1].closing
    for order in orders:
        total += order.size * end_price
        coins += order.size
    total += funds
    print('total ${:,.2f} - coins {:,.3f}'.format(total, coins))
    return [total, coins, low, high, buys, sells, gains, losses, draw_down]


def perfect(candles, funds):
    first_low = candles[0].closing
    high = candles[0].closing
    second_low = None
    for candle in candles:
        if second_low:
            if candle.closing < second_low:
                second_low = candle.closing
            elif candle.high > high:
                high = candle.closing
            if (high - second_low) / second_low > 0.01:
                funds *= (high - first_low) / first_low
        else:
            if candle.closing < first_low:
                first_low = candle.closing
            elif candle.high > high:
                high = candle.closing
            if (high - first_low) / first_low > 0.01:
                second_low = high
    return funds