package user

import (
	"context"
	"fmt"
	"log"

	"marketgrid/user/internal/domain/port"
)

type EmailSyncService struct {
	userRepo             port.UserRepository
	emailExistenceService port.EmailExistenceService
}

func NewEmailSyncService(userRepo port.UserRepository, emailExistenceService port.EmailExistenceService) *EmailSyncService {
	return &EmailSyncService{
		userRepo:             userRepo,
		emailExistenceService: emailExistenceService,
	}
}

func (s *EmailSyncService) InitializeBloomFilter(ctx context.Context) error {
	log.Println("Initializing Bloom Filter with existing emails...")

	emails, err := s.userRepo.GetAllEmails(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing emails: %w", err)
	}

	if len(emails) == 0 {
		log.Println("No existing emails found, Bloom Filter initialized empty")
		return nil
	}

	err = s.emailExistenceService.InitializeFromEmails(ctx, emails)
	if err != nil {
		return fmt.Errorf("failed to initialize Bloom Filter: %w", err)
	}

	log.Printf("Bloom Filter initialized with %d emails", len(emails))
	return nil
}