package logic

//go:generate go tool goverter gen .

import (
	"aqua/pkg/types"
	blog "aqua/services/blog/domain"
)

func nullStringToString(value types.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend aqua/services/blog/logic:nullStringToString
// goverter:extend aqua/pkg/utils/convert:CloneTime
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	ToAccessResponse func(item *blog.Access) *blog.AccessResponse
)
