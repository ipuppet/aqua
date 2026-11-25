package grpcserver

//go:generate go tool goverter gen .

import (
	blogv1 "github.com/ipuppet/aqua/server/gen/go/aqua/blog/v1"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:ToPBTime
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:IntToInt32
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:Int32ToInt
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	ToPBArticle             func(item *blog.Article) *blogv1.Article
	ToPBArticles            func(items []blog.Article) []*blogv1.Article
	ToPBArticleDescriptions func(items []blog.ArticleDescription) []*blogv1.ArticleDescription

	ToPBArticleContext func(item *blog.ArticleContext) *blogv1.ArticleContext

	// goverter:map ID Id
	ToPBTag  func(tag blog.Tag) *blogv1.Tag
	ToPBTags func(tags []blog.Tag) []*blogv1.Tag
	// goverter:map Id ID
	ToModelTag  func(tag *blogv1.Tag) blog.Tag
	ToModelTags func(tags []*blogv1.Tag) []blog.Tag

	// goverter:map ID Id
	// goverter:map UID Uid
	// goverter:map IP Ip
	ToPBComment  func(item blog.Comment) *blogv1.Comment
	ToPBComments func(items []blog.Comment) []*blogv1.Comment

	// goverter:map ID Id
	// goverter:map IPv4 Ipv4
	// goverter:map IPv6 Ipv6
	ToPBAccessRecord func(item blog.AccessResponse) *blogv1.AccessRecord
	ToPBAccess       func(items []blog.AccessResponse) []*blogv1.AccessRecord
)
