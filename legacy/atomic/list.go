package atomic

import (
	"math/big"
	"sync"
)

type List struct {
	m  *sync.Mutex
	ls []big.Rat
}
