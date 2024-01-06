package querier

import (
	"context"
	"state-proof-relayer/internal/servicestate"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Querier struct {
	client *algod.Client
}

func InitializeQuerier(nodePath string) (*Querier, error) {
	var algodAddress = "https://testnet-api.algonode.cloud"

	client, err := algod.MakeClient(algodAddress, "")
	if err != nil {
		return nil, err
	}

	return &Querier{
		client: client,
	}, nil
}

func (q *Querier) QueryNextStateProofData(state *servicestate.ServiceState) (*models.StateProof, error) {
	proof, err := q.client.GetStateProof(state.LatestCompletedAttestedRound + 1).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return &proof, nil
}
