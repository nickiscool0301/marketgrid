package redis

import (
	"context"
	"fmt"
	"strings"

	"marketgrid/user/domain/port"

	"github.com/redis/go-redis/v9"
)

const (
	bloomFilterKey = "emails:bloom"
	cacheKeyPrefix = "email:exists:"
)

type RedisEmailExistenceService struct {
	client *redis.Client
}

var _ port.EmailExistenceService = (*RedisEmailExistenceService)(nil)

func NewRedisEmailExistenceService(redisClient *RedisClient) *RedisEmailExistenceService {
	return &RedisEmailExistenceService{
		client: redisClient.GetClient(),
	}
}

func (s *RedisEmailExistenceService) EmailExists(ctx context.Context, email string) (bool, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	exists, err := s.client.Do(ctx, "BF.EXISTS", bloomFilterKey, normalizedEmail).Bool()
	if err != nil {
		return s.checkCache(ctx, normalizedEmail)
	}

	if !exists {
		return false, nil
	}

	return s.checkCache(ctx, normalizedEmail)
}

func (s *RedisEmailExistenceService) AddEmail(ctx context.Context, email string) error {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	_, err := s.client.Do(ctx, "BF.ADD", bloomFilterKey, normalizedEmail).Result()
	if err != nil {
		return fmt.Errorf("failed to add email to bloom filter: %w", err)
	}

	cacheKey := cacheKeyPrefix + normalizedEmail
	err = s.client.Set(ctx, cacheKey, "1", 0).Err()
	if err != nil {
		return fmt.Errorf("failed to add email to cache: %w", err)
	}

	return nil
}

func (s *RedisEmailExistenceService) InitializeFromEmails(ctx context.Context, emails []string) error {
	if len(emails) == 0 {
		return nil
	}

	// Create Bloom Filter with estimated capacity
	capacity := len(emails) * 2 // Allow for growth
	errorRate := 0.01           // 1% false positive rate

	_, err := s.client.Do(ctx, "BF.RESERVE", bloomFilterKey, errorRate, capacity).Result()
	if err != nil {
		if !strings.Contains(err.Error(), "item exists") {
			return fmt.Errorf("failed to create bloom filter: %w", err)
		}
	}

	pipe := s.client.Pipeline()

	for _, email := range emails {
		normalizedEmail := strings.ToLower(strings.TrimSpace(email))

		pipe.Do(ctx, "BF.ADD", bloomFilterKey, normalizedEmail)

		cacheKey := cacheKeyPrefix + normalizedEmail
		pipe.Set(ctx, cacheKey, "1", 0)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize emails: %w", err)
	}

	return nil
}

func (s *RedisEmailExistenceService) checkCache(ctx context.Context, normalizedEmail string) (bool, error) {
	cacheKey := cacheKeyPrefix + normalizedEmail
	result := s.client.Get(ctx, cacheKey)

	if result.Err() == redis.Nil {
		return false, nil
	} else if result.Err() != nil {
		return false, fmt.Errorf("failed to check email cache: %w", result.Err())
	}

	return true, nil
}
