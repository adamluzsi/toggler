// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// EnrollmentResponseBody EnrollmentResponse will be returned when feature flag status is requested.
//
// The content will be always given, regardless if the flag exists or not.
// This helps the developers to use it as a null object, regardless the toggler service state.
// swagger:model EnrollmentResponseBody
type EnrollmentResponseBody struct {

	// Rollout feature flag enrollment status.
	Enrollment bool `json:"enrollment,omitempty"`
}

// Validate validates this enrollment response body
func (m *EnrollmentResponseBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *EnrollmentResponseBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EnrollmentResponseBody) UnmarshalBinary(b []byte) error {
	var res EnrollmentResponseBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
