package stripe

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestPaymentIntentLastPaymentError_UnmarshalJSON(t *testing.T) {
	errorData := map[string]interface{}{
		"charge":       "ch_123",
		"code":         "card_declined",
		"decline_code": "generic_decline",
		"doc_url":      "https://stripe.com/docs/error-codes/card-declined",
		"message":      "Your card was declined.",
		"source": map[string]interface{}{
			"id":          "card_123",
			"object":      "card",
			"brand":       "Visa",
			"country":     "US",
			"customer":    "cus_123",
			"exp_month":   9,
			"exp_year":    2019,
			"fingerprint": "fingerprint",
			"last4":       "0341",
		},
		"type": "card_error",
	}
	bytes, err := json.Marshal(&errorData)
	assert.NoError(t, err)

	var lastPaymentError PaymentIntentLastPaymentError
	err = json.Unmarshal(bytes, &lastPaymentError)
	assert.NoError(t, err)

	assert.Equal(t, ErrorTypeCard, lastPaymentError.Type)
	assert.Equal(t, "ch_123", lastPaymentError.Charge)
	assert.Equal(t, "https://stripe.com/docs/error-codes/card-declined", lastPaymentError.DocURL)
	assert.Equal(t, PaymentSourceTypeCard, lastPaymentError.Source.Type)
	assert.Equal(t, "card_123", lastPaymentError.Source.Card.ID)
}

func TestPaymentIntentNextAction_UnmarshalJSON(t *testing.T) {
	actionData := map[string]interface{}{
		"redirect_to_url": map[string]interface{}{
			"return_url": "https://stripe.com/return",
			"url":        "https://stripe.com",
		},
		"type": "redirect_to_url",
	}

	bytes, err := json.Marshal(&actionData)
	assert.NoError(t, err)

	var action PaymentIntentNextAction
	err = json.Unmarshal(bytes, &action)
	assert.NoError(t, err)

	assert.Equal(t, PaymentIntentNextActionTypeRedirectToURL, action.Type)
	assert.Equal(t, "https://stripe.com", action.RedirectToURL.URL)
	assert.Equal(t, "https://stripe.com/return", action.RedirectToURL.ReturnURL)
}

func TestPaymentIntent_UnmarshalJSON(t *testing.T) {
	intentData := map[string]interface{}{
		"id":     "pi_123",
		"object": "payment_intent",
		"charges": map[string]interface{}{
			"object":   "list",
			"has_more": true,
			"data": []map[string]interface{}{
				{
					"id":     "ch_123",
					"object": "charge",
				},
				{
					"id":     "ch_234",
					"object": "charge",
				},
			},
		},
		"payment_method_types": []interface{}{
			"card",
		},
	}

	bytes, err := json.Marshal(&intentData)
	assert.NoError(t, err)

	var intent PaymentIntent
	err = json.Unmarshal(bytes, &intent)
	assert.NoError(t, err)

	assert.Equal(t, "pi_123", intent.ID)

	assert.Equal(t, 2, len(intent.Charges.Data))
	assert.Equal(t, 1, len(intent.PaymentMethodTypes))
}
