package domain

import (
	"math/big"
	"strings"
	"time"
)

type Transfer struct {
	ID              string    `json:"id"`
	ArtifactID      string    `json:"artifact_id"`
	PayerWallet     string    `json:"payer_wallet"`
	RecipientWallet string    `json:"recipient_wallet"`
	Token           string    `json:"token"`
	AmountAtomic    string    `json:"amount_atomic"`
	TxHash          string    `json:"tx_hash"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewTransfer(
	artifactID string,
	payerWallet string,
	recipientWallet string,
	token string,
	amountAtomic string,
	txHash string,
	createdAt time.Time,
) (Transfer, error) {
	artifactID = strings.TrimSpace(artifactID)
	payerWallet = strings.TrimSpace(payerWallet)
	recipientWallet = strings.TrimSpace(recipientWallet)
	token = strings.TrimSpace(token)
	amountAtomic = strings.TrimSpace(amountAtomic)
	txHash = strings.TrimSpace(txHash)
	if artifactID == "" || payerWallet == "" || recipientWallet == "" || token == "" || amountAtomic == "" || txHash == "" {
		return Transfer{}, ErrInvalidArgument
	}

	v, ok := new(big.Int).SetString(amountAtomic, 10)
	if !ok || v.Sign() <= 0 {
		return Transfer{}, ErrInvalidArgument
	}

	return Transfer{
		ArtifactID:      artifactID,
		PayerWallet:     payerWallet,
		RecipientWallet: recipientWallet,
		Token:           token,
		AmountAtomic:    amountAtomic,
		TxHash:          txHash,
		CreatedAt:       createdAt.UTC(),
	}, nil
}
