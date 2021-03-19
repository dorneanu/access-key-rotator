package entity

// AccessKey represents a key/credential/passwort used to authenticate against APIs/services
type AccessKey struct {
	ID     string
	Secret string
}

// EncryptedKey holds an encrypted representation of an AccessKey
type EncryptedKey struct {
	ID     string
	Secret []byte
}
