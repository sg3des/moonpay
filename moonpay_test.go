package moonpay

import (
	"os"
	"testing"

	"github.com/Pallinder/go-randomdata"
)

var testMoonpay = New(os.Getenv("MOONPAY_KEY"))

func TestCurrencies(t *testing.T) {
	list, err := testMoonpay.Currencies()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(list) == 0 {
		t.Error("currencies not found")
	}

	for _, c := range list {
		t.Log(c.ID, c.CreatedAt, c.UpdatedAt)
	}
}

func TestCurrencyPrice(t *testing.T) {
	prices, err := testMoonpay.CurrencyPrice("btc")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(prices)
}

func TestCurrenciesPrice(t *testing.T) {
	prices, err := testMoonpay.CurrenciesPrice([]string{"btc", "bch"}, []string{"eur", "gbp"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(prices)
}

//
//
//

func TestCountries(t *testing.T) {
	countries, err := testMoonpay.Countries()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(countries) == 0 {
		t.Error("countries not found")
	}

	t.Log(countries)
}

func TestIPaddress(t *testing.T) {
	ip, err := testMoonpay.IPaddress()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(ip)
}

//
//
//

func TestSecurityCode(t *testing.T) {
	preauth, err := testMoonpay.SecurityCode(os.Getenv("TEST_EMAIL"))
	if err != nil {
		t.Error(err)
	}

	t.Log(preauth)
}

func TestConfirmRegistration(t *testing.T) {
	code := os.Getenv("TEST_CODE")
	if code == "" {
		t.Skip("security code not specified, check email and set code to TEST_CODE environment variable")
	}

	c, err := testMoonpay.ConfirmRegistration(os.Getenv("TEST_EMAIL"), code, "some-test-id")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", c)
}

func TestRefreshToken(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Skip("auth token is not specified, set valid token to TEST_TOKEN environment variable")
	}

	c, err := testMoonpay.Customer(token).RefreshToken()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", c)
}

func TestCustomerInfo(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Skip("auth token is not specified, set valid token to TEST_TOKEN environment variable")
	}

	c, err := testMoonpay.Customer(token).Info()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", c)
}

func TestCustomerLimits(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Skip("auth token is not specified, set valid token to TEST_TOKEN environment variable")
	}

	l, err := testMoonpay.Customer(token).Limits()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", l)
}

func TestUpdateCustomer(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Skip("auth token is not specified, set valid token to TEST_TOKEN environment variable")
	}

	newfirstname := randomdata.FirstName(randomdata.Male)

	c, err := testMoonpay.Customer(token).Update(CustomerFields{FirstName: newfirstname})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if c.FirstName != newfirstname {
		t.Errorf("failed set firstname, exist: %s - expect: %s", c.FirstName, newfirstname)
	}

	t.Logf("%+v", c)
}
