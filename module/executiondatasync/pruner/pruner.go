package pruner

import (
	"fmt"

	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/executiondatasync/tracker"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/module/util"
)

const (
	defaultHeightRangeTarget = uint64(400000)
	defaultThreshold         = uint64(100000)
)

type Pruner struct {
	storage *tracker.Storage

	fulfilledHeightsIn    chan<- interface{}
	fulfilledHeightsOut   <-chan interface{}
	thresholdChan         chan uint64
	heightRangeTargetChan chan uint64

	lastPrunedHeight  uint64
	heightRangeTarget uint64
	threshold         uint64

	component.Component
	cm *component.ComponentManager
}

type PrunerOption func(*Pruner)

func WithHeightRangeTarget(heightRangeTarget uint64) PrunerOption {
	return func(p *Pruner) {
		p.heightRangeTarget = heightRangeTarget
	}
}

func WithThreshold(threshold uint64) PrunerOption {
	return func(p *Pruner) {
		p.threshold = threshold
	}
}

func NewPruner(storage *tracker.Storage, opts ...PrunerOption) (*Pruner, error) {
	lastPrunedHeight, err := storage.GetPrunedHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get pruned height: %w", err)
	}

	fulfilledHeight, err := storage.GetFulfilledHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to get fulfilled height: %w", err)
	}

	fulfilledHeightsIn, fulfilledHeightsOut := util.UnboundedChannel()
	fulfilledHeightsIn <- fulfilledHeight

	p := &Pruner{
		storage:               storage,
		fulfilledHeightsIn:    fulfilledHeightsIn,
		fulfilledHeightsOut:   fulfilledHeightsOut,
		thresholdChan:         make(chan uint64),
		heightRangeTargetChan: make(chan uint64),
		lastPrunedHeight:      lastPrunedHeight,
		heightRangeTarget:     defaultHeightRangeTarget,
		threshold:             defaultThreshold,
	}
	p.cm = component.NewComponentManagerBuilder().
		AddWorker(p.loop).
		Build()
	p.Component = p.cm

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

func (p *Pruner) NotifyFulfilledHeight(height uint64) {
	if util.CheckClosed(p.cm.ShutdownSignal()) {
		return
	}

	p.fulfilledHeightsIn <- height
}

func (p *Pruner) SetHeightRangeTarget(heightRangeTarget uint64) error {
	select {
	case p.heightRangeTargetChan <- heightRangeTarget:
		return nil
	case <-p.cm.ShutdownSignal():
		return component.ErrComponentShutdown
	}
}

func (p *Pruner) SetThreshold(threshold uint64) error {
	select {
	case p.thresholdChan <- threshold:
		return nil
	case <-p.cm.ShutdownSignal():
		return component.ErrComponentShutdown
	}
}

func (p *Pruner) loop(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
	ready()

	for {
		select {
		case <-ctx.Done():
			return
		case h := <-p.fulfilledHeightsOut:
			fulfilledHeight := h.(uint64)
			if fulfilledHeight-p.lastPrunedHeight > p.heightRangeTarget+p.threshold {
				pruneHeight := fulfilledHeight - p.heightRangeTarget

				if err := p.storage.Prune(pruneHeight); err != nil {
					ctx.Throw(fmt.Errorf("failed to prune: %w", err))
				}

				p.lastPrunedHeight = pruneHeight
			}
		case heightRangeTarget := <-p.heightRangeTargetChan:
			p.heightRangeTarget = heightRangeTarget
		case threshold := <-p.thresholdChan:
			p.threshold = threshold
		}
	}
}
