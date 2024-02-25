package ifood

import "time"

type (
	Orders []*Order
	Order  struct {
		ID         string    `json:"id"`
		ShortID    string    `json:"shortId"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  string    `json:"updatedAt"`
		ClosedAt   string    `json:"closedAt"`
		LastStatus string    `json:"lastStatus"`
		Details    struct {
			Mode        string `json:"mode"`
			Scheduled   bool   `json:"scheduled"`
			Tippable    bool   `json:"tippable"`
			Trackable   bool   `json:"trackable"`
			Boxable     bool   `json:"boxable"`
			PlacedAtBox bool   `json:"placedAtBox"`
			Reviewed    bool   `json:"reviewed"`
			DarkKitchen bool   `json:"darkKitchen"`
		} `json:"details"`
		Delivery struct {
			Address struct {
				City         string `json:"city"`
				Country      string `json:"country"`
				Neighborhood string `json:"neighborhood"`
				State        string `json:"state"`
				StreetName   string `json:"streetName"`
				StreetNumber string `json:"streetNumber"`
				Coordinates  struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"coordinates"`
				Reference  string `json:"reference"`
				Complement string `json:"complement"`
			} `json:"address"`
			EstimatedTimeOfArrival struct {
				DeliversAt string `json:"deliversAt"`
				UpdatedAt  string `json:"updatedAt"`
			} `json:"estimatedTimeOfArrival"`
			ExpectedDeliveryTime string `json:"expectedDeliveryTime"`
			ExpectedDuration     int    `json:"expectedDuration"`
		} `json:"delivery"`
		Merchant struct {
			Address struct {
				Establishment string `json:"establishment"`
				City          string `json:"city"`
				Country       string `json:"country"`
				Neighborhood  string `json:"neighborhood"`
				State         string `json:"state"`
				StreetName    string `json:"streetName"`
				StreetNumber  string `json:"streetNumber"`
				Coordinates   struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"coordinates"`
				Reference  string `json:"reference"`
				Complement string `json:"complement"`
			} `json:"address"`
			ID           string `json:"id"`
			Name         string `json:"name"`
			PhoneNumber  string `json:"phoneNumber"`
			Logo         string `json:"logo"`
			CompanyGroup string `json:"companyGroup"`
			Type         string `json:"type"`
		} `json:"merchant"`
		Payments struct {
			Methods []struct {
				ID    string `json:"id"`
				Brand struct {
					ID          string `json:"id"`
					Image       string `json:"image"`
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"brand"`
				Credit struct {
					CardNumber string `json:"cardNumber"`
				} `json:"credit"`
				Amount struct {
					Currency string `json:"currency"`
					Value    int    `json:"value"`
				} `json:"amount"`
				Transactions []struct {
					ID        string        `json:"id"`
					Type      string        `json:"type"`
					Status    string        `json:"status"`
					CreatedAt string        `json:"createdAt"`
					Value     int           `json:"value"`
					Refunds   []interface{} `json:"refunds"`
				} `json:"transactions"`
			} `json:"methods"`
			Total struct {
				Currency string `json:"currency"`
				Value    int    `json:"value"`
			} `json:"total"`
		} `json:"payments"`
		Bag struct {
			Benefits    []interface{} `json:"benefits"`
			DeliveryFee struct {
				Value             int `json:"value"`
				ValueWithDiscount int `json:"valueWithDiscount"`
			} `json:"deliveryFee"`
			Items []struct {
				ExternalID             string        `json:"externalId"`
				Name                   string        `json:"name"`
				Quantity               int           `json:"quantity"`
				SubItems               []interface{} `json:"subItems"`
				Tags                   []interface{} `json:"tags"`
				TotalPrice             int           `json:"totalPrice"`
				TotalPriceWithDiscount int           `json:"totalPriceWithDiscount"`
				UnitPrice              int           `json:"unitPrice"`
				UnitPriceWithDiscount  int           `json:"unitPriceWithDiscount"`
			} `json:"items"`
			SubTotal struct {
				Value             int `json:"value"`
				ValueWithDiscount int `json:"valueWithDiscount"`
			} `json:"subTotal"`
			Total struct {
				Value             int `json:"value"`
				ValueWithDiscount int `json:"valueWithDiscount"`
			} `json:"total"`
			Updated bool `json:"updated"`
		} `json:"bag"`
		Origin struct {
			Platform   string `json:"platform"`
			AppName    string `json:"appName"`
			AppVersion string `json:"appVersion"`
		} `json:"origin"`
		DeliveryMethod struct {
			ID   string `json:"id"`
			Mode string `json:"mode"`
		} `json:"deliveryMethod"`
	}
)

func (o Orders) FilterByStatus(status string) Orders {
	res := make([]*Order, 0, 0)

	for _, order := range o {
		if order.LastStatus == status {
			res = append(res, order)
		}
	}

	return res
}

func (o Orders) FilterFrom(from time.Time) Orders {
	res := make([]*Order, 0, 0)

	for _, order := range o {
		if order.CreatedAt.After(from) {
			res = append(res, order)
		}
	}

	return res
}

func (o Orders) FilterTo(to time.Time) Orders {
	res := make([]*Order, 0, 0)

	for _, order := range o {
		if order.CreatedAt.Before(to) {
			res = append(res, order)
		}
	}

	return res
}
