package valuemodifier

import (
	"fmt"
	"testing"
)

type testModifier1 struct{}

func (m testModifier1) Modify(value []byte) ([]byte, error) {
	return []byte(string(value) + "-modified1"), nil
}

type testModifier2 struct{}

func (m testModifier2) Modify(value []byte) ([]byte, error) {
	return []byte(string(value) + "-modified2"), nil
}

func Test_ValueModifier_TraverseJSON(t *testing.T) {
	testCases := []struct {
		ValueModifiers []ValueModifier
		IgnoreFields   []string
		Input          string
		Expected       string
	}{
		// Test case 1, a single modifier modifies all secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `{
  "noSecret1": "noSecret1",
  "pass1": "pass1"
}`,
			Expected: `{
  "noSecret1": "noSecret1-modified1",
  "pass1": "pass1-modified1"
}`,
		},
		// Test case 2, a single modifier modifies all numeric secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `{
  "noSecret1": "noSecret1",
  "pass1": 12345
}`,
			Expected: `{
  "noSecret1": "noSecret1-modified1",
  "pass1": "12345-modified1"
}`,
		},
		// Test case 3, a single modifier modifies all secrets inside lists.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `{
  "list1": [
    {
      "pass1": "pass1"
    }
  ]
}`,
			Expected: `{
  "list1": [
    {
      "pass1": "pass1-modified1"
    }
  ]
}`,
		},
		// Test case 4, a single modifier modifies all secrets, but ignores the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{
				"noSecret1",
			},
			Input: `{
  "noSecret1": "foo",
  "pass1": "pass1"
}`,
			Expected: `{
  "noSecret1": "foo",
  "pass1": "pass1-modified1"
}`,
		},
		// Test case 5, multiple modifiers modify all secrets, but ignore the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			Input: `{
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass1": "pass1",
  "pass2": "pass2"
}`,
			Expected: `{
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass1": "pass1-modified1-modified2",
  "pass2": "pass2-modified1-modified2"
}`,
		},
		// Test case 6, nested blocks, multiple modifiers modify all secrets, but
		// ignore the ones configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			Input: `{
  "block1": {
    "block11": {
      "pass1": "pass1"
    },
    "pass2": "pass2"
  },
  "block2": {
    "block21": {
      "pass3": "pass3"
    },
    "pass4": "pass4"
  },
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass5": "pass5",
  "pass6": 123456
}`,
			Expected: `{
  "block1": {
    "block11": {
      "pass1": "pass1-modified1-modified2"
    },
    "pass2": "pass2-modified1-modified2"
  },
  "block2": {
    "block21": {
      "pass3": "pass3-modified1-modified2"
    },
    "pass4": "pass4-modified1-modified2"
  },
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass5": "pass5-modified1-modified2",
  "pass6": "123456-modified1-modified2"
}`,
		},
	}

	for i, testCase := range testCases {
		config := DefaultConfig()
		config.ValueModifiers = testCase.ValueModifiers
		config.IgnoreFields = testCase.IgnoreFields
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.TraverseJSON([]byte(testCase.Input))
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if string(output) != testCase.Expected {
			t.Fatal("test", i+1, "expected", fmt.Sprintf("%q", testCase.Expected), "got", fmt.Sprintf("%q", output))
		}
	}
}

func Test_ValueModifier_TraverseYAML(t *testing.T) {
	testCases := []struct {
		ValueModifiers []ValueModifier
		IgnoreFields   []string
		Input          string
		Expected       string
	}{
		// Test case 1, a single modifier modifies all secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `noSecret1: noSecret1
pass1: pass1
`,
			Expected: `noSecret1: noSecret1-modified1
pass1: pass1-modified1
`,
		},
		// Test case 2, a single modifier modifies all numeric secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `noSecret1: noSecret1
pass1: 12345
`,
			Expected: `noSecret1: noSecret1-modified1
pass1: 12345-modified1
`,
		},
		// Test case 3, a single modifier modifies all secrets inside lists.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `list1:
- pass1: pass1
`,
			Expected: `list1:
- pass1: pass1-modified1
`,
		},
		// Test case 4, a single modifier modifies all secrets, but ignores the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{
				"noSecret1",
			},
			Input: `noSecret1: noSecret1
pass1: pass1
`,
			Expected: `noSecret1: noSecret1
pass1: pass1-modified1
`,
		},
		// Test case 5, multiple modifiers modify all secrets, but ignore the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			Input: `noSecret1: noSecret1
noSecret2: noSecret2
pass1: pass1
pass2: pass2
`,
			Expected: `noSecret1: noSecret1
noSecret2: noSecret2
pass1: pass1-modified1-modified2
pass2: pass2-modified1-modified2
`,
		},
		// Test case 6, nested blocks, multiple modifiers modify all secrets, but
		// ignore the ones configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			Input: `block1:
  block11:
    pass1: pass1
  pass2: pass2
block2:
  block21:
    pass3: pass3
  pass4: pass4
noSecret1: foo
noSecret2: bar
pass5: pass5
pass6: 123456
`,
			Expected: `block1:
  block11:
    pass1: pass1-modified1-modified2
  pass2: pass2-modified1-modified2
block2:
  block21:
    pass3: pass3-modified1-modified2
  pass4: pass4-modified1-modified2
noSecret1: foo
noSecret2: bar
pass5: pass5-modified1-modified2
pass6: 123456-modified1-modified2
`,
		},
		// Test case 7, modifiers modify string blocks.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `pass1: |
  foo
  bar
`,
			Expected: `pass1: |-
  foo
  bar
  -modified1
`,
		},
		// Test case 8, modifiers modify secrets of string blocks representing YAML.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `pass1: |
  bar:
    baz: pass2
  foo: pass3
`,
			Expected: `pass1: |
  bar:
    baz: pass2-modified1
  foo: pass3-modified1
`,
		},
		// Test case 9, modifiers modify secrets of string blocks representing JSON.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `pass1: |
  {
    "block1": {
      "block11": {
        "pass2": "pass2"
      }
    }
  }
`,
			Expected: `pass1: |-
  {
    "block1": {
      "block11": {
        "pass2": "pass2-modified1"
      }
    }
  }
`,
		},
	}

	for i, testCase := range testCases {
		if i != 8 {
			continue
		}
		config := DefaultConfig()
		config.ValueModifiers = testCase.ValueModifiers
		config.IgnoreFields = testCase.IgnoreFields
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.TraverseYAML([]byte(testCase.Input))
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if string(output) != testCase.Expected {
			t.Fatal("test", i+1, "expected", fmt.Sprintf("%q", testCase.Expected), "got", fmt.Sprintf("%q", output))
		}
	}
}
