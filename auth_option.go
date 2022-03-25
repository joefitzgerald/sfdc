package sfdc

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

type AuthOption interface {
	applyAuth(i *Instance) error
}

type withNoAuthentication struct{}

func (w *withNoAuthentication) applyAuth(i *Instance) error { return nil }

// WithNoAuthentication specifies that you are handling authentication yourself. You should consider using WithHTTPClient to provide your authenticated HTTPClient.
func WithNoAuthentication() AuthOption { return &withNoAuthentication{} }

// WithToken is an AuthOption that sets the token to use for authentication
func WithToken(ctx context.Context, config *oauth2.Config, token *oauth2.Token) AuthOption {
	return &withToken{
		ctx:    ctx,
		token:  token,
		config: config,
	}
}

func (w *withToken) applyAuth(i *Instance) error {
	if _, ok := w.ctx.Value(oauth2.HTTPClient).(*http.Client); !ok {
		w.ctx = context.WithValue(w.ctx, oauth2.HTTPClient, i.client)
	}
	if !w.token.Valid() {
		s := w.config.TokenSource(w.ctx, w.token)
		t, err := s.Token()
		if err != nil {
			return err
		}
		w.token = t
	}
	i.client = w.config.Client(w.ctx, w.token)
	instanceURL, ok := w.token.Extra("instance_url").(string)
	if !ok {
		return errors.New("instance_url not available in the token")
	}
	i.url = instanceURL
	return nil
}

type withToken struct {
	ctx    context.Context
	token  *oauth2.Token
	config *oauth2.Config
}
