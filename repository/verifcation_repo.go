package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"database/sql"
)

type VerificationRepository interface {
	UploadVerification(authID int, front string, back string) error

	UploadPlayerDocument(doc model.PlayerDocument) error

	GetVerification(authID int) (*model.Verification, error)

	GetPlayerDocuments(authID int) ([]model.PlayerDocument, error)

	GetPlayerDocumentByID(docID int) (*model.PlayerDocument, error)

	GetPendingVerifications() ([]model.Verification, error)

	ApproveVerification(authID int) error

	RejectVerification(authID int, reason string) error

	ApprovePlayerDocument(docID int) error

	RejectPlayerDocument(docID int, reason string) error
	GetPendingPlayerDocuments() ([]model.PlayerDocument, error)
}

type verificationRepository struct{}

func NewVerificationRepository() VerificationRepository {
	return &verificationRepository{}
}

// upload citizenship
func (r *verificationRepository) UploadVerification(authID int, front string, back string) error {

	query := `
INSERT INTO user_verification
(
auth_id,
citizenship_front,
citizenship_back,
verification_status
)
VALUES($1,$2,$3,'pending')
ON CONFLICT(auth_id)
DO UPDATE SET

citizenship_front=EXCLUDED.citizenship_front,
citizenship_back=EXCLUDED.citizenship_back,
verification_status='pending',
updated_at=NOW(),
rejection_reason=NULL,
verified_at=NULL;
`

	_, err := intailizer.DB.Exec(query,
		authID,
		front,
		back,
	)

	return err
}

// upload player document
func (r *verificationRepository) UploadPlayerDocument(doc model.PlayerDocument) error {

	query := `
INSERT INTO player_documents
(
auth_id,
document_name,
document_type,
document_url,
issued_by,
issue_date,
verification_status
)

VALUES($1,$2,$3,$4,$5,$6,'pending')
`

	_, err := intailizer.DB.Exec(
		query,
		doc.AuthID,
		doc.DocumentName,
		doc.DocumentType,
		doc.DocumentURL,
		doc.IssuedBy,
		doc.IssueDate,
	)

	return err
}

// get verification
func (r *verificationRepository) GetVerification(authID int) (*model.Verification, error) {

	query := `
SELECT

verification_id,
auth_id,
citizenship_front,
citizenship_back,
verification_status,
rejection_reason,
verified_at,
created_at,
updated_at

FROM user_verification
WHERE auth_id=$1
`

	var v model.Verification

	err := intailizer.DB.QueryRow(query, authID).Scan(

		&v.VerificationID,
		&v.AuthID,
		&v.CitizenshipFront,
		&v.CitizenshipBack,
		&v.Status,
		&v.RejectionReason,
		&v.VerifiedAt,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &v, nil

}

// get player certificates
func (r *verificationRepository) GetPlayerDocuments(authID int) ([]model.PlayerDocument, error) {

	query := `
SELECT

document_id,
auth_id,
document_name,
document_type,
document_url,
issued_by,
issue_date,
verification_status,
rejection_reason,
verified_at,
created_at

FROM player_documents
WHERE auth_id=$1
ORDER BY created_at DESC
`

	rows, err := intailizer.DB.Query(query, authID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var docs []model.PlayerDocument

	for rows.Next() {

		var d model.PlayerDocument

		err := rows.Scan(

			&d.DocumentID,
			&d.AuthID,
			&d.DocumentName,
			&d.DocumentType,
			&d.DocumentURL,
			&d.IssuedBy,
			&d.IssueDate,
			&d.Status,
			&d.RejectionReason,
			&d.VerifiedAt,
			&d.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		docs = append(docs, d)

	}

	return docs, nil

}

// get a single player document by id
func (r *verificationRepository) GetPlayerDocumentByID(docID int) (*model.PlayerDocument, error) {

	query := `
SELECT

document_id,
auth_id,
document_name,
document_type,
document_url,
issued_by,
issue_date,
verification_status,
rejection_reason,
verified_at,
created_at

FROM player_documents
WHERE document_id=$1
`

	var d model.PlayerDocument

	err := intailizer.DB.QueryRow(query, docID).Scan(

		&d.DocumentID,
		&d.AuthID,
		&d.DocumentName,
		&d.DocumentType,
		&d.DocumentURL,
		&d.IssuedBy,
		&d.IssueDate,
		&d.Status,
		&d.RejectionReason,
		&d.VerifiedAt,
		&d.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &d, nil
}

// get pending verifications(admin)
func (r *verificationRepository) GetPendingVerifications() ([]model.Verification, error) {

	query := `
SELECT

verification_id,
auth_id,
citizenship_front,
citizenship_back,
verification_status,
rejection_reason,
verified_at,
created_at,
updated_at

FROM user_verification
WHERE verification_status='pending'
ORDER BY created_at ASC
`

	rows, err := intailizer.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []model.Verification

	for rows.Next() {

		var v model.Verification

		err := rows.Scan(

			&v.VerificationID,
			&v.AuthID,
			&v.CitizenshipFront,
			&v.CitizenshipBack,
			&v.Status,
			&v.RejectionReason,
			&v.VerifiedAt,
			&v.CreatedAt,
			&v.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, v)

	}

	return list, nil

}

// approve verification(admin)
func (r *verificationRepository) ApproveVerification(authID int) error {

	tx, err := intailizer.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
UPDATE user_verification

SET

verification_status='verified',
verified_at=NOW(),
updated_at=NOW()

WHERE auth_id=$1
`, authID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
UPDATE cometoseeauth

SET is_verified=true

WHERE auth_id=$1
`, authID)

	if err != nil {
		return err
	}

	return tx.Commit()

}

// reject verification(admin)
func (r *verificationRepository) RejectVerification(authID int, reason string) error {

	tx, err := intailizer.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()
	_, err = tx.Exec(`
UPDATE user_verification

SET

verification_status='rejected',
rejection_reason=$1,
updated_at=NOW()

WHERE auth_id=$2
`, reason, authID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
UPDATE cometoseeauth

SET is_verified=false

WHERE auth_id=$1
`, authID)

	if err != nil {
		return err
	}

	return tx.Commit()

}

// approve player document(admin)
func (r *verificationRepository) ApprovePlayerDocument(docID int) error {

	_, err := intailizer.DB.Exec(`
UPDATE player_documents

SET

verification_status='verified',
verified_at=NOW(),
rejection_reason=NULL

WHERE document_id=$1
`, docID)

	return err
}

// reject player document(admin)
func (r *verificationRepository) RejectPlayerDocument(docID int, reason string) error {

	_, err := intailizer.DB.Exec(`
UPDATE player_documents

SET

verification_status='rejected',
rejection_reason=$1

WHERE document_id=$2
`, reason, docID)

	return err
}

// get pending player documents(admin)
func (r *verificationRepository) GetPendingPlayerDocuments() ([]model.PlayerDocument, error) {

	query := `
SELECT
	document_id,
	auth_id,
	document_name,
	document_type,
	document_url,
	issued_by,
	issue_date,
	verification_status,
	rejection_reason,
	verified_at,
	created_at

FROM player_documents
WHERE verification_status='pending'
ORDER BY created_at ASC
`

	rows, err := intailizer.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []model.PlayerDocument

	for rows.Next() {

		var d model.PlayerDocument

		err := rows.Scan(
			&d.DocumentID,
			&d.AuthID,
			&d.DocumentName,
			&d.DocumentType,
			&d.DocumentURL,
			&d.IssuedBy,
			&d.IssueDate,
			&d.Status,
			&d.RejectionReason,
			&d.VerifiedAt,
			&d.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		docs = append(docs, d)
	}

	return docs, rows.Err()
}
