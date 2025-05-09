package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_petshop/application/dtos"
	"github.com/henrygoeszanin/api_petshop/application/services"
	"github.com/henrygoeszanin/api_petshop/config"
	"github.com/segmentio/ksuid"
)

// TokenExtractor é um middleware que extrai o token JWT de diferentes fontes na requisição
// seguindo uma ordem de prioridade: cookies, header de autorização e por fim query parameters.
// Uma vez encontrado, o token é adicionado ao header Authorization para processamento pelos
// middlewares subsequentes.
func TokenExtractor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 1. Tenta obter do cookie - primeira fonte de prioridade
		token, _ = c.Cookie("jwt")
		if token == "" {
			token, _ = c.Cookie("token") // fallback para cookie alternativo
		}

		// 2. Tenta obter do header Authorization - segunda fonte de prioridade
		if token == "" {
			auth := c.GetHeader("Authorization")
			if auth != "" && strings.HasPrefix(auth, "Bearer ") {
				// Remove o prefixo "Bearer " para extrair apenas o token
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		// 3. Tenta obter de query parameter - última fonte de prioridade
		if token == "" {
			token = c.Query("token")
		}

		// Se um token foi encontrado em qualquer fonte, adicionamos ao header de Autorização
		// para que possa ser processado pelos middlewares de autenticação JWT
		if token != "" {
			c.Request.Header.Set("Authorization", "Bearer "+token)
			fmt.Printf("TokenExtractor: Token encontrado\n")
		} else {
			fmt.Printf("TokenExtractor: Nenhum token encontrado\n")
		}

		// Continua a execução da cadeia de middlewares
		c.Next()
	}
}

// PetshopRequired é um middleware que verifica se o usuário autenticado é um petshop
// através da propriedade "tipo" nas claims do JWT. Caso contrário, a requisição é abortada com
// status 403 Forbidden.
func PetshopRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims do token JWT da requisição atual
		claims := jwt.ExtractClaims(c)

		// Verifica se a claim "tipo" existe e se seu valor é "petshop"
		tipo, exists := claims["tipo"]
		if !exists || tipo != "petshop" {
			// Retorna erro 403 para usuários não petshop
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Acesso restrito a petshops",
			})
			c.Abort() // Interrompe a execução de middlewares subsequentes
			return
		}

		// Se for petshop, permite a continuação da requisição
		c.Next()
	}
}

// DonoRequired é um middleware que verifica se o usuário autenticado é um dono
// através da propriedade "tipo" nas claims do JWT. Caso contrário, a requisição é abortada com
// status 403 Forbidden.
func DonoRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai as claims do token JWT da requisição atual
		claims := jwt.ExtractClaims(c)

		// Verifica se a claim "tipo" existe e se seu valor é "dono"
		tipo, exists := claims["tipo"]
		if !exists || tipo != "dono" {
			// Retorna erro 403 para usuários não dono
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Acesso restrito a donos de pets",
			})
			c.Abort() // Interrompe a execução de middlewares subsequentes
			return
		}

		// Se for dono, permite a continuação da requisição
		c.Next()
	}
}

// login define a estrutura esperada para as requisições de autenticação
// contendo email, senha e tipo do usuário. As tags de binding garantem validação
// básica dos campos.
type login struct {
	Email    string `json:"email" binding:"required,email"` // Email validado pelo formato
	Password string `json:"password" binding:"required"`    // Senha obrigatória
	UserType string `json:"user_type" binding:"required"`   // Tipo de usuário: "dono" ou "petshop"
}

// SetupJWTMiddleware configura e retorna uma instância do middleware JWT para autenticação
// recebendo uma instância do serviço de autenticação e configurações do sistema.
func SetupJWTMiddleware(authService *services.AuthService, cfg *config.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "petshop-api",         // Nome do domínio de autenticação
		Key:         []byte(cfg.JWTSecret), // Chave secreta para assinatura dos tokens
		Timeout:     time.Hour * 24,        // Duração de validade do token: 24 horas
		MaxRefresh:  time.Hour * 24 * 7,    // Período máximo em que o token pode ser renovado: 7 dias
		IdentityKey: "id",                  // Chave que identifica o usuário nas claims

		// Configurações de cookies de autenticação
		SendCookie:     true,                     // Envia token como cookie
		CookieName:     "jwt",                    // Nome do cookie
		CookieMaxAge:   24 * time.Hour,           // Tempo de vida do cookie
		CookieDomain:   "",                       // Domínio do cookie (vazio = domínio atual)
		SecureCookie:   false,                    // Cookie não requer HTTPS (alterar para true em produção)
		CookieHTTPOnly: true,                     // Cookie não acessível via JavaScript (proteção XSS)
		CookieSameSite: http.SameSiteDefaultMode, // Política de SameSite do cookie

		// Configuração de lookup do token
		TokenLookup:   "cookie:jwt,header:Authorization", // Ordem de busca do token
		TokenHeadName: "Bearer",                          // Prefixo esperado no header
		TimeFunc:      time.Now,                          // Função para obter a hora atual

		// Authenticator: função responsável por validar as credenciais do usuário
		// e retornar os dados do usuário para geração do token JWT
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			// Lê e valida os dados de login do corpo da requisição
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			fmt.Printf("Login - Tentativa para email: %s (tipo: %s)\n", loginVals.Email, loginVals.UserType)

			// Com base no tipo de usuário, autentica como dono ou petshop
			switch loginVals.UserType {
			case "dono":
				dono, err := authService.AuthenticateDono(loginVals.Email, loginVals.Password)
				if err != nil {
					fmt.Printf("Login Dono - Falha: %v\n", err)
					return nil, jwt.ErrFailedAuthentication
				}
				fmt.Printf("Login Dono - Sucesso para: %s (ID: %s)\n", dono.Email, dono.ID)
				return dono, nil

			case "petshop":
				petshop, err := authService.AuthenticatePetshop(loginVals.Email, loginVals.Password)
				if err != nil {
					fmt.Printf("Login Petshop - Falha: %v\n", err)
					return nil, jwt.ErrFailedAuthentication
				}
				fmt.Printf("Login Petshop - Sucesso para: %s (ID: %s)\n", petshop.Email, petshop.ID)
				return petshop, nil

			default:
				fmt.Printf("Login - Tipo de usuário inválido: %s\n", loginVals.UserType)
				return nil, jwt.ErrFailedAuthentication
			}
		},

		// PayloadFunc: função que define quais dados do usuário serão incluídos no token JWT
		// como claims personalizadas
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Printf("PayloadFunc recebeu dados do tipo: %T\n", data)

			// Verificar o tipo de resposta e extrair os dados apropriados
			if dono, ok := data.(*dtos.DonoResponseDTO); ok {
				fmt.Printf("Convertido com sucesso para DonoResponseDTO: ID=%s, Email=%s\n",
					dono.ID, dono.Email)

				// Define as claims que serão adicionadas ao token JWT
				return jwt.MapClaims{
					"id":       dono.ID.String(),
					"email":    dono.Email,
					"nome":     dono.Nome,
					"tipo":     dono.Tipo,
					"telefone": dono.Telefone,
				}
			} else if petshop, ok := data.(*dtos.PetshopResponseDTO); ok {
				fmt.Printf("Convertido com sucesso para PetshopResponseDTO: ID=%s, Email=%s\n",
					petshop.ID, petshop.Email)

				// Define as claims que serão adicionadas ao token JWT
				return jwt.MapClaims{
					"id":        petshop.ID.String(),
					"email":     petshop.Email,
					"nome":      petshop.Nome,
					"tipo":      petshop.Tipo,
					"telefone":  petshop.Telefone,
					"descricao": petshop.Descricao,
					"nota":      petshop.Nota,
				}
			}

			// Tratamento de erro para tipo inesperado
			fmt.Printf("AVISO: Falha na conversão para DTO conhecido. Tentando alternativas...\n")
			fmt.Printf("Conteúdo do objeto recebido: %+v\n", data)

			// Retorna claims vazias se a conversão falhar
			return jwt.MapClaims{}
		},

		// IdentityHandler: extrai a identidade do usuário das claims do JWT
		// durante a validação do token
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			fmt.Printf("IdentityHandler - Claims recebidas: %+v\n", claims)

			// Verifica se as claims básicas estão presentes
			idStr, idExists := claims["id"].(string)
			emailVal, emailExists := claims["email"]
			nomeVal, nomeExists := claims["nome"]
			tipoVal, tipoExists := claims["tipo"]

			if !idExists || !emailExists || !nomeExists || !tipoExists {
				fmt.Printf("ALERTA: JWT incompleto ou inválido. Claims: %+v\n", claims)
				return nil
			}

			// Conversão segura do ID
			id, err := ksuid.Parse(idStr)
			if err != nil {
				fmt.Printf("ERRO: Campo 'id' não é um KSUID válido: %v\n", idStr)
				return nil
			}

			// Determina o tipo de usuário e constrói o objeto apropriado
			tipo := tipoVal.(string)

			authResponse := dtos.AuthResponseDTO{
				ID:    id,
				Email: emailVal.(string),
				Nome:  nomeVal.(string),
				Tipo:  tipo,
			}

			if tipo == "dono" {
				telefoneVal, ok := claims["telefone"].(string)
				if !ok {
					telefoneVal = ""
				}
				return &dtos.DonoResponseDTO{
					AuthResponseDTO: authResponse,
					Telefone:        telefoneVal,
				}
			} else if tipo == "petshop" {
				telefoneVal, ok := claims["telefone"].(string)
				if !ok {
					telefoneVal = ""
				}

				descricaoVal, ok := claims["descricao"].(string)
				if !ok {
					descricaoVal = ""
				}

				var notaVal float32
				if nota, ok := claims["nota"].(float64); ok {
					notaVal = float32(nota)
				}

				return &dtos.PetshopResponseDTO{
					AuthResponseDTO: authResponse,
					Telefone:        telefoneVal,
					Descricao:       descricaoVal,
					Nota:            notaVal,
				}
			}

			return nil
		},

		// Authorizator: define regras de autorização após a autenticação
		// Aqui estamos permitindo acesso a qualquer usuário autenticado
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// Nesta implementação básica, qualquer usuário autenticado está autorizado
			// Para regras de autorização específicas, este é o lugar para implementá-las
			return data != nil
		},

		// LoginResponse: personaliza a resposta HTTP em caso de login bem-sucedido
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			fmt.Println("==== Login Bem-Sucedido ====")
			fmt.Printf("Token gerado: %s\n", token)
			fmt.Printf("Expira em: %v\n", expire)

			// Retorna o token e sua data de expiração no corpo da resposta
			c.JSON(code, gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},

		// Unauthorized: personaliza a resposta HTTP em caso de falha na autenticação
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Printf("==== Falha na Autenticação ====\n")
			fmt.Printf("Rota: %s %s\n", c.Request.Method, c.Request.URL.Path)
			fmt.Printf("Código: %d, Mensagem: %s\n", code, message)

			// Retorna um erro JSON com o código e mensagem apropriados
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
}
