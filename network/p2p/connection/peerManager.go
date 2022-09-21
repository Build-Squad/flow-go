package connection

import (
	"context"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/network/p2p"
)

// DefaultPeerUpdateInterval is default duration for which the peer manager waits in between attempts to update peer connections
var DefaultPeerUpdateInterval = 10 * time.Minute

var _ p2p.PeerManager = (*PeerManager)(nil)
var _ component.Component = (*PeerManager)(nil)

// PeerManager adds and removes connections to peers periodically and on request
type PeerManager struct {
	component.Component

	logger             zerolog.Logger
	peersProvider      p2p.PeersProvider // callback to retrieve list of peers to connect to
	peerRequestQ       chan struct{}     // a channel to queue a peer update request
	connector          p2p.Connector     // connector to connect or disconnect from peers
	peerUpdateInterval time.Duration     // interval the peer manager runs on
}

// NewPeerManager creates a new peer manager which calls the peersProvider callback to get a list of peers to connect to
// and it uses the connector to actually connect or disconnect from peers.
func NewPeerManager(logger zerolog.Logger, updateInterval time.Duration, peersProvider p2p.PeersProvider, connector p2p.Connector) *PeerManager {
	pm := &PeerManager{
		logger:             logger,
		peersProvider:      peersProvider,
		connector:          connector,
		peerRequestQ:       make(chan struct{}, 1),
		peerUpdateInterval: updateInterval,
	}

	pm.Component = component.NewComponentManagerBuilder().
		AddWorker(func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
			// makes sure that peer update request is invoked once before returning
			pm.RequestPeerUpdate()

			ready()
			pm.updateLoop(ctx)
		}).
		AddWorker(func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
			ready()
			pm.periodicLoop(ctx)
		}).
		Build()

	return pm
}

// PeerManagerFactory generates a PeerManagerFunc that produces the default PeerManager with the given peer manager
// options and that uses the LibP2PConnector with the given LibP2P connector options
func PeerManagerFactory(connectionPruning bool, updateInterval time.Duration) p2p.PeerManagerFactoryFunc {
	return func(host host.Host, peersProvider p2p.PeersProvider, logger zerolog.Logger) (p2p.PeerManager, error) {
		connector, err := NewLibp2pConnector(logger, host, connectionPruning)
		if err != nil {
			return nil, fmt.Errorf("failed to create libp2pConnector: %w", err)
		}
		return NewPeerManager(logger, updateInterval, peersProvider, connector), nil
	}
}

// updateLoop triggers an update peer request when it has been requested
func (pm *PeerManager) updateLoop(ctx irrecoverable.SignalerContext) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-pm.peerRequestQ:
			pm.updatePeers(ctx)
		}
	}
}

// periodicLoop periodically triggers an update peer request
func (pm *PeerManager) periodicLoop(ctx irrecoverable.SignalerContext) {
	// add a random delay to initial launch to avoid synchronizing this
	// potentially expensive operation across the network
	delay := time.Duration(mrand.Int63n(pm.peerUpdateInterval.Nanoseconds()))

	ticker := time.NewTicker(pm.peerUpdateInterval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-time.After(delay):
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pm.RequestPeerUpdate()
		}
	}
}

// RequestPeerUpdate requests an update to the peer connections of this node.
// If a peer update has already been requested (either as a periodic request or an on-demand request) and is outstanding,
// then this call is a no-op.
func (pm *PeerManager) RequestPeerUpdate() {
	select {
	case pm.peerRequestQ <- struct{}{}:
	default:
	}
}

// updatePeers updates the peers by connecting to all the nodes provided by the peersProvider callback and disconnecting from
// previous nodes that are no longer in the new list of nodes.
func (pm *PeerManager) updatePeers(ctx context.Context) {

	// get all the peer ids to connect to
	peers := pm.peersProvider()

	pm.logger.Trace().
		Str("peers", fmt.Sprintf("%v", peers)).
		Msg("connecting to peers")

	// ask the connector to connect to all peers in the list
	pm.connector.UpdatePeers(ctx, peers)
}

// ForceUpdatePeers initiates an update to the peer connections of this node immediately
func (pm *PeerManager) ForceUpdatePeers(ctx context.Context) {
	pm.updatePeers(ctx)
}
