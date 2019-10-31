package moonpay

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/goware/urlx"
	"github.com/imroc/req"
)

const apiAddr = "https://api.moonpay.io/v3"

type Moonpay struct {
	u      url.URL
	pubkey string
}

func New(pubkey string) *Moonpay {
	u, _ := urlx.ParseWithDefaultScheme(apiAddr, "https")

	return &Moonpay{
		pubkey: pubkey,
		u:      *u,
	}
}

func (m *Moonpay) url(p string, a ...interface{}) string {
	u := m.u
	u.Path = path.Join("/v3/", fmt.Sprintf(p, a...))
	return u.String()
}

type RespError struct {
	Status string

	Errors []struct {
		Target      map[string]interface{}
		Value       string
		Propery     string
		Constraints map[string]interface{}
	}
	Message string
	Name    string
}

func (e *RespError) Error() string {
	if len(e.Errors) > 0 {
		var ss []string
		for _, er := range e.Errors {
			for _, v := range er.Constraints {
				ss = append(ss, fmt.Sprint(v))
			}
		}

		return strings.Join(ss, "; ")
	}
	if e.Message != "" {
		return e.Message
	}

	return e.Status
}

func (m *Moonpay) handleError(resp *req.Resp, err error) error {
	if err != nil {
		return err
	}
	if r := resp.Response(); r.StatusCode >= 400 {
		err := &RespError{Status: r.Status}
		resp.ToJSON(err)

		return err
	}

	return nil
}

//
//
//

// Currencies returns a list of all currencies supported by MoonPay.
// https://www.moonpay.io/api_reference/v3#list_currencies
func (m *Moonpay) Currencies() (list []Currency, err error) {
	resp, err := req.Get(m.url("/currencies"))
	if err := m.handleError(resp, err); err != nil {
		return nil, err
	}

	err = resp.ToJSON(&list)
	return
}

// CurrencyPrice get the current exchange rates of a currency. Supply the
// currency code, and MoonPay will return the corresponding exchange rates.
// https://www.moonpay.io/api_reference/v3#get_currency_exchange_rate
func (m *Moonpay) CurrencyPrice(crypto string) (prices map[string]float64, err error) {
	resp, err := req.Get(
		m.url("/currencies/%s/price", strings.ToLower(crypto)),
		req.QueryParam{"apiKey": m.pubkey},
	)
	if err := m.handleError(resp, err); err != nil {
		return nil, err
	}

	err = resp.ToJSON(&prices)
	return
}

// CurrenciesPrice get the current exchange rates of multiple cryptocurrencies
// against fiat currencies. Supply the codes of the cryptocurrencies and fiat
// currencies you are interested in, and MoonPay will return the relevant
// exchange rates.
// https://www.moonpay.io/api_reference/v3#get_multiple_exchange_rates
func (m *Moonpay) CurrenciesPrice(crypto, fiat []string) (prices map[string]map[string]float64, err error) {
	resp, err := req.Get(
		m.url("/currencies/price"),
		req.QueryParam{
			"apiKey":           m.pubkey,
			"cryptoCurrencies": strings.Join(crypto, ","),
			"fiatCurrencies":   strings.Join(fiat, ","),
		},
	)
	if err := m.handleError(resp, err); err != nil {
		return nil, err
	}

	err = resp.ToJSON(&prices)
	return
}

// https://www.moonpay.io/api_reference/v3#get_currency_quote
// func (m *Moonpay) CurrencyQuote(crypto, fiat string, fiatAmount, fee float64)

//
//
//

// Countries returnes a list of all countries supported by MoonPay
// https://www.moonpay.io/api_reference/v3#list_countries
func (m *Moonpay) Countries() (countries []Country, err error) {
	resp, err := req.Get(m.url("/countries"))
	if err := m.handleError(resp, err); err != nil {
		return nil, err
	}

	err = resp.ToJSON(&countries)
	return
}

// IPaddress returns information about an IP address
// https://www.moonpay.io/api_reference/v3#check_ip_address
func (m *Moonpay) IPaddress() (ip IPaddress, err error) {
	resp, err := req.Get(m.url("/ip_address"), req.QueryParam{"apiKey": m.pubkey})
	if err := m.handleError(resp, err); err != nil {
		return ip, err
	}

	err = resp.ToJSON(&ip)
	return
}

//
//
//

type email_login_data struct {
	Email              string `json:"email"`
	SecurityCode       string `json:"securityCode,omitempty"`
	ExternalCustomerId string `json:"externalCustomerId,omitempty"`
}

// SecurityCode sends a one-time security code to that email
// is first step for authentication process
// https://www.moonpay.io/api_reference/v3#authenticate_customer_email
func (m *Moonpay) SecurityCode(email string) (bool, error) {
	resp, err := req.Post(
		m.url("/customers/email_login"),
		req.QueryParam{"apiKey": m.pubkey},
		req.BodyJSON(email_login_data{Email: email}),
	)
	if err := m.handleError(resp, err); err != nil {
		return false, err
	}

	var respdata struct{ PreAuthenticated bool }
	err = resp.ToJSON(&respdata)
	return respdata.PreAuthenticated, err
}

// ConfirmRegistration validates the email and authenticates the customer
// is second step for authentication process
// https://www.moonpay.io/api_reference/v3#authenticate_customer_email
func (m *Moonpay) ConfirmRegistration(email, code, extid string) (c CustomerAuth, err error) {
	data := email_login_data{Email: email, SecurityCode: code, ExternalCustomerId: extid}

	resp, err := req.Post(
		m.url("/customers/email_login"),
		req.QueryParam{"apiKey": m.pubkey},
		req.BodyJSON(data),
	)
	if err := m.handleError(resp, err); err != nil {
		return c, err
	}

	err = resp.ToJSON(&c)
	return
}

//
//
//

// MoonpayCustomer is wrapper on the Moonpay instance for join requests for
// specified customer
type MoonpayCustomer struct {
	*Moonpay
	token string
}

func (m *Moonpay) Customer(token string) *MoonpayCustomer {
	return &MoonpayCustomer{m, token}
}

func (m *MoonpayCustomer) authHeader() req.Header {
	return req.Header{"Authorization": "Bearer " + m.token}
}

// RefreshToken refresh the logged-in customer's JWT
// https://www.moonpay.io/api_reference/v3#refresh_token
func (m *MoonpayCustomer) RefreshToken() (c CustomerAuth, err error) {
	resp, err := req.Get(
		m.url("/customers/refresh_token"),
		req.QueryParam{"apiKey": m.pubkey},
		m.authHeader(),
	)
	if err := m.handleError(resp, err); err != nil {
		return c, err
	}

	err = resp.ToJSON(&c)
	if err != nil {
		m.token = c.Token
	}

	return
}

// CustomerInfo retrieves the details of the logged-in customer.
// https://www.moonpay.io/api_reference/v3#retrieve_customer
func (m *MoonpayCustomer) Info() (c Customer, err error) {
	resp, err := req.Get(m.url("/customers/me"), m.authHeader())
	if err := m.handleError(resp, err); err != nil {
		return c, err
	}

	err = resp.ToJSON(&c)
	return
}

// CustomerLimits retrieve the logged-in customer's limits
// https://www.moonpay.io/api_reference/v3#retrieve_customer_limits
func (m *MoonpayCustomer) Limits() (l Limits, err error) {
	resp, err := req.Get(m.url("/customers/me/limits"), m.authHeader())
	if err := m.handleError(resp, err); err != nil {
		return l, err
	}

	err = resp.ToJSON(&l)
	return
}

// UpdateCustomer by setting the values of the parameters passed.
// https://www.moonpay.io/api_reference/v3#update_customer
func (m *MoonpayCustomer) Update(u CustomerFields) (c Customer, err error) {
	resp, err := req.Patch(
		m.url("/customers/me"),
		req.QueryParam{"apiKey": m.pubkey},
		m.authHeader(),
		req.BodyJSON(u),
	)
	log.Println(resp.Dump())
	if err := m.handleError(resp, err); err != nil {
		return c, err
	}

	err = resp.ToJSON(&c)
	return c, err
}

//
// Files
//

//
// Tokens
//

// CreateToken creates a single-use token that represents a credit card’s details.
// https://www.moonpay.io/api_reference/v3#create_token
func (m *Moonpay) CreateToken(data TokenRequest) (t Token, err error) {
	resp, err := req.Post(m.url("/tokens"), req.QueryParam{"apiKey": m.pubkey}, req.BodyJSON(data))
	if err := m.handleError(resp, err); err != nil {
		return t, err
	}

	err = resp.ToJSON(&t)
	return
}

//
// Cards
//

// CreateCard creates a new card object using a token.
// Note that you must provide the user's personal information before being able
// to create a card.
// https://www.moonpay.io/api_reference/v3#create_card
func (m *MoonpayCustomer) CreateCard(tokenid uuid.UUID) (card Card, err error) {
	type data struct {
		TokenID uuid.UUID `json:"tokenId"`
	}
	resp, err := req.Post(m.url("/cards"), m.authHeader(), req.BodyJSON(data{tokenid}))
	if err := m.handleError(resp, err); err != nil {
		return card, err
	}

	err = resp.ToJSON(&card)
	return
}

// Cards returns a list of the cards that you have stored for the logged-in user
// https://www.moonpay.io/api_reference/v3#list_cards
func (m *MoonpayCustomer) Cards() (cards []Card, err error) {
	resp, err := req.Get(m.url("/cards"), m.authHeader())
	if err := m.handleError(resp, err); err != nil {
		return nil, err
	}

	err = resp.ToJSON(&cards)
	return
}

// DeleteCard permanently deletes a card. It cannot be undone
// https://www.moonpay.io/api_reference/v3#delete_card
func (m *MoonpayCustomer) DeleteCard(id uuid.UUID) (card Card, err error) {
	resp, err := req.Delete(m.url("/cards/%s", id), m.authHeader())
	if err := m.handleError(resp, err); err != nil {
		return card, err
	}

	err = resp.ToJSON(&card)
	return
}

//
// Transactions
//

// CreateTransaction creates a new transaction object
// https://www.moonpay.io/api_reference/v3#create_transaction
func (m *MoonpayCustomer) CreateTransaction(data TransactionRequest) (tx Transaction, err error) {
	resp, err := req.Get(m.url("/transactions"), m.authHeader(), req.BodyJSON(data))
	if err := m.handleError(resp, err); err != nil {
		return tx, err
	}

	err = resp.ToJSON(&tx)
	return
}

// Transaction retrieve the details of an existing transaction. Supply the unique
// transaction identifier that was returned upon transaction creation.
// https://www.moonpay.io/api_reference/v3#retrieve_transaction
func (m *MoonpayCustomer) Transaction(id uuid.UUID) (tx Transaction, err error) {
	resp, err := req.Get(m.url("/transactions/%s", id), m.authHeader())
	if err := m.handleError(resp, err); err != nil {
		return tx, err
	}

	err = resp.ToJSON(&tx)
	return
}

// Transactions returns a list of the logged-in customer's transactions
// https://www.moonpay.io/api_reference/v3#list_transactions
func (m *MoonpayCustomer) Transactions() (txs []Transaction, err error) {
	resp, err := req.Get(m.url("/transactions"), m.authHeader())
	if err := m.handleError(resp, err); err != nil {
		return nil, err
	}

	err = resp.ToJSON(&txs)
	return
}
