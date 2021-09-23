package memorystorage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	accFromID        = "Ivanov"
	accFromAmount    = 1000
	accToID          = "Petrov"
	accToAmount      = 0
	amountToTransfer = 1
)

func TestTransfer(t *testing.T) {
	t.Run("check for concurrent transfers", func(t *testing.T) {
		storage := New()
		err := storage.Set(accFromID, accFromAmount)
		require.NoErrorf(t, err, "storage.Set has got an err: ", err)
		err = storage.Set(accToID, accToAmount)
		require.NoErrorf(t, err, "storage.Set has got an err: ", err)
		sumBefore := accFromAmount + accToAmount
		storage.Print()
		wg := &sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = storage.Transfer(accFromID, accToID, amountToTransfer)
				require.NoErrorf(t, err, "storage.Transfer has got an err: ", err)
			}()
		}
		wg.Wait()
		storage.Print()
		sum1, err := storage.Get(accFromID)
		require.NoErrorf(t, err, "storage.Get has got an err: ", err)
		sum2, err := storage.Get(accToID)
		require.NoErrorf(t, err, "storage.Get has got an err: ", err)
		sumAfter := sum1 + sum2
		require.Equalf(t, sumBefore, int(sumAfter), "expected %v, have %v", sumBefore, sumAfter)
	})
}
