package azure

// ResponseHeaders struct
type ResponseHeaders struct {
	ContentLength             []string `json:"Content-Length"`
	OperationLocation         []string `json:"Operation-Location"`
	XEnvoyUpstreamServiceTime []string `json:"X-Envoy-Upstream-Service-Time"`
	APIMRequestID             []string `json:"Apim-Request-Id"`
	StrictTransportSecurity   []string `json:"Strict-Transport-Security"`
	XContentTypeOptions       []string `json:"X-Content-Type-Options"`
	XMsRegion                 []string `json:"X-Ms-Region"`
	Date                      []string `json:"Date"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Error struct {
		Code       string `json:"code"`
		Innererror struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"innererror"`
		Message string `json:"message"`
	} `json:"error"`
}

// ResultAnalyzeResponse struct
type ResultAnalyzeResponse struct {
	AnalyzeResult struct {
		APIVersion string `json:"apiVersion"`
		Content    string `json:"content"`
		Documents  []struct {
			BoundingRegions []struct {
				PageNumber float64   `json:"pageNumber"`
				Polygon    []float64 `json:"polygon"`
			} `json:"boundingRegions"`
			Confidence float64 `json:"confidence"`
			DocType    string  `json:"docType"`
			Fields     struct {
				Items struct {
					Type       string `json:"type"`
					ValueArray []struct {
						BoundingRegions []struct {
							PageNumber float64   `json:"pageNumber"`
							Polygon    []float64 `json:"polygon"`
						} `json:"boundingRegions"`
						Confidence float64 `json:"confidence"`
						Content    string  `json:"content"`
						Spans      []struct {
							Length float64 `json:"length"`
							Offset float64 `json:"offset"`
						} `json:"spans"`
						Type        string                 `json:"type"`
						ValueObject map[string]interface{} `json:"valueObject"`
					} `json:"valueArray"`
				} `json:"Items"`
				MerchantAddress struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type         string `json:"type"`
					ValueAddress struct {
						City          string `json:"city"`
						Road          string `json:"road"`
						StreetAddress string `json:"streetAddress"`
					} `json:"valueAddress"`
				} `json:"MerchantAddress"`
				MerchantName struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type        string `json:"type"`
					ValueString string `json:"valueString"`
				} `json:"MerchantName"`
				MerchantPhoneNumber struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type             string `json:"type"`
					ValuePhoneNumber string `json:"valuePhoneNumber"`
				} `json:"MerchantPhoneNumber"`
				TaxDetails struct {
					Type       string `json:"type"`
					ValueArray []struct {
						BoundingRegions []struct {
							PageNumber float64   `json:"pageNumber"`
							Polygon    []float64 `json:"polygon"`
						} `json:"boundingRegions"`
						Confidence float64 `json:"confidence"`
						Content    string  `json:"content"`
						Spans      []struct {
							Length float64 `json:"length"`
							Offset float64 `json:"offset"`
						} `json:"spans"`
						Type        string `json:"type"`
						ValueObject struct {
							Amount struct {
								BoundingRegions []struct {
									PageNumber float64   `json:"pageNumber"`
									Polygon    []float64 `json:"polygon"`
								} `json:"boundingRegions"`
								Confidence float64 `json:"confidence"`
								Content    string  `json:"content"`
								Spans      []struct {
									Length float64 `json:"length"`
									Offset float64 `json:"offset"`
								} `json:"spans"`
								Type          string `json:"type"`
								ValueCurrency struct {
									Amount         float64 `json:"amount"`
									CurrencyCode   string  `json:"currencyCode"`
									CurrencySymbol string  `json:"currencySymbol"`
								} `json:"valueCurrency"`
							} `json:"Amount"`
						} `json:"valueObject"`
					} `json:"valueArray"`
				} `json:"TaxDetails"`
				Total struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type        string  `json:"type"`
					ValueNumber float64 `json:"valueNumber"`
				} `json:"Total"`
				TotalTax struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type        string  `json:"type"`
					ValueNumber float64 `json:"valueNumber"`
				} `json:"TotalTax"`
				TransactionDate struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type      string `json:"type"`
					ValueDate string `json:"valueDate"`
				} `json:"TransactionDate"`
				TransactionTime struct {
					BoundingRegions []struct {
						PageNumber float64   `json:"pageNumber"`
						Polygon    []float64 `json:"polygon"`
					} `json:"boundingRegions"`
					Confidence float64 `json:"confidence"`
					Content    string  `json:"content"`
					Spans      []struct {
						Length float64 `json:"length"`
						Offset float64 `json:"offset"`
					} `json:"spans"`
					Type      string `json:"type"`
					ValueTime string `json:"valueTime"`
				} `json:"TransactionTime"`
			} `json:"fields"`
			Spans []struct {
				Length float64 `json:"length"`
				Offset float64 `json:"offset"`
			} `json:"spans"`
		} `json:"documents"`
		ModelID string `json:"modelId"`
		Pages   []struct {
			Angle  float64 `json:"angle"`
			Height float64 `json:"height"`
			Lines  []struct {
				Content string    `json:"content"`
				Polygon []float64 `json:"polygon"`
				Spans   []struct {
					Length float64 `json:"length"`
					Offset float64 `json:"offset"`
				} `json:"spans"`
			} `json:"lines"`
			PageNumber float64 `json:"pageNumber"`
			Spans      []struct {
				Length float64 `json:"length"`
				Offset float64 `json:"offset"`
			} `json:"spans"`
			Unit  string  `json:"unit"`
			Width float64 `json:"width"`
			Words []struct {
				Confidence float64   `json:"confidence"`
				Content    string    `json:"content"`
				Polygon    []float64 `json:"polygon"`
				Span       struct {
					Length float64 `json:"length"`
					Offset float64 `json:"offset"`
				} `json:"span"`
			} `json:"words"`
		} `json:"pages"`
		StringIndexType string        `json:"stringIndexType"`
		Styles          []interface{} `json:"styles"`
	} `json:"analyzeResult"`
	CreatedDateTime     string `json:"createdDateTime"`
	LastUpdatedDateTime string `json:"lastUpdatedDateTime"`
	Status              string `json:"status"`
}

type (
	// DueDate struct
	DueDate struct {
		BoundingRegions []struct {
			PageNumber float64   `json:"pageNumber"`
			Polygon    []float64 `json:"polygon"`
		} `json:"boundingRegions"`
		Confidence float64 `json:"confidence"`
		Content    string  `json:"content"`
		Spans      []struct {
			Length float64 `json:"length"`
			Offset float64 `json:"offset"`
		} `json:"spans"`
		Type      string `json:"type"`
		ValueDate string `json:"valueDate"`
	}

	// PaymentTerm struct
	PaymentTerm struct {
		BoundingRegions []struct {
			PageNumber float64   `json:"pageNumber"`
			Polygon    []float64 `json:"polygon"`
		} `json:"boundingRegions"`
		Confidence float64 `json:"confidence"`
		Content    string  `json:"content"`
		Spans      []struct {
			Length float64 `json:"length"`
			Offset float64 `json:"offset"`
		} `json:"spans"`
		Type        string `json:"type"`
		ValueString string `json:"valueString"`
	}
)
