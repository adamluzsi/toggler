// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/toggler-io/toggler/lib/go/models"
)

// NewClientConfigPOSTParams creates a new ClientConfigPOSTParams object
// with the default values initialized.
func NewClientConfigPOSTParams() *ClientConfigPOSTParams {
	var ()
	return &ClientConfigPOSTParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewClientConfigPOSTParamsWithTimeout creates a new ClientConfigPOSTParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewClientConfigPOSTParamsWithTimeout(timeout time.Duration) *ClientConfigPOSTParams {
	var ()
	return &ClientConfigPOSTParams{

		timeout: timeout,
	}
}

// NewClientConfigPOSTParamsWithContext creates a new ClientConfigPOSTParams object
// with the default values initialized, and the ability to set a context for a request
func NewClientConfigPOSTParamsWithContext(ctx context.Context) *ClientConfigPOSTParams {
	var ()
	return &ClientConfigPOSTParams{

		Context: ctx,
	}
}

// NewClientConfigPOSTParamsWithHTTPClient creates a new ClientConfigPOSTParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewClientConfigPOSTParamsWithHTTPClient(client *http.Client) *ClientConfigPOSTParams {
	var ()
	return &ClientConfigPOSTParams{
		HTTPClient: client,
	}
}

/*ClientConfigPOSTParams contains all the parameters to send to the API endpoint
for the client config p o s t operation typically these are written to a http.Request
*/
type ClientConfigPOSTParams struct {

	/*Body*/
	Body *models.ClientConfigRequestInput

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the client config p o s t params
func (o *ClientConfigPOSTParams) WithTimeout(timeout time.Duration) *ClientConfigPOSTParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the client config p o s t params
func (o *ClientConfigPOSTParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the client config p o s t params
func (o *ClientConfigPOSTParams) WithContext(ctx context.Context) *ClientConfigPOSTParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the client config p o s t params
func (o *ClientConfigPOSTParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the client config p o s t params
func (o *ClientConfigPOSTParams) WithHTTPClient(client *http.Client) *ClientConfigPOSTParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the client config p o s t params
func (o *ClientConfigPOSTParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the client config p o s t params
func (o *ClientConfigPOSTParams) WithBody(body *models.ClientConfigRequestInput) *ClientConfigPOSTParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the client config p o s t params
func (o *ClientConfigPOSTParams) SetBody(body *models.ClientConfigRequestInput) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *ClientConfigPOSTParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}