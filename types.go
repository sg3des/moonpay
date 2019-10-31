package moonpay

import (
	"time"

	"github.com/google/uuid"
)

// Currency objects represent the cryptocurrencies supported by MoonPay.
// https://www.moonpay.io/api_reference/v3#currency_object
type Currency struct {
	ID                  uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Type                string
	Name                string
	Code                string
	Precision           int
	AddressRegex        string
	TestnetAddressRegex string
	SupportsAddressTag  bool
	AddressTagRegex     string
	SupportsTestMode    bool
	IsSuspended         bool
	IsSupportedInUS     bool
}

// Country objects represent the countries supported by MoonPay. If the isAllowed
// flag is set to false, it means that MoonPay accepts citizens of this country
// but not residents.
// https://www.moonpay.io/api_reference/v3#country_object
type Country struct {
	Alpha2             string
	Alpha3             string
	IsAllowed          bool
	Name               string
	SupportedDocuments []string
}

const (
	DocumentPassport       = "passport"
	DocumentIDcard         = "national_identity_card"
	DocumentDrivingLicence = "driving_licence"
	DocumentSelfie         = "selfie"

	SideFront = "front"
	SideBack
)

// IP address objects represent the end user's IP address. If the isAllowed flag
// is set to false, it means that MoonPay accepts citizens of this country but
// not residents.
// https://www.moonpay.io/api_reference/v3#ip_address_object
type IPaddress struct {
	Alpha2    string
	Alpha3    string
	State     string
	IPaddress string
	IsAllowed bool
}

// Customer objects represent your end users.
// https://www.moonpay.io/api_reference/v3#customer_object
type Customer struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	FirstName             string    `json:"firstName,omitempty"`
	LastName              string    `json:"lastName,omitempty"`
	Email                 string    `json:"email,omitempty"`
	Phone                 string    `json:"phoneNumber,omitempty"`
	IsPhoneNumberVerified bool      `json:"isPhoneNumberVerified,omitempty"`
	DateOfBirth           time.Time `json:"dateOfBirth,omitempty"`
	SocialSecurityNumber  string    `json:"socialSecurityNumber,omitempty"`

	// LiveMode has the value true if the object exists in live mode or the value
	// false if the object exists in test mode.
	LiveMode bool `json:"liveMode,omitempty"`

	DefaultCurrencyId string `json:"defaultCurrencyId,omitempty"`

	Address Address `json:"address,omitempty"`

	ExternalCustomerID string
}

type CustomerAuth struct {
	CSRFtoken string
	Token     string
	Customer  Customer
}

// CustomerFields is fields allowed to update
type CustomerFields struct {
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	Email                string `json:"email,omitempty"`
	Phone                string `json:"phoneNumber,omitempty"`
	DateOfBirth          string `json:"dateOfBirth,omitempty"`
	SocialSecurityNumber string `json:"socialSecurityNumber,omitempty"`

	DefaultCurrencyId string `json:"defaultCurrencyId,omitempty"`

	Address *Address `json:"address,omitempty"`
}

type Address struct {
	Street    string `json:"street,omitempty"`
	SubStreet string `json:"subStreet,omitempty"`
	Town      string `json:"town,omitempty"`
	PostCode  string `json:"postCode,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
}

// Limits describing the verification levels and limits of the logged-in customer.
type Limits struct {
	Limits []struct {
		Type                  string
		DailyLimit            int
		DailyLimitRemaining   int
		MonthlyLimit          int
		MonthlyLimitRemaining int
	}

	VerificationLevels []struct {
		Name         string
		Requirements []struct {
			Completed  bool
			Identifier string
		}
	}

	LimitIncreaseEligible bool
}

// Card objects represent your end user's credit or debit cards.
// You can save multiple cards on a customer and use them to create transactions.
type Card struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	ExpiryMonth int
	ExpiryYear  int

	Brand          string
	Bin            string
	LastDigits     string
	BillingAddress Address

	CustomerId uuid.UUID
}

// Token objects represent tokenized credit or debit cards.
// Tokenization is the process MoonPay uses to collect sensitive card details
// directly from your customers in a secure manner.
// Tokens cannot be stored or used more than once and they expire after one hour.
type Token struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`

	ExpiryMonth int
	ExpiryYear  int

	Brand          string
	Bin            string
	LastDigits     string
	BillingAddress Address
}

type TokenRequest struct {
	Number     string  `json:"number"`
	ExpiryDate string  `json:"expiryDate"`
	CVC        string  `json:"cvc"`
	Address    Address `json:"address"`
}

// Transaction objects represent cryptocurrency purchases by your end users.
// Cryptocurrency purchases and withdrawals are performed asynchronously.
// You must set up a webhook to be notified of a status change.
// https://www.moonpay.io/api_reference/v3#transaction_object
type Transaction struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	BaseCurrencyAmount  float64 `json:"baseCurrencyAmount"`
	QuoteCurrencyAmount float64 `json:"quoteCurrencyAmount"`
	FeeAmount           float64 `json:"feeAmount"`
	ExtraFeeAmount      float64 `json:"extraFeePercentage"`
	AreFeesIncluded     bool    `json:"areFeesIncluded"`

	Status        string
	FailureReason string

	WalletAddress       string `json:"walletAddress"`
	WalletAddressTag    string `json:"walletAddressTag,omitempty"`
	CryptoTransactionId string

	ReturnURL   string `json:"returnUrl,omitempty"`
	RedirectURL string `json:"redirectUrl,omitempty"`

	BaseCurrencyID uuid.UUID `json:"baseCurrencyId"`
	CurrencyID     uuid.UUID `json:"currencyId"`
	CustomerID     uuid.UUID `json:"customerId"`
	CardID         uuid.UUID `json:"cardId"`

	EURrate float64 `json:"eurRate"`
	USDrate float64 `json:"usdRate"`
	GBPrate float64 `json:"gbpRate"`
}

type TransactionRequest struct {
	BaseCurrencyAmount float64 `json:"baseCurrencyAmount"`
	ExtraFeeAmount     float64 `json:"extraFeePercentage"`
	AreFeesIncluded    bool    `json:"areFeesIncluded"`

	WalletAddress    string `json:"walletAddress"`
	WalletAddressTag string `json:"walletAddressTag,omitempty"`

	BaseCurrencyCode string `json:"baseCurrencyCode"`
	CurrencyCode     string `json:"currencyCode"`

	ReturnURL string `json:"returnUrl"`

	TokenID string `json:"tokenId"`
	CardID  string `json:"cardId"`
}
