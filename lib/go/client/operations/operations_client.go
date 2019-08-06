// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
CreateRolloutFeatureFlag creates rollout feature flag

This operation allows you to create a new rollout feature flag.
*/
func (a *Client) CreateRolloutFeatureFlag(params *CreateRolloutFeatureFlagParams) (*CreateRolloutFeatureFlagOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateRolloutFeatureFlagParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "CreateRolloutFeatureFlag",
		Method:             "POST",
		PathPattern:        "/rollout/flag/create.json",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &CreateRolloutFeatureFlagReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*CreateRolloutFeatureFlagOK), nil

}

/*
IsFeatureEnabled checks rollout feature status for pilot

Reply back whether the feature for a given pilot id is enabled or not.
By Default, this will be determined whether the flag exist,
the pseudo random dice roll enrolls the pilot,
or if there any manually set enrollment status for the pilot.
The endpoint can be called with HTTP GET method as well,
POST is used officially only to support most highly abstracted http clients.
*/
func (a *Client) IsFeatureEnabled(params *IsFeatureEnabledParams) (*IsFeatureEnabledOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewIsFeatureEnabledParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "IsFeatureEnabled",
		Method:             "POST",
		PathPattern:        "/rollout/is-feature-enabled.json",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &IsFeatureEnabledReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*IsFeatureEnabledOK), nil

}

/*
IsFeatureGloballyEnabled checks rollout feature status for global use

Reply back whether the feature rolled out globally or not.
This is especially useful for cases where you don't have pilot id.
Such case is batch processing, or dark launch flips.
By Default, this will be determined whether the flag exist,
Then  whether the release id done to everyone or not by percentage.
The endpoint can be called with HTTP GET method as well,
POST is used officially only to support most highly abstracted http clients.
*/
func (a *Client) IsFeatureGloballyEnabled(params *IsFeatureGloballyEnabledParams) (*IsFeatureGloballyEnabledOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewIsFeatureGloballyEnabledParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "IsFeatureGloballyEnabled",
		Method:             "POST",
		PathPattern:        "/rollout/is-feature-globally-enabled.json",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &IsFeatureGloballyEnabledReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*IsFeatureGloballyEnabledOK), nil

}

/*
RolloutClientConfig checks multiple rollout feature status for a certain pilot

Return all the flag states that was requested in the favor of a Pilot.
This endpoint especially useful for Mobile & SPA apps.
The endpoint can be called with HTTP GET method as well,
POST is used officially only to support most highly abstracted http clients.
*/
func (a *Client) RolloutClientConfig(params *RolloutClientConfigParams) (*RolloutClientConfigOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRolloutClientConfigParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "RolloutClientConfig",
		Method:             "POST",
		PathPattern:        "/rollout/config.json",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &RolloutClientConfigReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*RolloutClientConfigOK), nil

}

/*
Websocket sockets API to check rollout feature flag status

This endpoint currently meant to used by servers and not by clients.
The  reason behind is that it is much more easy to calculate with server quantity,
than with client quantity, and therefore the load balancing is much more deterministic for the service.
The websocket based communication allows for servers to do low latency quick requests,
which is ideal to check flag status for individual requests that the server receives.
Because the nature of the persistent connection, TCP connection overhead is minimal.
The endpoint able to serve back whether the feature for a given pilot id is enabled or not.
The endpoint also able to serve back global flag state checks as well.
The flag enrollment interpretation use the same logic as it is described in the documentation.
*/
func (a *Client) Websocket(params *WebsocketParams, authInfo runtime.ClientAuthInfoWriter) (*WebsocketOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewWebsocketParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Websocket",
		Method:             "GET",
		PathPattern:        "/ws",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &WebsocketReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*WebsocketOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
