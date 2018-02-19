
package analyst

import (
    "../historian"
)

func SimpleMovingAverage(periods int, history []historian.Period) ([]float64) {
    size := len(history)
    sma := make([]float64, size)
    for i := 0; i < size; i++ {
        if i < periods {
            sma[i] = history[i].Closing
            continue
        }
        sum := float64(0.0)
        for j := i - periods; j < i; j++ {
            sum += history[j].Closing
        }
        sma[i] = sum / float64(periods)
    }
    return sma
}

func ExponentialMovingAverage(periods int, history []historian.Period) ([]float64) {
    size := len(history)
    ema := make([]float64, size)
    weight := 2.0 / (float64(periods) + 1.0)
    for i := 0; i < size; i++ {
        if i < periods {
            ema[i] = history[i].Closing
            continue
        }
        previous := ema[i - 1]
        ema[i] = (history[i].Closing - previous) * weight + previous
    }
    return ema
}

func MovingAverageConvergenceDivergence(periodsA int, periodsB int, history []historian.Period) ([]float64) {
    emaA := ExponentialMovingAverage(periodsA, history)
    emaB := ExponentialMovingAverage(periodsB, history)
    size := len(history)
    macd := make([]float64, size)
    for i := 0; i < size; i++ {
        macd[i] = emaA[i] - emaB[i]
    }
    return macd
}

func RelativeStrengthIndex(periods int, history []historian.Period) ([]float64) {
    size := len(history)
    u := make([]float64, size)
    d := make([]float64, size)
    rsi := make([]float64, size)
    u[0] = 0.0
    d[0] = 0.0
    rsi[0] = 0.0
    for i := 1; i < size; i++ {
        prev := history[i - 1].Closing
        now := history[i].Closing
        if now > prev {
            u[i] = now - prev
            d[i] = 0.0
        } else {
            u[i] = 0.0
            d[i] = prev - now
        }
        if i < periods {
            rsi[i] = 0.0
            continue
        }
        smaU := float64(0.0)
        smaD := float64(0.0)
        for j := i - periods; j < i; j++ {
            smaU += u[j]
            smaD += d[j]
        }
        smaU /= float64(periods)
        smaD /= float64(periods)
        var rs float64
        if smaU == 0.0 || smaD == 0.0 {
            rs = 0.0
        } else {
            rs = smaU / smaD
        }
        rsi[i] = 100.0 - (100.0 / (1.0 + rs)) 
    }
    return rsi
}