package framework

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
)

// Secret is a type of secret that can be returned from a backend.
type Secret struct {
	// Type is the name of this secret type. This is used to setup the
	// vault ID and to look up the proper secret structure when revocation/
	// renewal happens. Once this is set this should not be changed.
	//
	// The format of this must match (case insensitive): ^a-Z0-9_$
	Type string

	// Fields is the mapping of data fields and schema that comprise
	// the structure of this secret.
	Fields map[string]*FieldSchema

	// DefaultDuration and DefaultGracePeriod are the default values for
	// the duration of the lease for this secret and its grace period. These
	// can be manually overwritten with the result of Response().
	DefaultDuration    time.Duration
	DefaultGracePeriod time.Duration

	// Below are the operations that can be called on the secret.
	//
	// Renew, if not set, will mark the secret as not renewable.
	//
	// Revoke is required.
	Renew  OperationFunc
	Revoke OperationFunc
}

// SecretType is the type of the secret with the given ID.
func SecretType(id string) (string, string) {
	idx := strings.Index(id, "-")
	if idx < 0 {
		return "", id
	}

	return id[:idx], id[idx+1:]
}

func (s *Secret) Response(
	data map[string]interface{}) (*logical.Response, error) {
	uuid, err := logical.UUID()
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s-%s", s.Type, uuid)
	return &logical.Response{
		IsSecret: true,
		Lease: &logical.Lease{
			VaultID:     id,
			Renewable:   s.Renew != nil,
			Duration:    s.DefaultDuration,
			GracePeriod: s.DefaultGracePeriod,
		},
		Data: data,
	}, nil
}