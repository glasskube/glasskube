package ctrlpkg

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func IsSameResource(a, b Package) bool {
	return a.GetName() == b.GetName() &&
		a.GroupVersionKind() == b.GroupVersionKind() &&
		a.GetNamespace() == b.GetNamespace()
}

func HasSpecChanged(pkg Package) (bool, string, error) {
	if specBytes, err := json.Marshal(pkg.GetSpec()); err != nil {
		return false, "", fmt.Errorf("failed to marshal package spec: %w", err)
	} else {
		var currentSpecHash string
		h := sha256.New()
		if _, err := h.Write(specBytes); err != nil {
			return false, "", fmt.Errorf("failed to hash package spec: %w", err)
		} else {
			currentSpecHash = hex.EncodeToString(h.Sum(nil))
		}
		return pkg.GetStatus().PreviousSpec != currentSpecHash, currentSpecHash, nil
	}
}
