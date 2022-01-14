package main

type PayerIdentification struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Payer struct {
	FirstName      string              `json:"firstName"`
	LastName       string              `json:"lastName"`
	Email          string              `json:"email"`
	Identification PayerIdentification `json:"identification"`
}

type PaymentPixRequest struct {
	Description       string `json:"description"`
	TransactionAmount int    `json:"transactionAmount"`
	Payer             Payer  `json:"payer"`
}

type PaymentPixResponse struct {
	ID           int    `json:"id"`
	Status       string `json:"status"`
	Detail       string `json:"detail"`
	QrCode       string `json:"qrCode"`
	QrCodeBase64 string `json:"qrCodeBase64"`
}

// ******************
// ** MERCADO PAGO **
// ******************

type MercadoPagoRequestPayer struct {
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	Email          string              `json:"email"`
	Identification PayerIdentification `json:"identification"`
}

type MercadoPagoBodyRequest struct {
	TransactionAmount int                     `json:"transaction_amount"`
	Description       string                  `json:"description"`
	PaymentMethodID   string                  `json:"payment_method_id"`
	Payer             MercadoPagoRequestPayer `json:"payer"`
}

type MercadoPagoPostResponse struct {
	ID                 int    `json:"id"`
	Status             string `json:"status"`
	StatusDetail       string `json:"status_detail"`
	PointOfInteraction struct {
		TransactionData struct {
			QrCode       string `json:"qr_code"`
			QrCodeBase64 string `json:"qr_code_base64"`
		} `json:"transaction_data"`
	} `json:"point_of_interaction"`
}
