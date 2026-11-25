package logic

//go:generate go tool goverter gen .

import (
	"github.com/ipuppet/aqua/server/pkg/types"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

func nullStringToString(value types.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend github.com/ipuppet/aqua/server/services/blog/logic:nullStringToString
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:CloneTime
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	ToAccessResponse func(item *blog.Access) *blog.AccessResponse
)
