package service

import (
	"cometosee/model"
	"cometosee/repository"
	"errors"
	"strings"
)

type VerificationService interface {
	UploadVerification(authID int, frontURL string, backURL string) error
	UploadPlayerDocument(doc model.PlayerDocument) error
	GetVerification(authID int) (*model.Verification, error)
	GetPlayerDocuments(authID int) ([]model.PlayerDocument, error)
	GetPendingVerifications() ([]model.Verification, error)
	ApproveVerification(authID int) error
	RejectVerification(authID int, reason string) error
}

type verificationService struct {
	repo repository.VerificationRepository
}

func NewVerificationService(
	repo repository.VerificationRepository,
) VerificationService {

	return &verificationService{
		repo: repo,
	}

}

// upload citizenship
func (s *verificationService) UploadVerification(

	authID int,
	frontURL string,
	backURL string,

) error {

	if authID <= 0 {
		return errors.New("invalid user")
	}

	if strings.TrimSpace(frontURL) == "" {
		return errors.New("citizenship front image is required")
	}

	if strings.TrimSpace(backURL) == "" {
		return errors.New("citizenship back image is required")
	}

	return s.repo.UploadVerification(

		authID,
		frontURL,
		backURL,
	)

}

// upload player docs
func (s *verificationService) UploadPlayerDocument(

	doc model.PlayerDocument,

) error {

	if doc.AuthID <= 0 {
		return errors.New("invalid user")
	}

	if strings.TrimSpace(doc.DocumentName) == "" {
		return errors.New("document name is required")
	}

	if strings.TrimSpace(doc.DocumentType) == "" {
		return errors.New("document type is required")
	}

	if strings.TrimSpace(doc.DocumentURL) == "" {
		return errors.New("document url is required")
	}

	return s.repo.UploadPlayerDocument(doc)

}

// get verfication status
func (s *verificationService) GetVerification(

	authID int,

) (*model.Verification, error) {

	if authID <= 0 {
		return nil, errors.New("invalid user")
	}

	return s.repo.GetVerification(authID)

}

// get player docs
func (s *verificationService) GetPlayerDocuments(

	authID int,

) ([]model.PlayerDocument, error) {

	if authID <= 0 {
		return nil, errors.New("invalid user")
	}

	return s.repo.GetPlayerDocuments(authID)

}

// get pending verifications(admin)
func (s *verificationService) GetPendingVerifications() ([]model.Verification, error) {

	return s.repo.GetPendingVerifications()

}

// approve verification(admin)
func (s *verificationService) ApproveVerification(

	authID int,

) error {

	if authID <= 0 {
		return errors.New("invalid auth id")
	}

	user, err := s.repo.GetVerification(authID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("verification not found")
	}

	if user.Status == "verified" {
		return errors.New("user already verified")
	}

	return s.repo.ApproveVerification(authID)

}

// reject verfication (admin)
func (s *verificationService) RejectVerification(

	authID int,
	reason string,

) error {

	if authID <= 0 {
		return errors.New("invalid auth id")
	}

	if strings.TrimSpace(reason) == "" {
		return errors.New("rejection reason is required")
	}

	user, err := s.repo.GetVerification(authID)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("verification not found")
	}

	if user.Status == "rejected" {
		return errors.New("verification already rejected")
	}

	return s.repo.RejectVerification(authID, reason)

}
