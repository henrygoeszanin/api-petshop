package errors

import "errors"

// Erros genéricos
var (
	ErrNotFound      = errors.New("registro não encontrado")
	ErrAlreadyExists = errors.New("registro já existe")
	ErrInvalidData   = errors.New("dados inválidos")
	ErrUnauthorized  = errors.New("não autorizado")
	ErrForbidden     = errors.New("acesso proibido")
)

// Erros de autenticação e credenciais
var (
	ErrInvalidCredentials  = errors.New("credenciais inválidas")
	ErrInvalidID           = errors.New("ID inválido")
	ErrSetPassword         = errors.New("falha ao definir senha")
	ErrFailedToSetPassword = errors.New("falha ao definir senha")
)

// Erros relacionados a Dono
var (
	ErrFailedToFetchDono   = errors.New("falha ao buscar dono")
	ErrFailedToCheckDono   = errors.New("falha ao verificar dono")
	ErrFailedToCreateDono  = errors.New("falha ao criar dono")
	ErrCheckExistingOwner  = errors.New("falha ao verificar dono existente")
	ErrCreateOwner         = errors.New("falha ao criar dono")
	ErrUpdateOwner         = errors.New("falha ao atualizar dono")
	ErrUpdateOwnerLocation = errors.New("falha ao atualizar localização do dono")
	ErrDonoNotFound        = errors.New("dono não encontrado")
)

// Erros relacionados a Petshop
var (
	ErrFailedToFetchPetshop  = errors.New("falha ao buscar petshop")
	ErrFailedToCheckPetshop  = errors.New("falha ao verificar petshop")
	ErrFailedToCreatePetshop = errors.New("falha ao criar petshop")
	ErrCheckExistingPetshop  = errors.New("falha ao verificar petshop existente")
	ErrUpdatePetshop         = errors.New("falha ao atualizar petshop")
	ErrUpdatePetshopLocation = errors.New("falha ao atualizar localização do petshop")
	ErrPetshopNotFound       = errors.New("petshop não encontrado")
)

// Erros relacionados a Pet
var (
	ErrFailedToCheckPet  = errors.New("falha ao verificar pet")
	ErrFailedToCreatePet = errors.New("falha ao criar pet")
	ErrUpdatePet         = errors.New("falha ao atualizar pet")
	ErrPetNotFound       = errors.New("pet não encontrado")
	ErrPetNotOwnedByDono = errors.New("o pet não pertence ao dono informado")
)

// Erros relacionados a Serviço
var (
	ErrFailedToCheckService  = errors.New("falha ao verificar serviço existente")
	ErrFailedToCreateService = errors.New("falha ao criar serviço")
	ErrFailedToUpdateService = errors.New("falha ao atualizar serviço")
	ErrFailedToFetchServices = errors.New("falha ao buscar serviços")
	ErrServiceNotFound       = errors.New("serviço não encontrado")
	ErrServiceNotFromPetshop = errors.New("o serviço não pertence ao petshop informado")
)

// Erros relacionados a Procedimento
var (
	ErrInvalidDate            = errors.New("formato de data inválido, use ISO 8601")
	ErrFutureDate             = errors.New("a data de realização não pode ser futura")
	ErrTotalMismatch          = errors.New("o total informado não corresponde à soma dos preços finais")
	ErrFailedToCheckProcedure = errors.New("falha ao verificar procedimento")
)

// Erros relacionados a Agendamento
var (
	ErrInvalidAgendamentoStatus   = errors.New("status de agendamento inválido")
	ErrPastDate                   = errors.New("a data de agendamento não pode ser no passado")
	ErrFailedToCheckAgendamento   = errors.New("falha ao verificar agendamento")
	ErrServiceInactive            = errors.New("o serviço não está ativo")
	ErrTotalPrevistoMismatch      = errors.New("o total previsto não corresponde à soma dos preços dos itens")
	ErrFailedToCreateAgendamento  = errors.New("falha ao criar agendamento")
	ErrFailedToUpdateAgendamento  = errors.New("falha ao atualizar agendamento")
	ErrFailedToUpdateStatus       = errors.New("falha ao atualizar status do agendamento")
	ErrFailedToFetchAgendamentos  = errors.New("falha ao buscar agendamentos")
	ErrFailedToFetchPetInfo       = errors.New("falha ao buscar informações do pet")
	ErrFailedToFetchDonoInfo      = errors.New("falha ao buscar informações do dono")
	ErrFailedToFetchPetshopInfo   = errors.New("falha ao buscar informações do petshop")
	ErrAgendamentoUpdateForbidden = errors.New("não é possível atualizar um agendamento cancelado ou concluído")
	ErrUpdateCanceledAgendamento  = errors.New("não é possível alterar o status de um agendamento cancelado")
	ErrUpdateCompletedAgendamento = errors.New("não é possível alterar o status de um agendamento concluído")
)
