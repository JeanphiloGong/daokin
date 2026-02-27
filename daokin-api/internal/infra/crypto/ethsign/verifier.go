package ethsign

import (
	"context"
	"strings"

	"daokin-api/internal/domain"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type PersonalSignVerifier struct{}

func NewPersonalSignVerifier() *PersonalSignVerifier {
	return &PersonalSignVerifier{}
}

func (v *PersonalSignVerifier) VerifyWalletSignature(_ context.Context, wallet, message, signature string) error {
	wallet = strings.TrimSpace(wallet)
	message = strings.TrimSpace(message)
	signature = strings.TrimSpace(signature)

	if wallet == "" || message == "" || signature == "" {
		return domain.ErrInvalidArgument
	}
	if !common.IsHexAddress(wallet) {
		return domain.ErrInvalidArgument
	}

	sig, err := hexutil.Decode(signature)
	if err != nil {
		return domain.ErrUnauthorized
	}
	if len(sig) != 65 {
		return domain.ErrUnauthorized
	}

	// Wallet providers often return V as 27/28, but crypto.SigToPub expects 0/1.
	sigCopy := append([]byte(nil), sig...)
	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}
	if sigCopy[64] > 1 {
		return domain.ErrUnauthorized
	}

	hash := accounts.TextHash([]byte(message))
	pubKey, err := crypto.SigToPub(hash, sigCopy)
	if err != nil {
		return domain.ErrUnauthorized
	}

	recovered := crypto.PubkeyToAddress(*pubKey)
	if !strings.EqualFold(recovered.Hex(), common.HexToAddress(wallet).Hex()) {
		return domain.ErrUnauthorized
	}
	return nil
}
