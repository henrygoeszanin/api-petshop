package services

import (
	"fmt"

	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_petshop/domain/entities"
	"github.com/henrygoeszanin/api_petshop/domain/errors"
)

// AuthService fornece métodos para autenticação de donos e petshops
type AuthService struct {
	donoRepository    repositories.DonoRepository
	petshopRepository repositories.PetshopRepository
}

// NewAuthService cria uma nova instância de AuthService
func NewAuthService(donoRepo repositories.DonoRepository, petshopRepo repositories.PetshopRepository) *AuthService {
	return &AuthService{
		donoRepository:    donoRepo,
		petshopRepository: petshopRepo,
	}
}

// AuthenticateDono autentica um dono com base no email e senha
func (s *AuthService) AuthenticateDono(email, password string) (*dtos.DonoResponseDTO, error) {
	dono, err := s.donoRepository.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar dono: %w", err)
	}

	if dono == nil {
		return nil, errors.ErrUnauthorized
	}

	if !dono.CheckPassword(password) {
		return nil, errors.ErrUnauthorized
	}

	// Criar DTO de resposta
	response := &dtos.DonoResponseDTO{
		AuthResponseDTO: dtos.AuthResponseDTO{
			ID:    dono.ID,
			Email: dono.Email,
			Nome:  dono.Nome,
			Tipo:  "dono",
		},
		Telefone: dono.Telefone,
	}

	return response, nil
}

// AuthenticatePetshop autentica um petshop com base no email e senha
func (s *AuthService) AuthenticatePetshop(email, password string) (*dtos.PetshopResponseDTO, error) {
	petshop, err := s.petshopRepository.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar petshop: %w", err)
	}

	if petshop == nil {
		return nil, errors.ErrUnauthorized
	}

	if !petshop.CheckPassword(password) {
		return nil, errors.ErrUnauthorized
	}

	// Criar DTO de resposta
	response := &dtos.PetshopResponseDTO{
		AuthResponseDTO: dtos.AuthResponseDTO{
			ID:    petshop.ID,
			Email: petshop.Email,
			Nome:  petshop.Nome,
			Tipo:  "petshop",
		},
		Telefone:  petshop.Telefone,
		Descricao: petshop.Descricao,
		Nota:      petshop.Nota,
	}

	return response, nil
}

// RegisterDono cria uma nova conta de Dono
func (s *AuthService) RegisterDono(dto *dtos.DonoRegisterDTO) (*dtos.DonoResponseDTO, error) {
	// Verificar duplicidade por email
	existing, err := s.donoRepository.GetByEmail(dto.Email)
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar dono existente: %w", err)
	}
	if existing != nil {
		return nil, errors.ErrAlreadyExists
	}
	// Criar entidade Dono
	dono := &entities.Dono{
		Nome:        dto.Nome,
		Email:       dto.Email,
		Telefone:    dto.Telefone,
		CEP:         dto.CEP,
		Rua:         dto.Rua,
		Bairro:      dto.Bairro,
		Cidade:      dto.Cidade,
		Estado:      dto.Estado,
		Numero:      dto.Numero,
		Complemento: dto.Complemento,
	}
	// Gerar hash da senha
	if err := dono.SetPassword(dto.Password); err != nil {
		return nil, fmt.Errorf("falha ao definir senha: %w", err)
	}
	// Salvar no repositório
	if err := s.donoRepository.Create(dono); err != nil {
		return nil, fmt.Errorf("falha ao criar dono: %w", err)
	}
	// Preparar DTO de resposta
	resp := &dtos.DonoResponseDTO{
		AuthResponseDTO: dtos.AuthResponseDTO{
			ID:    dono.ID,
			Email: dono.Email,
			Nome:  dono.Nome,
			Tipo:  "dono",
		},
		Telefone: dono.Telefone,
	}
	return resp, nil
}

// RegisterPetshop cria uma nova conta de Petshop
func (s *AuthService) RegisterPetshop(dto *dtos.PetshopRegisterDTO) (*dtos.PetshopResponseDTO, error) {
	// Verificar duplicidade por email
	existing, err := s.petshopRepository.GetByEmail(dto.Email)
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar petshop existente: %w", err)
	}
	if existing != nil {
		return nil, errors.ErrAlreadyExists
	}
	// Criar entidade Petshop
	petshop := &entities.Petshop{
		Nome:        dto.Nome,
		Email:       dto.Email,
		Telefone:    dto.Telefone,
		CEP:         dto.CEP,
		Rua:         dto.Rua,
		Bairro:      dto.Bairro,
		Cidade:      dto.Cidade,
		Estado:      dto.Estado,
		Numero:      dto.Numero,
		Complemento: dto.Complemento,
		Descricao:   dto.Descricao,
	}
	// Gerar hash da senha
	if err := petshop.SetPassword(dto.Password); err != nil {
		return nil, fmt.Errorf("falha ao definir senha: %w", err)
	}
	// Salvar no repositório
	if err := s.petshopRepository.Create(petshop); err != nil {
		return nil, fmt.Errorf("falha ao criar petshop: %w", err)
	}
	// Preparar DTO de resposta
	resp := &dtos.PetshopResponseDTO{
		AuthResponseDTO: dtos.AuthResponseDTO{
			ID:    petshop.ID,
			Email: petshop.Email,
			Nome:  petshop.Nome,
			Tipo:  "petshop",
		},
		Telefone:  petshop.Telefone,
		Descricao: petshop.Descricao,
		Nota:      petshop.Nota,
	}
	return resp, nil
}
