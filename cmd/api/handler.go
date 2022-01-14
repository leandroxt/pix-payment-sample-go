package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const ApiMercadoPago = "https://api.mercadopago.com/v1/payments"

func (app *application) processPayment(w http.ResponseWriter, r *http.Request) {
	var input PaymentPixRequest

	if err := readJSON(w, r, &input); err != nil {
		log.Println("Error parse request body", err)
		writeJSON(w, 500, envelope{"error": "cant parse request body"}, nil)
		return
	}

	pix, err := app.createMercadoPagoPayment(input)
	if err != nil {
		log.Println("Error request mercado pago", err)
		writeJSON(w, 500, envelope{"error": "cant create payment"}, nil)
		return
	}

	if err = writeJSON(w, http.StatusCreated, envelope{"pix": pix}, nil); err != nil {
		writeJSON(w, 500, envelope{"error": "internal server error"}, nil)
	}
}

func (app *application) createMercadoPagoPayment(payment PaymentPixRequest) (*PaymentPixResponse, error) {
	mpBodyRequest := &MercadoPagoBodyRequest{
		TransactionAmount: payment.TransactionAmount,
		Description:       payment.Description,
		PaymentMethodID:   "pix",
		Payer: MercadoPagoRequestPayer{
			FirstName: payment.Payer.FirstName,
			LastName:  payment.Payer.LastName,
			Email:     payment.Payer.Email,
			Identification: PayerIdentification{
				Type:   payment.Payer.Identification.Type,
				Number: payment.Payer.Identification.Number,
			},
		},
	}

	body, err := json.Marshal(mpBodyRequest)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, ApiMercadoPago, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.config.mpAccessToken))

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)
	decoder := json.NewDecoder(r)

	pix := &MercadoPagoPostResponse{}
	if err = decoder.Decode(pix); err != nil {
		return nil, err
	}

	return &PaymentPixResponse{
		ID:           pix.ID,
		Status:       pix.Status,
		Detail:       pix.StatusDetail,
		QrCode:       pix.PointOfInteraction.TransactionData.QrCode,
		QrCodeBase64: pix.PointOfInteraction.TransactionData.QrCodeBase64,
	}, nil
}
