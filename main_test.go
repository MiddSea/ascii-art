package main

import (
	"testing"
)

func TestProcessQuotesRegEx(t *testing.T) {
	tests := []struct {
		tstNo    int
		name     string
		input    string
		expected string
	}{
		{
			tstNo:    1,
			name:     "2 quotes at end simple and short",
			input:    "'<open close>'",
			expected: "'<open close>'",
		},
		{
			tstNo:    2,
			name:     "empty string",
			input:    "",
			expected: "", // does nothing - no error
		},
		{
			tstNo:    3,
			name:     "single quote",
			input:    "'",
			expected: "", // ERROR
		},
		{
			tstNo:    4,
			name:     "one quote begin short",
			input:    "'<Open CloseCXX ",
			expected: "", // ERROR
		},
		{
			tstNo:    5,
			name:     "one quote end short",
			input:    "XXOpen CloseC>'",
			expected: "", // ERROR
		},
		{
			tstNo:    6,
			name:     "no quote simple and short",
			input:    "just short text",
			expected: "just short text",
		},
		{
			tstNo:    7,
			name:     "extreme quote handling",
			input:    "don't even  ' <Open CloseC> ' and ' open won't close 'and' <open close> 'j",
			expected: "don't even '<Open CloseC>' and 'open won't close' and '<open close>' j",
		},
		{
			tstNo:    8,
			name:     "quote at beginning",
			input:    "'Now, don't even dare ' and ' open won't close 'and' <open close> '",
			expected: "'Now, don't even dare' and 'open won't close' and '<open close>' ",
		},
		{
			tstNo:    9,
			name:     "quote handling",
			input:    "They said: ' hello world '",
			expected: "They said: 'hello world'",
		},
		{
			tstNo:    10,
			name:     "two consecutive quotes",
			input:    "They said: ' hello world ''",
			expected: "They said: 'hello world'",
		},
		{
			tstNo:    11,
			name:     "additional test case",
			input:    "",
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := processQuotesRegEx(tt.input)
			if result != tt.expected {
				inpStr, procStr, wantStr := "input", "processQuotes()", "...want"
				t.Errorf("\n%14d %5s = %q\n%20s = %q\n%20s = %q",
							tt.tstNo, inpStr, tt.input,
										  procStr, result,
													wantStr, tt.expected)
			}
		})
	}
}


func TestProcessText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hex conversion",
			input:    "1E (hex) files were added",
			expected: "30 files were added",
		},
		{
			name:     "binary conversion",
			input:    "It has been 10 (bin) years",
			expected: "It has been 2 years",
		},
		{
			name:     "uppercase conversion",
			input:    "Ready, set, go (up) !",
			expected: "Ready, set, GO!",
		},
		{
			name:     "lowercase conversion",
			input:    "I should stop SHOUTING (low)",
			expected: "I should stop shouting",
		},
		{
			name:     "capitalize conversion",
			input:    "Welcome to the Brooklyn bridge (cap)",
			expected: "Welcome to the Brooklyn Bridge",
		},
		{
			name:     "multiple word transformation",
			input:    "This is so exciting (up, 2)",
			expected: "This is SO EXCITING",
		},
		{
			name:     "punctuation spacing",
			input:    "Hello ,world ! How are you ?",
			expected: "Hello, world! How are you?",
		},
		{
			name:     "ellipsis handling",
			input:    "I was thinking ... You were right",
			expected: "I was thinking... You were right",
		},
		{
			name:     "ellipsis handling lonely periods",
			input:    "I was thinking . .. You were.  . .right",
			expected: "I was thinking... You were... right",
		},
		{
			name:     "extreme quote handling",
			input:    "don't even  ' <Open CloseC> ' and ' open won't close 'and' <open close> '",
			expected: "don't even '<Open CloseC>' and 'open won't close' and '<open close>' ",
		},
		{
			name:     "quote at begining handling",
			input:    "' Now, don't even dare ' and ' open won't close 'and' <open close> '",
			expected: "'Now, don't even dare' and 'open won't close' and '<open close>' ",
		},
		{
			name:     "quote handling",
			input:    "They said: ' hello world '",
			expected: "They said: 'hello world'",
		},
		{
			name:     "article fixing",
			input:    "This is a excellent day",
			expected: "This is an excellent day",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := processText(tt.input)
			if result != tt.expected {
				t.Errorf("       Input = %q\n             processText() = %q\n                   ...want = %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessHex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1E (hex)", "30"},
		{"FF (hex)", "255"},
		{"A5 (hex)", "165"},
	}

	for _, tt := range tests {
		result, _ := processHex(tt.input)
		if result != tt.expected {
			t.Errorf("processHex(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestProcessBin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1010 (bin)", "10"},
		{"1100100 (bin)", "100"},
		{"11111111 (bin)", "255"},
	}

	for _, tt := range tests {
		result, _ := processBin(tt.input)
		if result != tt.expected {
			t.Errorf("processBin(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
