package keycloak

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	url2 "net/url"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"
)

// Service define las operaciones del servicio de Keycloak
type Service interface {
	ValidateToken(tokenString string) (*KeycloakClaims, error)
	GetUserInfo(tokenString string) (*UserInfo, error)
	GetUserByID(userID string) (*UserInfo, error)
	Login(username, password string) (string, error)
	CreateUser(user *CreateUserRequest) error
	UpdateUser(userID string, user *UpdateUserRequest) error
	DeleteUser(userID string) error
	GetUsers() ([]*UserInfo, error)
	GetUserGroups(userID string) ([]*Group, error)
	AddUserToGroup(userID, groupID string) error
	RemoveUserFromGroup(userID, groupID string) error
}

// service implementa el servicio de Keycloak
type service struct {
	baseURL      string
	realm        string
	clientID     string
	clientSecret string
	httpClient   *http.Client

	publicKey       interface{}
	publicKeyExpiry time.Time
}

// NewService crea una nueva instancia del servicio de Keycloak
func NewService(baseURL, realm, clientID, clientSecret string) Service {
	return &service{
		baseURL:      strings.TrimSuffix(baseURL, "/"),
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

const publicKeyTTL = time.Hour

// KeycloakClaims representa los claims del token JWT de Keycloak
type KeycloakClaims struct {
	Sub               string                 `json:"sub"`
	Email             string                 `json:"email"`
	EmailVerified     bool                   `json:"email_verified"`
	Name              string                 `json:"name"`
	PreferredUsername string                 `json:"preferred_username"`
	GivenName         string                 `json:"given_name"`
	FamilyName        string                 `json:"family_name"`
	RealmAccess       *RealmAccess           `json:"realm_access"`
	ResourceAccess    map[string]*ClientRole `json:"resource_access"`
	ClientID          string                 `json:"client_id"`
	Username          string                 `json:"username"`
	Active            bool                   `json:"active"`
	Exp               int64                  `json:"exp"`
	Iat               int64                  `json:"iat"`
	Iss               string                 `json:"iss"`
	Aud               interface{}            `json:"aud"`
	Typ               string                 `json:"typ"`
	AuthTime          int64                  `json:"auth_time"`
	SessionState      string                 `json:"session_state"`
	Acr               string                 `json:"acr"`
	AllowedOrigins    []string               `json:"allowed-origins"`
	Realm             string                 `json:"realm"`
	TokenType         string                 `json:"token_type"`
	Nonce             string                 `json:"nonce"`
	Jti               string                 `json:"jti"`
	Azp               string                 `json:"azp"`
	Scope             string                 `json:"scope"`
	ClientHost        string                 `json:"client_host"`
	ClientAddress     string                 `json:"client_address"`
	CustomClaims      map[string]interface{} `json:"-"`
}

// RealmAccess representa los roles del realm
type RealmAccess struct {
	Roles []string `json:"roles"`
}

// ClientRole representa los roles de un cliente específico
type ClientRole struct {
	Roles []string `json:"roles"`
}

// UserInfo representa la información de un usuario de Keycloak
type UserInfo struct {
	ID            string                 `json:"id"`
	Username      string                 `json:"username"`
	Email         string                 `json:"email"`
	EmailVerified bool                   `json:"emailVerified"`
	FirstName     string                 `json:"firstName"`
	LastName      string                 `json:"lastName"`
	Enabled       bool                   `json:"enabled"`
	Created       int64                  `json:"createdTimestamp"`
	Attributes    map[string][]string    `json:"attributes,omitempty"`
	Groups        []string               `json:"groups,omitempty"`
	CustomClaims  map[string]interface{} `json:"-"`
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Username      string              `json:"username"`
	Email         string              `json:"email"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	Enabled       bool                `json:"enabled"`
	EmailVerified bool                `json:"emailVerified"`
	Credentials   []*Credential       `json:"credentials,omitempty"`
	Groups        []string            `json:"groups,omitempty"`
	Attributes    map[string][]string `json:"attributes,omitempty"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Username      string              `json:"username,omitempty"`
	Email         string              `json:"email,omitempty"`
	FirstName     string              `json:"firstName,omitempty"`
	LastName      string              `json:"lastName,omitempty"`
	Enabled       *bool               `json:"enabled,omitempty"`
	EmailVerified *bool               `json:"emailVerified,omitempty"`
	Groups        []string            `json:"groups,omitempty"`
	Attributes    map[string][]string `json:"attributes,omitempty"`
}

// Credential representa las credenciales de un usuario
type Credential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

// Group representa un grupo de Keycloak
type Group struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Path       string              `json:"path"`
	SubGroups  []*Group            `json:"subGroups,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
}

// ValidateToken valida un token JWT de Keycloak
func (s *service) ValidateToken(tokenString string) (*KeycloakClaims, error) {
	// Obtener la clave pública de Keycloak
	publicKey, err := s.getPublicKey()
	if err != nil {
		return nil, fmt.Errorf("error getting public key: %w", err)
	}

	// Parsear y validar el token
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	var claims KeycloakClaims
	if err := token.Claims(publicKey, &claims); err != nil {
		// reintentar refrescando la clave
		if publicKey, err = s.refreshPublicKey(); err != nil {
			return nil, fmt.Errorf("error verifying token: %w", err)
		}
		if err := token.Claims(publicKey, &claims); err != nil {
			return nil, fmt.Errorf("error verifying token: %w", err)
		}
	}

	// Validar tiempo de expiración
	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	// Validar issuer
	expectedIssuer := fmt.Sprintf("%s/realms/%s", s.baseURL, s.realm)
	if claims.Iss != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	return &claims, nil
}

// GetUserInfo obtiene información del usuario desde el token
func (s *service) GetUserInfo(tokenString string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", s.baseURL, s.realm)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting user info: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *service) GetUserByID(userID string) (*UserInfo, error) {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", s.baseURL, s.realm, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting user: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// CreateUser crea un nuevo usuario en Keycloak
func (s *service) CreateUser(user *CreateUserRequest) error {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users", s.baseURL, s.realm)

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(userData)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating user: %d", resp.StatusCode)
	}

	return nil
}

// UpdateUser actualiza un usuario existente
func (s *service) UpdateUser(userID string, user *UpdateUserRequest) error {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", s.baseURL, s.realm, userID)

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(userData)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error updating user: %d", resp.StatusCode)
	}

	return nil
}

// DeleteUser elimina un usuario
func (s *service) DeleteUser(userID string) error {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", s.baseURL, s.realm, userID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting user: %d", resp.StatusCode)
	}

	return nil
}

// GetUsers obtiene todos los usuarios
func (s *service) GetUsers() ([]*UserInfo, error) {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users", s.baseURL, s.realm)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting users: %d", resp.StatusCode)
	}

	var users []*UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserGroups obtiene los grupos de un usuario
func (s *service) GetUserGroups(userID string) ([]*Group, error) {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s/groups", s.baseURL, s.realm, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting user groups: %d", resp.StatusCode)
	}

	var groups []*Group
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// AddUserToGroup agrega un usuario a un grupo
func (s *service) AddUserToGroup(userID, groupID string) error {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s/groups/%s", s.baseURL, s.realm, userID, groupID)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error adding user to group: %d", resp.StatusCode)
	}

	return nil
}

// RemoveUserFromGroup remueve un usuario de un grupo
func (s *service) RemoveUserFromGroup(userID, groupID string) error {
	accessToken, err := s.getAdminToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s/groups/%s", s.baseURL, s.realm, userID, groupID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error removing user from group: %d", resp.StatusCode)
	}

	return nil
}

// getPublicKey obtiene la clave pública de Keycloak
func (s *service) getPublicKey() (interface{}, error) {
	if s.publicKey != nil && time.Now().Before(s.publicKeyExpiry) {
		return s.publicKey, nil
	}
	url := fmt.Sprintf("%s/realms/%s", s.baseURL, s.realm)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("requesting realm info: %w", err)
	}
	defer resp.Body.Close()

	var realmInfo struct {
		PublicKey string `json:"public_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&realmInfo); err != nil {
		return nil, fmt.Errorf("decoding realm info: %w", err)
	}

	pemData := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", realmInfo.PublicKey)
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("failed to decode public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %w", err)
	}

	s.publicKey = key
	s.publicKeyExpiry = time.Now().Add(publicKeyTTL)

	return key, nil
}

func (s *service) refreshPublicKey() (interface{}, error) {
	s.publicKey = nil
	s.publicKeyExpiry = time.Time{}
	return s.getPublicKey()
}

// getAdminToken obtiene un token de administrador
func (s *service) getAdminToken() (string, error) {
	url := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", s.baseURL)

	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", s.clientID, s.clientSecret)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting admin token: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// Login authenticates a user using username and password and returns an access token
func (s *service) Login(username, password string) (string, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.baseURL, s.realm)

	data := url2.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting user token: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
