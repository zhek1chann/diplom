package payment

import (
	"diploma/modules/cart/client/payment/converter"
	clientModel "diploma/modules/cart/client/payment/model"
	"diploma/modules/cart/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type PaymentClient struct {
	checkoutURL      string
	merchantID       string
	merchantPassword string
	callbackURL      string
}

func NewPaymentClient(checkoutURL, merchantID, merchantPassword, callbackURL string) *PaymentClient {
	return &PaymentClient{
		checkoutURL:      checkoutURL,
		merchantID:       merchantID,
		merchantPassword: merchantPassword,
		callbackURL:      callbackURL,
	}
}

func (p *PaymentClient) PaymentRequest(orderID, amount, currency, orderDesc string) (model.CheckoutResponse, error) {
	checkoutRequest := clientModel.CheckoutRequest{
		OrderID:           orderID,
		MerchantID:        p.merchantID,
		OrderDesc:         orderDesc,
		Amount:            amount,
		Currency:          currency,
		ServerCallbackURL: p.callbackURL,
	}

	checkoutRequest.SetSignature(p.merchantPassword)

	apiRequest := clientModel.APIRequest{
		Request: checkoutRequest,
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		return model.CheckoutResponse{}, fmt.Errorf("error encoding request: %w", err)
	}

	resp, err := http.Post(p.checkoutURL, "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return model.CheckoutResponse{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.CheckoutResponse{}, fmt.Errorf("error reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return model.CheckoutResponse{}, fmt.Errorf("error response from server: %s", string(body))
	}
	var apiResponse clientModel.APIResponse
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return model.CheckoutResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return converter.ToServiceCheckoutFromClient(apiResponse.Response)
}
