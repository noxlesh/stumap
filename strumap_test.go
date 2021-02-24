package strumap

import (
	"encoding/json"
	"reflect"
	"testing"
)

var a = struct {
	BlackFox   bool
	Name       string
	TestNilPtr *int
	TestPtr    *int
	LeastName  string
	Title      string
	Slice      []string
	Inside     struct {
		SomeOf   int
		MostName string
	}
}{
	BlackFox:   false,
	Name:       "abbr",
	TestPtr:    func() *int { a := 1; return &a }(),
	TestNilPtr: nil,
	LeastName:  "Washere",
	Title:      "Nothing there",
	Slice:      []string{"asf", "cde", "dsdsdsdosdsdosjsdklj"},
	Inside: struct {
		SomeOf   int
		MostName string
	}{
		SomeOf:   24565434,
		MostName: "Tohere"},
}

func TestConvert(t *testing.T) {
	type args struct {
		s interface{}
	}
	cases := []struct {
		name    string
		args    args
		wantMap map[string]interface{}
		wantErr error
	}{
		{
			name: "not a struct but nil",
			args: args{
				nil,
			},
			wantMap: nil,
			wantErr: ErrNotStruct,
		},
		{
			name: "not a struct but a string",
			args: args{
				s: "1",
			},
			wantMap: nil,
			wantErr: ErrNotStruct,
		},
		{
			name: "empty struct",
			args: args{
				s: struct{}{},
			},
			wantMap: nil,
			wantErr: ErrEmptyStruct,
		},
		{
			name: "struct with a string field",
			args: args{
				s: struct {
					Field string
				}{
					Field: "field",
				},
			},
			wantMap: map[string]interface{}{
				"Field": "field",
			},
			wantErr: nil,
		},
		{
			name: "struct with a string field. convert to snake case",
			args: args{
				s: struct {
					FieldOf string
				}{
					FieldOf: "field",
				},
			},
			wantMap: map[string]interface{}{
				"FieldOf": "field",
			},
			wantErr: nil,
		},
		{
			name: "struct with an unexported string field",
			args: args{
				s: struct {
					field string
				}{
					field: "field",
				},
			},
			wantMap: map[string]interface{}{},
			wantErr: nil,
		},
		{
			name: "struct has a struct field witch has a string filed",
			args: args{
				s: struct {
					Field struct {
						Str string
					}
				}{
					Field: struct {
						Str string
					}{Str: "field"},
				},
			},
			wantMap: map[string]interface{}{
				"Field": map[string]interface{}{
					"Str": "field",
				},
			},
			wantErr: nil,
		},
		{
			name: "struct has a pointer to a struct field witch has a string field",
			args: args{
				s: struct {
					Field *struct {
						Str string
					}
				}{
					Field: &struct {
						Str string
					}{Str: "field"},
				},
			},
			wantMap: map[string]interface{}{
				"Field": map[string]interface{}{
					"Str": "field",
				},
			},
			wantErr: nil,
		},
		{
			name: "struct has a pointer to a struct field witch has a nil pointer to int field",
			args: args{
				s: struct {
					Field *struct {
						Str *int
					}
				}{
					Field: &struct {
						Str *int
					}{Str: nil},
				},
			},
			wantMap: map[string]interface{}{
				"Field": map[string]interface{}{
					"Str": (*int)(nil),
				},
			},
			wantErr: nil,
		},
		{
			name: "most comprehensive",
			args: args{
				s: struct {
					TestPtr *int
					Abc     *string
					Slice   []string
				}{
					TestPtr: func() *int { a := 1; return &a }(),
					Abc:     func() *string { a := "abc"; return &a }(),
					Slice:   []string{"asf", "cde", "dsdsdsdosdsdosjsdklj"},
				},
			},
			wantMap: map[string]interface{}{
				"TestPtr": 1,
				"Abc":     "abc",
				"Slice":   []string{"asf", "cde", "dsdsdsdosdsdosjsdklj"},
			},
			wantErr: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotMap, gotErr := Convert(tc.args.s)
			if gotErr != tc.wantErr {
				t.Errorf("Convert() \n\t got error: %v \n\t want error: %v", gotErr, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMap, tc.wantMap) {
				t.Errorf("Convert() \n\t got = %v, \n\t want = %v", gotMap, tc.wantMap)
			}
		})
	}
}

func TestConvertSnakeCase(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantMap map[string]interface{}
		wantErr error
	}{

		{
			name: "struct with a string field. convert to camel case",
			args: args{
				s: struct {
					FieldOf string
				}{
					FieldOf: "field",
				},
			},
			wantMap: map[string]interface{}{
				"field_of": "field",
			},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMap, gotErr := ConvertSnakeCase(tc.args.s)
			if gotErr != tc.wantErr {
				t.Errorf("Convert() \n\t got error: %v \n\t want error: %v", gotErr, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMap, tc.wantMap) {
				t.Errorf("Convert() \n\t got = %v, \n\t want = %v", gotMap, tc.wantMap)
			}
		})
	}
}

func BenchmarkConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, _ := Convert(a)
		_ = d
	}
}

func BenchmarkConvertUsingJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(a)
		d := map[string]interface{}{}
		_ = json.Unmarshal(data, &d)
		_ = d
	}
}
