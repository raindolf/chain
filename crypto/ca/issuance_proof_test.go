package ca

import (
	"testing"

	"chain/crypto/ed25519/ecmath"
)

func TestIssuanceProof(t *testing.T) {
	var (
		y        [3]ecmath.Scalar // issuance private keys
		assetIDs [3]AssetID       // issuance asset ids
		Y        [3]ecmath.Point  // issuance public keys
	)
	for i := 0; i < len(assetIDs); i++ {
		y[i][0] = byte(10 + i)
		assetIDs[i][0] = byte(20 + i)
		Y[i].ScMulBase(&y[i])
	}
	candidates := []AssetIssuanceCandidate{
		&testAssetIssuanceCandidate{
			assetID:     &assetIDs[0],
			issuanceKey: &Y[0],
		},
		&testAssetIssuanceCandidate{
			assetID:     &assetIDs[1],
			issuanceKey: &Y[1],
		},
		&testAssetIssuanceCandidate{
			assetID:     &assetIDs[2],
			issuanceKey: &Y[2],
		},
	}

	var nonce [32]byte
	copy(nonce[:], []byte("nonce"))
	msg := []byte("message")

	j := uint64(1) // secret index
	aek := []byte("asset encryption key")
	ac, c := CreateAssetCommitment(assetIDs[j], aek)
	iarp := CreateConfidentialIARP(ac, *c, candidates, nonce, msg, j, y[j])
	ip := CreateIssuanceProof(ac, iarp, candidates, msg, nonce, y[j])
	valid, yj := ip.Validate(ac, iarp, candidates, msg, nonce, j)
	if !valid {
		t.Error("failed to validate issuance proof")
	}
	if !yj {
		t.Error("validated issuance proof but not yj")
	}

	if valid, _ = ip.Validate(ac, iarp, candidates, msg[1:], nonce, j); valid {
		t.Error("validated invalid issuance proof")
	}
	if valid, _ = ip.Validate(ac, iarp, candidates, msg, nonce, 0); valid {
		t.Error("validated invalid issuance proof")
	}

	nonce2 := nonce
	nonce2[0] ^= 1
	if valid, _ = ip.Validate(ac, iarp, candidates, msg, nonce2, j); valid {
		t.Error("validated invalid issuance proof")
	}

	ip2 := *ip
	ip2.e1[0] ^= 1
	if valid, _ = ip2.Validate(ac, iarp, candidates, msg, nonce, j); valid {
		t.Error("validated invalid issuance proof")
	}
	ip2 = *ip
	ip2.s1[0] ^= 1
	if valid, _ = ip2.Validate(ac, iarp, candidates, msg, nonce, j); valid {
		t.Error("validated invalid issuance proof")
	}
	ip2 = *ip
	ip2.e2[0] ^= 1
	if valid, _ = ip2.Validate(ac, iarp, candidates, msg, nonce, j); valid {
		t.Error("validated invalid issuance proof")
	}
	ip2 = *ip
	ip2.s2[0] ^= 1
	if valid, _ = ip2.Validate(ac, iarp, candidates, msg, nonce, j); valid {
		t.Error("validated invalid issuance proof")
	}
}