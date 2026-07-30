package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type AdminRepository interface {
	GetAdminByEmail(email string) (*model.Admin, error)
}

type adminRepository struct{}

func NewAdminRepository() AdminRepository {
	return &adminRepository{}
}

func (r *adminRepository) GetAdminByEmail(email string) (*model.Admin, error) {
	admin := &model.Admin{}
	err := intailizer.DB.QueryRow(`
		SELECT admin_id, username, email, password, is_active
		FROM admin_auth
		WHERE LOWER(email) = LOWER($1)
	`, email).Scan(
		&admin.AdminID,
		&admin.Username,
		&admin.Email,
		&admin.Password,
		&admin.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
