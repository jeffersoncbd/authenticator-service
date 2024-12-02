package utils

import (
	"context"
	"errors"
	"slices"
)

type key int

const ContextKey key = 6789

type types string

const KeyToRead types = "read"
const KeyToWrite types = "write"
const KeyToDelete types = "delete"

func CheckPermissions(ctx context.Context, identifier string, need types) error {
	permissions := map[string]int{
		"read":   1,
		"write":  2,
		"delete": 4,
	}

	permissionInt := permissions[string(need)]

	userPermissions := ctx.Value(ContextKey).(map[string]*int)

	keys := make([]string, 0, len(userPermissions))
	for k := range userPermissions {
			keys = append(keys, k)
	}

	if !slices.Contains(keys, identifier) {
		return errors.New("usuário não possui a autorização necessária")
	}

	permissionLevel := userPermissions[identifier]

	if permissionLevel == nil || *permissionLevel^permissionInt == permissionInt {
		return errors.New("usuário não possui autorização necessária")
	}

	return nil
}
