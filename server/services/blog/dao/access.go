package dao

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/database"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
	"github.com/ipuppet/aqua/server/services/blog/internal/netutil"
)

// LogAccessInfo records an article access with IPv4/IPv6 support.
func (r *repository) LogAccessInfo(ctx context.Context, aid int, ip string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	parsedIP, flag := netutil.ParseIP(ip)
	var sqlQuery string
	var err error

	switch flag {
	case 4:
		sqlQuery = "insert into blog.access (ipv4,aid) values ($1,$2)"
		ipInt := binary.BigEndian.Uint32(parsedIP)
		_, err = r.db.ExecContext(ctx, sqlQuery, ipInt, aid)
	case 6:
		sqlQuery = "insert into blog.access (ipv6,aid) values ($1,$2)"
		_, err = r.db.ExecContext(ctx, sqlQuery, parsedIP, aid)
	default:
		return fmt.Errorf("invalid IP address format for logging access: %s", ip)
	}

	if err != nil {
		return fmt.Errorf("logging access for aid %d, ip %s: %w", aid, ip, err)
	}
	return nil
}

// GetAccess returns article access records for the last N days.
func (r *repository) GetAccess(ctx context.Context, day int) ([]blog.Access, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	accesses := []blog.Access{}
	err := database.Retrieve(
		ctx,
		r.db,
		&accesses,
		`select access.*, article.title, EXTRACT(EPOCH FROM access.timestamp) as timestamp_epoch
        from blog.access
        left join blog.article on access.aid=article.aid
        where $1::date - ($2 || ' DAY')::INTERVAL <= DATE(access.timestamp)
        order by access.timestamp DESC`,
		time.Now().Format("2006-01-02"), day)
	if err != nil {
		return nil, fmt.Errorf("getting access data for last %d days: %w", day, err)
	}
	return accesses, nil
}
