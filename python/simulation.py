class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price


def round(candles, intervals, funds, fees, algorithm, conditions, print_trades):
    candle_count = len(candles)
    orders = []
    limits = []
    low = funds
    high = funds
    coins = 0.0
    buys = 0
    sells = 0
    index = intervals
    while index < candle_count:
        ticker_price = candles[index].closing
        for limit_order in limits[:]:
            if limit_order.is_buy:
                if ticker_price < limit_order.coin_price:
                    orders.append(SimOrder(limit_order.coin_price, None, limit_order.usd))    
            else:
                if ticker_price > limit_order.price:
                    orders.append(SimOrder(limit_order.coin_price, None, limit_order.usd))
        signal = algorithm(candles, index)
        if signal == 'buy':
            usd = funds * conditions['fund_percent']
            if usd > 10.0:
                orders.append(SimOrder(ticker_price, None, usd))
                usd *= (1.0 + fees)
                funds -= usd
                coins += orders[-1].size
                buys += 1
                total = funds + coins * ticker_price
                if total < low:
                    low = total
                if print_trades:
                    print('time - {} - ticker ${:,.2f} - spent ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, usd, funds, coins))
        elif signal == 'sell':
            for order_to_sell in orders[:]:
                change = (ticker_price - order_to_sell.coin_price) / order_to_sell.coin_price
                if change > conditions['min_sell']:
                    orders.remove(order_to_sell)
                    usd = (ticker_price * order_to_sell.size) * (1.0 - fees)
                    funds += usd
                    coins -= order_to_sell.size
                    sells += 1
                    total = funds + coins * ticker_price
                    if total > high:
                        high = total
                    if print_trades:
                        profit = usd - order_to_sell.usd * (1.0 + fees)
                        print('time - {} - ticker ${:,.2f} - profit ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, profit, funds, coins))
        index += 1
    total = 0.0
    coins = 0.0
    end_price = candles[-1].closing
    for order in orders:
        total += order.size * end_price
        coins += order.size
    total += funds
    print('total ${:,.2f} - coins {:,.3f}'.format(total, coins))
    return [total, coins, low, high, buys, sells]
