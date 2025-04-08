package jwt

import "testing"

func TestDecode(t *testing.T) {
	data := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiZ3VhY2FkX2FkZHJcIjpcIlwiLFwiYXNzZXRfcHJvdG9jb2xcIjpcIlwiLFwiYXNzZXRfaG9zdFwiOlwiXCIsXCJhc3NldF9wb3J0XCI6XCJcIixcImFzc2V0X3VzZXJcIjpcIlwiLFwiYXNzZXRfcGFzc3dvcmRcIjpcIlwiLFwiYXNzZXRfcmVtb3RlX2FwcFwiOlwiXCIsXCJhc3NldF9yZW1vdGVfYXBwX2FyZ3NcIjpcIlwiLFwiYXNzZXRfcmVtb3RlX2FwcF9kaXJcIjpcIlwiLFwic2NyZWVuX3dpZHRoXCI6MCxcInNjcmVlbl9oZWlnaHRcIjowLFwic3RvcmFnZV9pZFwiOlwiXCJ9In0.9O3yI3xckwSAqGV4ObcUtVbOxCHQ_ferfL67cU2zSbk"

	s, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(s)

}
