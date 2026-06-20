package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	esewaProductCode = "EPAYTEST"
)

type PaymentController struct {
	paymentService service.PaymentService
	subService     service.SubscriptionService
	esewaSecret    string
	successURL     string
	failureURL     string
}

func NewPaymentController(
	paymentService service.PaymentService,
	subService service.SubscriptionService,
	esewaSecret, successURL, failureURL string,
) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
		subService:     subService,
		esewaSecret:    esewaSecret,
		successURL:     successURL,
		failureURL:     failureURL,
	}
}

var plans = map[string]float64{
	"monthly": 100.00,
	"yearly":  1000.00,
}

func (pc *PaymentController) InitiateHandler(w http.ResponseWriter, r *http.Request) {
	authID := common.GetAuthid(r.Context())

	var req struct {
		Plan string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Plan == "" {
		common.WriteJSONError(w, "plan is required", http.StatusBadRequest)
		return
	}

	amount, ok := plans[req.Plan]
	if !ok {
		common.WriteJSONError(w, "invalid plan", http.StatusBadRequest)
		return
	}
	amountStr := strconv.FormatFloat(amount, 'f', 0, 64)

	txnUUID := uuid.New().String()
	payload := &model.EsewaPayload{
		Amount:                amountStr,
		TaxAmount:             "0",
		ProductServiceCharge:  "0",
		ProductDeliveryCharge: "0",
		ProductCode:           esewaProductCode,
		TotalAmount:           amountStr,
		TransactionUUID:       txnUUID,
		SuccessURL:            pc.successURL,
		FailureURL:            pc.failureURL,
		SignedFieldNames:      "total_amount,transaction_uuid,product_code",
	}

	client, err := New(pc.esewaSecret, payload)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sig, err := client.GenerateSignature()
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	payload.Signature = sig

	if err := pc.paymentService.InitiatePayment(int(authID), txnUUID, amount, req.Plan); err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func (pc *PaymentController) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if data == "" {
		common.WriteJSONError(w, "missing data param", http.StatusBadRequest)
		return
	}

	client, _ := New(pc.esewaSecret, &model.EsewaPayload{})
	if err := client.VerifySignature(data); err != nil {
		common.WriteJSONError(w, "invalid signature: "+err.Error(), http.StatusBadRequest)
		return
	}
	resp := client.ReponsePayload

	payment, err := pc.paymentService.GetPayment(resp.TransactionUUID)
	if err != nil {
		common.WriteJSONError(w, "unknown transaction", http.StatusBadRequest)
		return
	}

	if payment.Status != "pending" {
		common.WriteJSONMessage(w, "transaction already processed")
		return
	}

	expectedAmountStr := strconv.FormatFloat(payment.Amount, 'f', 0, 64)
	if resp.TotalAmount != expectedAmountStr {
		pc.paymentService.FailPayment(resp.TransactionUUID)
		common.WriteJSONError(w, "amount mismatch", http.StatusBadRequest)
		return
	}

	if err := pc.paymentService.ConfirmPayment(resp.TransactionUUID); err != nil {
		common.WriteJSONError(w, "failed to record payment", http.StatusInternalServerError)
		return
	}

	if err := pc.subService.SubscribeUser(payment.AuthID, payment.Plan); err != nil {
		common.WriteJSONError(w, "failed to update subscription", http.StatusInternalServerError)
		return
	}

	common.WriteJSONMessage(w, "payment verified, subscription updated")
}

func (pc *PaymentController) FailureHandler(w http.ResponseWriter, r *http.Request) {
	txnUUID := r.URL.Query().Get("transaction_uuid")
	if txnUUID != "" {
		pc.paymentService.FailPayment(txnUUID)
	}
	common.WriteJSONMessage(w, "payment failed")
}
