package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis 사용 cmd 정리
// - DONE
// ZADD (messageId 최초 create 할때 호출)
// ZADDNX (번역결과 add 할때 "ko", "en" 필드에 각각 번역됐던게 있으면 추가로 set 할 필요 없어서 NX 사용)
// ZRANGEBYSCORE (room 안에서 생성된 message id 가져올때)
// HGETALL (hashes 에서 message 단건 정보 가져올때, 번역결과)
// HGET (번역결과 단건 필요할때)
// HGETALL Multi (timestamp 이후로 발생한 message 들을 hashes 에서 한번에 긁어올때 MGET 같이.. , rdb.pipelined 를 이용한다)
// EXPIRE (해당 message, room 관련 expire 를 설정 가능할듯)
// HSET (HMSET 은 deprecated 됨, HSET 에 values 를 map 으로 전달 하면 multi 로 동작하게됨)

// HSet

func (c *Client) Del(ctx context.Context, keys ...string) error {
	if _, err := c.cl.Del(ctx, keys...).Result(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	res, err := c.cl.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) HSet(ctx context.Context, key string, v ...string) error {
	if _, err := c.cl.HSet(ctx, key, v).Result(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Expire(ctx context.Context, key string, seconds int) error {
	cmd := c.cl.Expire(ctx, key, time.Duration(seconds)*time.Second)
	if _, err := cmd.Result(); err != nil {
		return err
	}
	return nil
}

func (c *Client) MultiHGetAll(ctx context.Context, keys []string) []map[string]string {
	res := make([]map[string]string, 0)
	c.cl.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, k := range keys {
			if r, err := p.HGetAll(ctx, k).Result(); err != nil {
				res = append(res, nil)
			} else {
				res = append(res, r)
			}
		}
		return nil
	})
	return res
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	res, err := c.cl.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	res, err := c.cl.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	res, err := c.cl.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Client) HDel(ctx context.Context, key, field string) error {
	_, err := c.cl.HDel(ctx, key, field).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ZRem(ctx context.Context, key string, members ...string) error {
	_, err := c.cl.ZRem(ctx, key, members).Result()
	if err != nil {
		return err
	}
	return nil
}

// ZRangeByScoreOver 는, redis 의 ZRangeByScore cmd 를 이용해서, score 이상의 결과를 전부 return
func (c *Client) ZRangeByScoreGreaterThan(ctx context.Context, key string, gte float64) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min: fmt.Sprintf("%v", gte),
		Max: "+inf",
	}

	res, err := c.cl.ZRangeByScore(ctx, key, opt).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) ZAddNX(ctx context.Context, key string, members ...Member) error {
	var z []*redis.Z
	for _, m := range members {
		z = append(z, &redis.Z{
			Score:  m.Score,
			Member: m.Value,
		})
	}

	if err := c.cl.ZAddNX(ctx, key, z...).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Client) ZAdd(ctx context.Context, key string, members ...Member) error {
	var z []*redis.Z
	for _, m := range members {
		z = append(z, &redis.Z{
			Score:  m.Score,
			Member: m.Value,
		})
	}

	if err := c.cl.ZAdd(ctx, key, z...).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Client) ZScore(ctx context.Context, key, member string) (float64, error) {
	res, err := c.cl.ZScore(ctx, key, member).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

var (
	TTLNotSet       = errors.New("ttl is not yet set")
	KeyDoesNotExist = errors.New("key does not exist")
)

func (c *Client) TTL(ctx context.Context, key string) (int, error) {
	ttl := c.cl.TTL(ctx, key).Val()
	switch ttl {
	case time.Duration(-2): // -2 if the key does not exist
		return 0, KeyDoesNotExist
	case time.Duration(-1): // -1 if the key exists but has no associated expire
		return 0, TTLNotSet
	default:
		return int(ttl.Seconds()), nil
	}
}
