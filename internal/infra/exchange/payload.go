package exchange

// Helper Types

type ErrorResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors"`
}

type ErrorResponseData struct {
	Error ErrorResponse `json:"error"`
}

type ExchangeErrorResponse struct {
	Errors []string `json:"errors"`
}

// DTOs

type LinkAccountPayload struct {
	Partner       string `json:"partner"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	PartnerUserId string `json:"partner_user_id"`
}

type LinkAccountResponse struct {
	Status      bool   `json:"status"`
	Method      string `json:"method"`
	Partner     string `json:"partner"`
	ChallengeId string `json:"challenge_id"`
}

type SummaryData struct {
	Balance        float64 `json:"balance"`
	DepositAddress string  `json:"deposit_address"`
	IsActive       bool    `json:"is_active"`
}

type ConfirmOtpMetadata struct {
	UserId        string `json:"user_id"`
	WalletAddress string `json:"wallet_address"`
}
type ConfirmOtpPayload struct {
	Partner          string             `json:"partner"`
	PhoneNumber      string             `json:"phone_number"`
	ChallengeId      string             `json:"challenge_id"`
	VerificationCode string             `json:"verification_code"`
	PartnerUserId    string             `json:"partner_user_id"`
	Metadata         ConfirmOtpMetadata `json:"metadata"`
}

type ConfirmOtpResponse struct {
	CustomerId int64  `json:"customer_id"`
	Status     string `json:"status"`
	Linked     bool   `json:"linked"`
	Partner    string `json:"partner"`
	Group      string `json:"group"`
}

type UnlinkAccountPayload struct {
	Partner       string `json:"partner"`
	PartnerUserId string `json:"partner_user_id"`
	CustomerId    string `json:"customer_id"`
}

type WithdrawPayload struct {
	Partner       string  `json:"partner"`
	Amount        float64 `json:"amount"`
	PartnerUserId string  `json:"partner_user_id"`
}

type WithdrawLimit struct {
	Last24Hours      string `json:"last_24_hours"`
	Last24HoursLimit string `json:"last_24_hours_limit"`
}

// *string-field return null, so not sure amount the type
type WithdrawResponse struct {
	Id       int64  `json:"id"`
	MemberId int64  `json:"member_id"`
	TId      string `json:"tid"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Amount   string `json:"amount"`
	Fee      string `json:"fee"`
	// BlockchainTxId *string `json:"blockchain_txid"`
	RId   string `json:"rid"`
	State string `json:"state"`
	// Note        *string `json:"note"`
	// Error       *string `json:"error"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	// DoneAt    *string `json:"done_at"`
}

type RequestWithdrawOtpPayload struct {
	Partner       string `json:"partner"`
	PartnerUserId string `json:"partner_user_id"`
	Channel       string `json:"channel"`
}

type RequestWithdrawOtpResponse struct {
	Sent bool `json:"sent"`
}

type SubmitWithdrawOtpPayload struct {
	Partner       string `json:"partner"`
	PartnerUserId string `json:"partner_user_id"`
	PhoneOtp      string `json:"phone_otp"`
}

type SubmitWithdrawOtpResponse struct {
	Status string `json:"status"`
}
