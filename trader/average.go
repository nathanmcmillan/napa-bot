package trader

// MovingAverage average move value of move product
type MovingAverage struct {
	Current float64
	Limit   int
	List    []float64
}

// NewMovingAverage constructor
func NewMovingAverage(limit int) *MovingAverage {
	move := MovingAverage{}
	move.Limit = limit
	move.List = make([]float64, 0)
	return &move
}

// Rolling new average using queue
func (move *MovingAverage) Rolling(value float64) {
	num := len(move.List)
	count := float64(num)
	if num == move.Limit {
		move.Current = ((move.Current * count) - move.List[0]) / (count - 1.0)
		move.List = move.List[1:]
		count -= 1.0
	}
	move.Current = (move.Current*count + value) / (count + 1.0)
	move.List = append(move.List, value)
}
