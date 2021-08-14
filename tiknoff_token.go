package main

const tokenUrl = "https://www.tinkoff.ru/api/common/v1/session?appName=invest&appVersion=1.153.2&origin=web%2Cib5%2Cplatform"

type TokenResponse struct {
	ResultCode string `json:"resultCode"`
	Payload    string `json:"payload"`
	TrackingId string `json:"trackingId"`
}

func getToken() (response string, ok bool) {
	tokenParcel := TokenResponse{}

	err := getJson(tokenUrl, &tokenParcel)
	if err != nil {
		return "", false
	}

	return tokenParcel.Payload, true
}
