package writer

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-stateproof-verification/stateproof"

	"state-proof-relayer/internal/servicestate"
	"state-proof-relayer/internal/utilities"
)

type Writer struct {
}

func InitializeWriter() *Writer {
	return &Writer{}
}

func (w *Writer) UploadStateProof(state *servicestate.ServiceState, proof *models.StateProof) error {
	var stateProof stateproof.StateProof
	err := msgpack.Decode(proof.Stateproof, &stateProof)
	if err != nil {
		return err
	}

	objectName := fmt.Sprintf("proof_%d_to_%d", proof.Message.Firstattestedround, proof.Message.Lastattestedround)

	utilities.EncodeToFile(stateProof, "stateProof_"+objectName+".json")

	state.LatestCompletedAttestedRound = proof.Message.Lastattestedround

	return nil
}
