package clients

import (
	"context"
	"fmt"

	"github.com/cloudradar-monitoring/rport/share/logger"
)

// LoadInitialClients returns an initial Client Repository state populated with clients from the internal storage.
func LoadInitialClients(ctx context.Context, p ClientStore, logger *logger.Logger) ([]*Client, error) {
	logger.Debugf("loading existing clients")

	// setup a logger for the clients
	clientLogger := logger.Fork("client")

	all, err := p.GetAll(ctx, clientLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %v", err)
	}

	logger.Debugf("loaded %d clients", len(all))

	// mark previously connected clients as disconnected with current time
	now := now()

	for _, client := range all {
		if client.IsConnected() {
			client.SetDisconnectedAt(&now)
			err := p.Save(ctx, client)
			if err != nil {
				return nil, fmt.Errorf("failed to save client: %v", err)
			}
		}
	}

	return all, nil
}
