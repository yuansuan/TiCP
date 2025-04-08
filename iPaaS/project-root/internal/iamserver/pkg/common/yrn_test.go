package common

import (
	"testing"
)

func TestParseYRN(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *YRN
		err      bool
	}{
		{
			name:  "Valid YRN",
			input: "yrn:ys:service:region:account-id:resource-type/resource-id",
			expected: &YRN{
				Partition:    "ys",
				Service:      "service",
				Region:       "region",
				AccountID:    "account-id",
				ResourceType: "resource-type",
				ResourceID:   "resource-id",
			},
			err: false,
		},
		{
			name:  "Valid YRN - windows format",
			input: "yrn:ys:service:region:account-id:resource-type/D:\\LD\\2024\\17_NX",
			expected: &YRN{
				Partition:    "ys",
				Service:      "service",
				Region:       "region",
				AccountID:    "account-id",
				ResourceType: "resource-type",
				ResourceID:   "D:\\LD\\2024\\17_NX",
			},
			err: false,
		},
		{
			name:     "Invalid YRN - wrong partition",
			input:    "yrn:wrong-partition:service:region:account-id:resource-type/resource-id",
			expected: nil,
			err:      true,
		},
		{
			name:     "Invalid YRN - missing resource type",
			input:    "yrn:ys:service:region:account-id:/resource-id",
			expected: nil,
			err:      true,
		},
		{
			name:     "Invalid YRN - missing resource id",
			input:    "yrn:ys:service:region:account-id:resource-type/",
			expected: nil,
			err:      true,
		},
		{
			name:     "Invalid YRN - missing partition",
			input:    "yrn::service:region:account-id:resource-type/resource-id",
			expected: nil,
			err:      true,
		},
		{
			name:     "Invalid YRN - missing fields",
			input:    "yrn:ys:service:region:account-id",
			expected: nil,
			err:      true,
		},
		{
			name:     "Invalid YRN - wrong format",
			input:    "invalid-yrn-format",
			expected: nil,
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseYRN(tc.input)

			if tc.err && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.err && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if tc.expected == nil && actual != nil {
				t.Errorf("Expected nil but got %v", actual)
			}

			if tc.expected != nil && actual == nil {
				t.Errorf("Expected %v but got nil", tc.expected)
			}

			if tc.expected != nil && actual != nil {
				if *tc.expected != *actual {
					t.Errorf("Expected %v but got %v", tc.expected, actual)
				}
			}
		})
	}
}
