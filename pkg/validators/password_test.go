package validators

import "testing"

func TestPasswordValidator_Validate(t *testing.T) {
	type fields struct {
		MinLength  int
		MinNumeric int
		MinSpecial int
		MinUpper   int
		MinLower   int
	}

	tests := []struct {
		name     string
		fields   fields
		password string
		wantErr  bool
	}{
		{name: "No constraint", fields: fields{}, password: "", wantErr: false},
		{name: "Min len success", fields: fields{MinLength: 6}, password: "123456", wantErr: false},
		{name: "Min len failure", fields: fields{MinLength: 6}, password: "12345", wantErr: true},
		{name: "Min special failure", fields: fields{MinSpecial: 1}, password: "12345", wantErr: true},
		{name: "Min special success", fields: fields{MinSpecial: 1}, password: "123@45", wantErr: false},
		{name: "Min Numeric success", fields: fields{MinNumeric: 1}, password: "asdf1g", wantErr: false},
		{name: "Min Numeric failure", fields: fields{MinNumeric: 1}, password: "asdfg", wantErr: true},
		{name: "Min Upper failure", fields: fields{MinUpper: 1}, password: "asdfg", wantErr: true},
		{name: "Min Upper success", fields: fields{MinUpper: 1}, password: "aSdfg", wantErr: false},
		{name: "Min Upper success", fields: fields{MinLower: 1}, password: "SSSSd", wantErr: false},
		{name: "Min Upper failure", fields: fields{MinLower: 1}, password: "SSSS", wantErr: true},
		{name: "Strong password success", fields: fields{MinUpper: 1, MinNumeric: 1, MinSpecial: 1, MinLength: 8}, password: "pAssw@rd5", wantErr: false},
		{name: "Strong password failure", fields: fields{MinUpper: 1, MinNumeric: 1, MinSpecial: 1, MinLength: 8}, password: "pAssword5", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := PasswordValidator{
				MinLength:  tt.fields.MinLength,
				MinNumeric: tt.fields.MinNumeric,
				MinSpecial: tt.fields.MinSpecial,
				MinUpper:   tt.fields.MinUpper,
				MinLower:   tt.fields.MinLower,
			}
			if err := v.Validate(tt.password); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
