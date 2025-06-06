package api

type RequestBody struct {
	Query string `json:"query"`
}

type ResponseBody struct {
	Answer string `json:"answer"`
}

type MultipleCityRequestBody struct {
	Query string `json:"query"`
}

type MultipleCityResponseBody struct {
	Reports map[string]string `json:"reports"`
}

type MultipleAsyncRequestBody struct {
	Cities []string `json:"cities"`
}

type MultipleAsyncResponseBody struct {
	Reports map[string]string `json:"reports"`
}
