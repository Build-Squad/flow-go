package approvals

import (
	"sync"

	"github.com/onflow/flow-go/model/flow"
)

// AggregatedSignatures is an utility struct that provides concurrency safe access
// to map of aggregated signatures indexed by chunk index
type AggregatedSignatures struct {
	signatures     map[uint64]flow.AggregatedSignature // aggregated signature for each chunk
	lock           sync.RWMutex                        // lock for modifying aggregatedSignatures
	numberOfChunks uint64
}

func NewAggregatedSignatures(chunks uint64) *AggregatedSignatures {
	return &AggregatedSignatures{
		signatures:     make(map[uint64]flow.AggregatedSignature, chunks),
		lock:           sync.RWMutex{},
		numberOfChunks: chunks,
	}
}

// PutSignature adds the AggregatedSignature from the collector to `aggregatedSignatures`.
// The returned int is the resulting number of approved chunks.
func (as *AggregatedSignatures) PutSignature(chunkIndex uint64, aggregatedSignature flow.AggregatedSignature) uint64 {
	as.lock.Lock()
	defer as.lock.Unlock()
	if chunkIndex < as.numberOfChunks {
		if _, found := as.signatures[chunkIndex]; !found {
			as.signatures[chunkIndex] = aggregatedSignature
		}
	}
	return uint64(len(as.signatures))
}

// HasSignature returns boolean depending if we have signature for particular chunk
func (as *AggregatedSignatures) HasSignature(chunkIndex uint64) bool {
	as.lock.RLock()
	defer as.lock.RUnlock()
	_, found := as.signatures[chunkIndex]
	return found
}

// Collect returns array with aggregated signature for each chunk
func (as *AggregatedSignatures) Collect() []flow.AggregatedSignature {
	as.lock.RLock()
	defer as.lock.RUnlock()

	// if there aren't signatures for each chunk we can't collect signatures
	if uint64(len(as.signatures)) < as.numberOfChunks {
		return nil
	}

	aggregatedSigs := make([]flow.AggregatedSignature, as.numberOfChunks)
	for chunkIndex, sig := range as.signatures {
		aggregatedSigs[chunkIndex] = sig
	}

	return aggregatedSigs
}

// ChunksWithoutAggregatedSignature returns indexes of chunks that don't have an aggregated signature
func (as *AggregatedSignatures) ChunksWithoutAggregatedSignature() []uint64 {
	// provide enough capacity to avoid allocations while we hold the lock
	missingChunks := make([]uint64, 0, as.numberOfChunks)
	as.lock.RLock()
	defer as.lock.RUnlock()
	for i := uint64(0); i < as.numberOfChunks; i++ {
		chunkIndex := uint64(i)
		if _, found := as.signatures[chunkIndex]; found {
			// skip if we already have enough valid approvals for this chunk
			continue
		}
		missingChunks = append(missingChunks, chunkIndex)
	}
	return missingChunks
}
