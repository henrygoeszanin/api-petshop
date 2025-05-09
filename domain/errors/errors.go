package errors

import "errors"

// Erros comuns do domínio
var (
	ErrNotFound              = errors.New("registro não encontrado")
	ErrAlreadyExists         = errors.New("registro já existe")
	ErrInvalidData           = errors.New("dados inválidos")
	ErrUnauthorized          = errors.New("não autorizado")
	ErrForbidden             = errors.New("acesso proibido")
	ErrInvalidCredentials    = errors.New("credenciais inválidas")
	ErrFailedToFetchDono     = errors.New("falha ao buscar dono")
	ErrFailedToFetchPetshop  = errors.New("falha ao buscar petshop")
	ErrFailedToCheckDono     = errors.New("falha ao verificar dono")
	ErrFailedToCheckPetshop  = errors.New("falha ao verificar petshop")
	ErrFailedToSetPassword   = errors.New("falha ao definir senha")
	ErrFailedToCreateDono    = errors.New("falha ao criar dono")
	ErrFailedToCreatePetshop = errors.New("falha ao criar petshop")
	ErrCheckExistingOwner    = errors.New("falha ao verificar dono existente")
	ErrSetPassword           = errors.New("falha ao definir senha")
	ErrCreateOwner           = errors.New("falha ao criar dono")
	ErrUpdateOwner           = errors.New("falha ao atualizar dono")
	ErrUpdateOwnerLocation   = errors.New("falha ao atualizar localização do dono")
	ErrCheckExistingPetshop  = errors.New("falha ao verificar petshop existente")
)
