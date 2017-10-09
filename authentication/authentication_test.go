package authentication

import (
	"bytes"
	"encoding/json"
	"gopkg.in/square/go-jose.v2/jwt"
	"testing"
)

func TestJWT(t *testing.T) {

	var td = []byte(`{
  "tx_id": "0f534ffc-9442-414c-b39f-a756b4adc6cb",
  "type" : "uk.gov.ons.edc.eq:surveyresponse",
  "version" : "0.0.1",
  "origin" : "uk.gov.ons.edc.eq",
  "survey_id": "021",
  "flushed": false,
  "collection":{
    "exercise_sid": "hfjdskf",
    "instrument_id": "yui789",
    "period": "2016-02-01"
  },
  "submitted_at": "2016-03-07T15:28:05Z",
  "metadata": {
    "user_id": "789473423",
    "ru_ref": "432423423423"
  },
  "data": [{
    "value": "Joe",
    "block_id": "household-composition",
    "answer_id": "household-first-name",
    "group_id": "multiple-questions-group",
    "group_instance": 0,
    "answer_instance": 0
  },
  {
    "value": ["Eggs", "Bacon", "Spam"],
    "block_id": "breakfast-block",
    "answer_id": "favourite-breakfast-food",
    "group_id": "breakfast-group",
    "group_instance": 0,
    "answer_instance": 0
  }]
}`)
	var mappedData map[string]interface{}

	expected := new(bytes.Buffer)
	json.Compact(expected, td)

	err := json.Unmarshal(td, &mappedData)
	if err != nil {
		t.Error("unmarshal error", err)
	}

	jwe, jerr := GetJwe(mappedData, "/Users/andrewtorrance/go/src/sdx-encrypted-submitter/authentication/testPrivateKey.pem",
		"/Users/andrewtorrance/go/src/sdx-encrypted-submitter/authentication/testPublicKey.pem")
	if jerr != nil {
		t.Error("GetJwe returned Error: ", jerr)
	}

	publicKeyResult, _ := loadEncryptionKey("/Users/andrewtorrance/go/src/sdx-encrypted-submitter/authentication/testPublicKey.pem")
	privateKeyResult, _ := loadSigningKey("/Users/andrewtorrance/go/src/sdx-encrypted-submitter/authentication/testPrivateKey.pem")

	parsedData, err := jwt.ParseSignedAndEncrypted(jwe)
	if err != nil {
		t.Error("Could not Parse encrypted Jwe: ", err)
	}

	decryptedData, err := parsedData.Decrypt(privateKeyResult.key)
	if err != nil {
		t.Error("Could not Parse Jwe:", err)
	}

	claimData := make(map[string]interface{})
	err = decryptedData.Claims(publicKeyResult.key, &claimData)
	if err != nil {
		t.Error("Could not extract claims Jwe:", err)
	}

	result, err := json.Marshal(claimData)
	if err != nil {
		t.Error("Could not marshal  Jwe:", err)
	}

	// Note to self:
	// It is not obvious how to compare result with expected , the mapping process does not maintain order (by design)
	// and we do not have a type that we can use therefore this test is far from satisfactory at this point
	if len(string(result)) != len(expected.String()) {
		t.Error("Result and Expected Do not Match ")
	}

}
