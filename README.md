# `sfdc`

A `go` client library for the Salesforce REST API.

## Usage

### Get the Package

```shell
go get -u github.com/joefitzgerald/sfdc
```

**Note:** this package requires Go 1.18 or later.

### Connect to the Salesforce API

⚠️ You should make sure you have completed the steps outlined in [Getting Started with the Salesforce API](#getting-started-with-the-salesforce-api) section.

```go
package main

func main() {
	// First: Fetch a Token
	// Note: if you request the `refresh_token` scope, this token includes a 
	// refresh token, which can be used to fetch new tokens in the future without
	// re-authenticating the user.
	// https://help.salesforce.com/s/articleView?id=sf.remoteaccess_oauth_refresh_token_flow.htm
	clientID := "your-client-id"
	clientSecret := "your-client-secret"
	domain := "your-domain"
	scopes := []string{"api", "openid", "id", "profile", "email", "refresh_token"}
	config := devicecode.NewWithDomain(domain, clientID, clientSecret, scopes)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Minute)
	defer cancel()
	token, err := config.Token(ctx, http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	// Next: Use the token to construct a sfdc.Instance and then an sfdc.Entity
	// for each SObject you wish to query or interact with.
	instance, err := sfdc.New(sfdc.WithToken(context.Background(), config.Config, token), sfdc.WithURL(fmt.Sprintf("https://%s.my.salesforce.com", domain)))
	type Opportunity struct {
		ID          string `json:"Id,omitempty"`
		Name        string `json:"Name,omitempty"`
		Description string `json:"Description,omitempty"`
	}
	opportunityEntity := sfdc.NewEntity[Opportunity](instance)

	// Finally: Make requests to the Salesforce API
	opportunities, err := opportunityEntity.Query(context.Background(), "SELECT Id, Name, Description FROM Opportunity where CloseDate = 2019-01-01")
	if err != nil {
		log.Fatal(err)
	}
	for i := range opportunities {
		fmt.Printf("%s: %s\n", opportunities[i].ID, opportunities[i].Name)
	}
}
```

## Getting Started with the Salesforce API

Accessing the Salesforce REST API is straightforward, but requires some one-time preparation:

1. Sign up for Salesforce `Developer Edition`
2. Create a `Connected App`
3. [Access the Salesforce REST API](#usage)

### Step 1: Sign up for Salesforce `Developer Edition`

You will need a developer edition organization so that you can register a connected app. If your company already has organization(s) that they use for development, just validate the `API Enabled` permission, below.

1. Go to https://developer.salesforce.com/signup
2. Follow the instructions to create an organization
3. Verify that your user profile has the `API Enabled` permission set (this is enabled by default, but an administrator can modify it)

### Step 2: Create a `Connected App`

You need a `Connected App` so that you can get an OAuth 2.0 `Client ID` and `Client Secret`, and configure the `scopes` that your API client will be able to request and make use of.

1. In your Developer Edition organization, select “Setup”, and then go to `Platform Tools` > `Apps` > `Apps Manager`.
1. Select `New Connected App`.
1. Fill in the required fields, and then check the `Enable OAuth Settings` option. In the resulting section:
    * Check the `Enable for Device Flow` option
    * Add the following scopes:
      * Access the identity URL service (id, profile, email, address, phone): optional, if you want to be able to identify the user by name
      * Access unique user identifiers (openid): required for an OpenID connect payload in your token
      * Manage user data via APIs (api): required for all API access
      * Perform requests at any time (refresh_token, offline_access): required for you to receive a refresh token
    * Check the `Configure ID Token` option:
      * Set the `Token Valid for` option to the desired number of minutes
      * Set the `Include Standard Claims` option
1. Select `Save` and then `Continue`
1. Note the Client ID and Client Secret that you will make use of with this API client:
    * Client ID: Copy the `Consumer Key` field and use it as your `ClientID`
    * Client Secret: Copy the `Consumer Secret` field and make use of it as your `ClientSecret`