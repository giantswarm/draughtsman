[![CircleCI](https://circleci.com/gh/giantswarm/valuemodifier.svg?&style=shield&circle-token=2f991b43e9c0147771ebe0575885501a7184dea9)](https://circleci.com/gh/giantswarm/valuemodifier)

# valuemodifier
Package valuemodifier provides an interface to modify values of arbitrary
structures in custom ways. Currently JSON and YAML formats are supported.

### usage
This is an working example of how to use this package. For more examples check
the go tests.
```
package main

import (
	"fmt"

	"github.com/giantswarm/valuemodifier"
	encodemodifier "github.com/giantswarm/valuemodifier/base64/encode"
	encryptmodifier "github.com/giantswarm/valuemodifier/gpg/encrypt"
)

func main() {
	var err error

	var encodeModifier valuemodifier.ValueModifier
	{
		modifierConfig := encodemodifier.DefaultConfig()
		encodeModifier, err = encodemodifier.New(modifierConfig)
		if err != nil {
			// error handling
		}
	}

	var encryptModifier valuemodifier.ValueModifier
	{
		modifierConfig := encryptmodifier.DefaultConfig()
		modifierConfig.Pass = "somesecretpassphrase"
		encryptModifier, err = encryptmodifier.New(modifierConfig)
		if err != nil {
			// error handling
		}
	}

	var newTraverser *valuemodifier.Service
	{
		config := valuemodifier.DefaultConfig()
		config.ValueModifiers = []valuemodifier.ValueModifier{encryptModifier, encodeModifier}
		newTraverser, err = valuemodifier.New(config)
		if err != nil {
			// error handling
		}
	}

	decrypted, err := newTraverser.TraverseJSON([]byte(`{"key":"val"}`))
	if err != nil {
		// error handling
	}

	// Variable decrypted holds the string "val" in GPG encrypted and base64 encoded
	// form.
	//
	//     {
	//       "key": "LS0tLS1CRUdJTiBQR1AgU0lHTkFUVVJFLS0tLS0KCnd4NEVCd01JdWk3S1BHcmRWK0JnU0U1RFJMVVY1L3l0UXp2dWwwUHhXS2JTNEFIa3Znc2ZreFJSZFZuUzdLVG4KTVkvZnJ1RXlCZURwNEp2aC90cmdBZUp6dDZGWjRQTGhZQ2JneXVCdTRKVGsrcWVzNm96dCtvUVJCRDcvTXNpSQpBZUxEeWk5bjRWYUdBQT09Cj14YTQ5Ci0tLS0tRU5EIFBHUCBTSUdOQVRVUkUtLS0tLQ=="
	//     }
	//
	fmt.Printf("%s\n", decrypted)
}
```
