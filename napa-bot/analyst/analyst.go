package analyst

import (
    "fmt"
    "../gdax"
)

func RelativeStrengthIndex(history []gdax.Candle) (float64) {
    for i := 1; i < len(history); i++ {
        prev := history[i - 1]
        now := history[i]
        diff := now.Closing - prev.Closing
        fmt.Println("closing ", now.Closing, " | ", diff)
    }
    rs := float64(50)
    return 100.0 / (100.0 - (1.0 + rs))
}

func SimpleMovingAverage() {
    fmt.Println("cutler non exponential is not dependent on length or starting point")   
}