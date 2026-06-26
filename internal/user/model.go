package user

import "time"

type User struct {
	ID                         string    `json:"id"`
	Username                   string    `json:"username"`
	DisplayName                string    `json:"displayName"`
	AuthKeyHash                string    `json:"-"`
	EncryptedMasterKey         []byte    `json:"encryptedMasterKey"`
	Salt                       []byte    `json:"salt"`
	IV                         []byte    `json:"iv"`
	RecoveryEncryptedMasterKey []byte    `json:"recoveryEncryptedMasterKey,omitempty"`
	RecoverySalt               []byte    `json:"recoverySalt,omitempty"`
	RecoveryIV                 []byte    `json:"recoveryIv,omitempty"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	AuthKeyHash string `json:"auth_key_hash"`
}
