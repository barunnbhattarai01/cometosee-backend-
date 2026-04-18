package controller

import (
	"cometosee/model"
	"cometosee/repository"
	"cometosee/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mukezhz/pay-np/errorz"
	"github.com/mukezhz/pay-np/esewa"
	"github.com/mukezhz/pay-np/utils"
)

type EsewaClient struct {
	Payload        *model.EsewaPayload
	Signature      string
	Secret         string
	signatureMap   map[string]string
	ReponsePayload *model.EsewaVerifyPayload
}

func New(secret string, payload *model.EsewaPayload) (*EsewaClient, error) {
	p := EsewaClient{
		Payload: payload,
		Secret:  secret,
	}
	sm, err := setupSignatureMap[model.EsewaPayload](*p.Payload)
	if err != nil {
		return nil, err
	}
	p.signatureMap = *sm
	return &p, nil
}

func setupSignatureMap[T model.EsewaPayload | model.EsewaVerifyPayload](p T) (*map[string]string, error) {
	j, _ := json.Marshal(p)

	var data map[string]any
	err := json.Unmarshal(j, &data)
	if err != nil {
		return nil, err
	}

	stringMap := make(map[string]string)
	for key, value := range data {
		stringMap[key] = fmt.Sprintf("%v", value)
	}
	return &stringMap, nil
}

func (e *EsewaClient) validate() error {
	// total_amount,transaction_uuid,product_code mandatory fields
	if e.Payload.TotalAmount == "" {
		return errorz.ErrEsewaTotalAmount
	}
	if e.Payload.TransactionUUID == "" {
		return errorz.ErrEsewaTransactionUUID
	}
	if e.Payload.ProductCode == "" {
		return errorz.ErrEsewaProductCode
	}
	return nil
}

func (e *EsewaClient) getInputForSignature(signedFieldNames string) (string, error) {
	splittedSignedFieldNames := strings.Split(signedFieldNames, ",")
	if len(splittedSignedFieldNames) < 3 {
		return "", errorz.ErrEsewaInvalidDataForSignature
	}

	var signatureDate []string
	for _, signedFieldName := range splittedSignedFieldNames {
		if e.signatureMap[signedFieldName] == "" {
			return "", errorz.ErrEsewaInvalidDataForSignature
		}
		signatureDate = append(signatureDate, fmt.Sprintf("%s=%s", signedFieldName, e.signatureMap[signedFieldName]))
	}
	return strings.Join(signatureDate, ","), nil
}

func (e *EsewaClient) GenerateSignature() (string, error) {
	err := e.validate()
	if err != nil {
		return "", err
	}
	data, err := e.getInputForSignature(e.Payload.SignedFieldNames)
	if err != nil {
		return "", err
	}
	return utils.HmacSHA256(e.Secret, data), nil
}

func (e *EsewaClient) VerifySignature(data string) error {
	d, err := utils.Base64Decode(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, &e.ReponsePayload)
	if err != nil {
		return err
	}
	sm, err := setupSignatureMap[model.EsewaVerifyPayload](*e.ReponsePayload)
	if err != nil {
		return err
	}
	e.signatureMap = *sm
	i, err := e.getInputForSignature(e.ReponsePayload.SignedFieldNames)
	if err != nil {
		return err
	}
	signature := utils.HmacSHA256(e.Secret, i)
	if e.ReponsePayload.Signature != signature {
		return errorz.ErrEsewaInvalidSignature
	}

	return nil
}

// helper fucntion
func InitiateHandler(w http.ResponseWriter, r *http.Request) {
	payload := &esewa.EsewaPayload{
		Amount:                "100",
		TaxAmount:             "0",
		ProductServiceCharge:  "0",
		ProductDeliveryCharge: "0",
		ProductCode:           "EPAYTEST", // eSewa sandbox merchant code
		TotalAmount:           "100",
		TransactionUUID:       uuid.New().String(),
		SuccessURL:            "http://localhost:8080/esewa/verify",
		FailureURL:            "http://localhost:8080/esewa/failure",
		SignedFieldNames:      "total_amount,transaction_uuid,product_code",
	}

	client, err := esewa.New("8gBm/:&EnhH.1/q", payload) // sandbox secret
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sig, err := client.GenerateSignature()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	payload.Signature = sig
	w.Header().Set("Content-Type", "application/json")

	subrepo := repository.NewSubscriptionRepository()
	subservice := service.NewSubscriptionService(subrepo)
	c := NewSubscriptionController(subservice)
	if c.GetSubscriptionStatus("barunnbhattarai@gmail.com") == false {
		c.SubscribeUser("barunnbhattarai@gmail.com")
	} else {
		c.UpdateSubscriptionEndDate("barunnbhattarai@gmail.com")
	}

	json.NewEncoder(w).Encode(payload)

}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "missing data param", 400)
		return
	}

	client, _ := esewa.New("8gBm/:&EnhH.1/q", &esewa.EsewaPayload{})
	err := client.VerifySignature(data)
	if err != nil {
		http.Error(w, "invalid signature: "+err.Error(), 400)
		return
	}

	//suncrintion logic
	subrepo := repository.NewSubscriptionRepository()
	subservice := service.NewSubscriptionService(subrepo)
	c := NewSubscriptionController(subservice)

	if c.GetSubscriptionStatus("barunnbhattarai@gmail.com") == false {
		c.SubscribeUser("barunnbhattarai@gmail.com")
		return
	}
	w.Write([]byte("payment verified ok"))

}
