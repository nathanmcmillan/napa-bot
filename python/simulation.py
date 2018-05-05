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
    index = intervals
    while index < candle_count:
        ticker_price = candles[index].closing

        sell = strategy.sell(candles, index)
        for order in orders[:]:
            if sell or ticker_price < order.stop_limit:
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
    return [total, coins, low, high, buys, sells, gains, losses]
