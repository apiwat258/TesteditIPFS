package ipfs

import (
    "github.com/ipfs/go-ipfs-api"
    "strings"
    "fmt"
)

// StoreInIPFS stores data in IPFS and returns the CID
func StoreInIPFS(document string) (string, error) {
    // Connect to the local IPFS node (replace with your IPFS node URL if needed)
    sh := shell.NewShell("localhost:5001")

    // Store the document on IPFS
    cid, err := sh.Add(strings.NewReader(document))
    if err != nil {
        return "", fmt.Errorf("failed to add document to IPFS: %w", err)
    }

    return cid, nil
}
