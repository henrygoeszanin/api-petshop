package errors

import "errors"

// Erros comuns do domínio
var (
	ErrNotFound      = errors.New("registro não encontrado")
	ErrAlreadyExists = errors.New("registro já existe")
	ErrInvalidData   = errors.New("dados inválidos")
	ErrUnauthorized  = errors.New("não autorizado")
	ErrForbidden     = errors.New("acesso proibido")
)
