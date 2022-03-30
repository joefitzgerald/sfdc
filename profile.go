package sfdc

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/oauth2"
)

type Profile struct {
	ID             string `json:"id,omitempty"`
	AssertedUser   bool   `json:"asserted_user,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Username       string `json:"username,omitempty"`
	Nickname       string `json:"nick_name,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	Email          string `json:"email,omitempty"`
	EmailVerified  bool   `json:"email_verified,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	Photos         struct {
		Picture   string `json:"picture,omitempty"`
		Thumbnail string `json:"thumbnail,omitempty"`
	} `json:"photos,omitempty"`
	AddrStreet           string      `json:"addr_street,omitempty"`
	AddrCity             string      `json:"addr_city,omitempty"`
	AddrState            string      `json:"addr_state,omitempty"`
	AddrCountry          string      `json:"addr_country,omitempty"`
	AddrZip              interface{} `json:"addr_zip,omitempty"`
	MobilePhone          interface{} `json:"mobile_phone,omitempty"`
	MobilePhoneVerified  bool        `json:"mobile_phone_verified,omitempty"`
	IsLightningLoginUser bool        `json:"is_lightning_login_user,omitempty"`
	Status               struct {
		CreatedDate string `json:"created_date,omitempty"`
		Body        string `json:"body,omitempty"`
	} `json:"status,omitempty"`
	URLs struct {
		Enterprise   string `json:"enterprise,omitempty"`
		Metadata     string `json:"metadata,omitempty"`
		Partner      string `json:"partner,omitempty"`
		Rest         string `json:"rest,omitempty"`
		Sobjects     string `json:"sobjects,omitempty"`
		Search       string `json:"search,omitempty"`
		Query        string `json:"query,omitempty"`
		Recent       string `json:"recent,omitempty"`
		ToolingSoap  string `json:"tooling_soap,omitempty"`
		ToolingRest  string `json:"tooling_rest,omitempty"`
		Profile      string `json:"profile,omitempty"`
		Feeds        string `json:"feeds,omitempty"`
		Groups       string `json:"groups,omitempty"`
		Users        string `json:"users,omitempty"`
		FeedItems    string `json:"feed_items,omitempty"`
		FeedElements string `json:"feed_elements,omitempty"`
		CustomDomain string `json:"custom_domain,omitempty"`
	} `json:"urls,omitempty"`
	Active           bool      `json:"active,omitempty"`
	UserType         string    `json:"user_type,omitempty"`
	Language         string    `json:"language,omitempty"`
	Locale           string    `json:"locale,omitempty"`
	UtcOffset        int       `json:"utcOffset,omitempty"`
	LastModifiedDate time.Time `json:"last_modified_date,omitempty"`
}

func GetProfile(ctx context.Context, c *oauth2.Config, t *oauth2.Token) (*Profile, error) {
	id, ok := t.Extra("id").(string)
	if !ok {
		return nil, errors.New("id not available in the token")
	}

	client := c.Client(ctx, t)
	r, err := client.Get(id)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var result Profile
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
