package permissions

import (
	"context"
	"errors"
)

type key int

const Key key = 6789

type types string

const ToRead types = "read"
const ToWrite types = "write"
const ToDelete types = "delete"

func Check(ctx context.Context, identifier string, need types) error {
	permissions := map[string]int{
		"read":   1,
		"write":  2,
		"delete": 4,
	}

	permissionInt := permissions[string(need)]

	permissionLevel := ctx.Value(Key).(map[string]*int)[identifier]
	if permissionLevel == nil || *permissionLevel^permissionInt == permissionInt {
		return errors.New("usuário não possui autorização necessária")
	}

	return nil
}
