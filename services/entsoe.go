package services

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"spotit-backend/config"
	"time"
)

type PublicationMarketDocument struct {
	XMLName                     xml.Name `xml:"Publication_MarketDocument"`
	Text                        string   `xml:",chardata"`
	Xmlns                       string   `xml:"xmlns,attr"`
	MRID                        string   `xml:"mRID"`
	RevisionNumber              string   `xml:"revisionNumber"`
	Type                        string   `xml:"type"`
	SenderMarketParticipantMRID struct {
		Text         string `xml:",chardata"`
		CodingScheme string `xml:"codingScheme,attr"`
	} `xml:"sender_MarketParticipant.mRID"`
	SenderMarketParticipantMarketRoleType string `xml:"sender_MarketParticipant.marketRole.type"`
	ReceiverMarketParticipantMRID         struct {
		Text         string `xml:",chardata"`
		CodingScheme string `xml:"codingScheme,attr"`
	} `xml:"receiver_MarketParticipant.mRID"`
	ReceiverMarketParticipantMarketRoleType string `xml:"receiver_MarketParticipant.marketRole.type"`
	CreatedDateTime                         string `xml:"createdDateTime"`
	PeriodTimeInterval                      struct {
		Text  string `xml:",chardata"`
		Start string `xml:"start"`
		End   string `xml:"end"`
	} `xml:"period.timeInterval"`
	TimeSeries struct {
		Text         string `xml:",chardata"`
		MRID         string `xml:"mRID"`
		BusinessType string `xml:"businessType"`
		InDomainMRID struct {
			Text         string `xml:",chardata"`
			CodingScheme string `xml:"codingScheme,attr"`
		} `xml:"in_Domain.mRID"`
		OutDomainMRID struct {
			Text         string `xml:",chardata"`
			CodingScheme string `xml:"codingScheme,attr"`
		} `xml:"out_Domain.mRID"`
		CurrencyUnitName     string `xml:"currency_Unit.name"`
		PriceMeasureUnitName string `xml:"price_Measure_Unit.name"`
		CurveType            string `xml:"curveType"`
		Period               struct {
			Text         string `xml:",chardata"`
			TimeInterval struct {
				Text  string `xml:",chardata"`
				Start string `xml:"start"`
				End   string `xml:"end"`
			} `xml:"timeInterval"`
			Resolution string `xml:"resolution"`
			Point      []struct {
				Text        string  `xml:",chardata"`
				Position    int     `xml:"position"`
				PriceAmount float32 `xml:"price.amount"`
			} `xml:"Point"`
		} `xml:"Period"`
	} `xml:"TimeSeries"`
}

func getEntsoeURL() string {
	time := time.Now().Format(config.GetEntsoeDateFormat()) + "0100"
	return config.GetEntsoeBase() + fmt.Sprintf("&periodStart=%v&periodEnd=%v", time, time)
}

func convertPrice(MWh float32) float32 {
	return (MWh / 10)
}

func parseXML(xmlBytes []byte) (*map[int]float32, error) {
	var e PublicationMarketDocument
	err := xml.Unmarshal(xmlBytes, &e)
	if err != nil {
		return nil, err
	}

	priceSeries := make(map[int]float32)
	for _, e := range e.TimeSeries.Period.Point {
		priceSeries[e.Position] = convertPrice(e.PriceAmount)
	}
	return &priceSeries, nil
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func GetDayAhead() (*map[int]float32, error) {
	xmlBytes, err := getXML(getEntsoeURL())
	if err != nil {
		return nil, err
	}
	return parseXML(xmlBytes)
}
