# MoonPay API Client

API Client for [MoonPay](https://www.moonpay.io/)

### Install

```sh
go get github.com/sg3des/moonpay
```

### Usage


```go
mpay := moonpay.New("....key....")


cardtoken, err := mpay.CreateToken(...TokenRequest...)
if err != nil {...}


// create new `customer`
email := "email@something.com"
mpay.SecurityCode(emai)

// then approve email with recieved code
cauth, err:= mpay.ConfirmRegistration(email, "EMAIL-CODE", "any-uuid")
if err != nil {...}

// cauth.Token - is token of this customer


// create customer card
card, err := mpay.Customer(cauth.Token).CreateCard(cardtoken.ID)
if err != nil {...}



// create transaction
tx, err := mpay.Customer(cauth.Token).CreateTransaction(...TransactionRequest...)
if err != nil {...}
```