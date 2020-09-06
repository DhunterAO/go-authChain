package crypto

import (
	"crypto/ecdsa"
	"github.com/DhunterAO/goAuthChain/crypto/types"
)

func CalcHash(data []byte) *types.Hash {
	hash := types.Keccak256Hash(data)
	return &hash
}

func recoverEcdsaFromSig(hash *types.Hash, sig *types.Signature) (*ecdsa.PublicKey, error) {
	uncompressedPk, err := types.Ecrecover(hash.Bytes(), sig.Bytes())
	if err != nil {
		return nil, err
	}
	pk, err := types.UnmarshalPubkey(uncompressedPk)
	if err != nil {
		return nil, err
	}
	return pk, err
}

func RecoverPkFromSig(hash *types.Hash, sig *types.Signature) (*types.Pubkey, error) {
	pk, err := recoverEcdsaFromSig(hash, sig)
	if err != nil {
		return &types.Pubkey{}, err
	}
	return types.BytesToPubkey(types.CompressPubkey(pk)), err
}

func PubkeyToAddress(pk *types.Pubkey) (*types.Address, error) {
	ecdsaPk, err := types.DecompressPubkey(pk.Bytes())
	if err != nil {
		return nil, err
	}
	addr := types.EcdsaPubkeyToAddress(*ecdsaPk)
	return addr, nil
}

func Sign(hash *types.Hash, prv *ecdsa.PrivateKey) (*types.Signature, error) {
	sig, err := types.Sign(hash.Bytes(), prv)
	return types.BytesToSignature(sig), err
}

func VerifySignature(pubkey *types.Pubkey, hash *types.Hash, sig *types.Signature) bool {
	return types.VerifySignature(pubkey[:], hash[:], sig[:types.SignatureLength-1])
}