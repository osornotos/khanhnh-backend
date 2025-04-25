package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// @khanh.nguyen_dev should fix this failed tests
func testArgon2idHash_GenerateAndCompare(t *testing.T) {
	hash := NewArgon2idHash(1, 32, 64, 4, 128)
	salt := "random"

	type input struct {
		raw, hash string
		salt      string
	}
	tb := []input{
		{salt: salt, raw: `12345678`, hash: `$a$v19$m64+t1+p4$cmFuZG9t$PkJhoQrWz2uMtu3pmO46/iA+TE5bMP58u8uUYCXLoH9m0hUuaZNfED2EAT5BWummIBB2DtCRtDUtoUaYxdz4P8TtISA+r9gori+/E8ffTLVBhO0nQOgHrtt8Z1WpZ9NPMaqo4mls4pGADOl8vz1ZtiXsLyWYdCe4MGk24L7F/sQ`},
		{salt: salt, raw: `password`, hash: `$a$v19$m64+t1+p4$cmFuZG9t$aPPhyvveujNJQQV4OQhvDhSE/a6fEp9Qhsv2a6PVLWdumd6s4yIiPl5pOXj/UJl69v4NDjJgrVv/UHgbdQHrclHx55oQuzYAo1q1sNPrMETNvi7YRQ3UvZDnlkId3w7baZO/Cwfhy6J/8bGY3jziNaIaLJXhJNn7ccai4XLtabw`},
		{salt: salt, raw: `sup3r s3cr3t`, hash: `$a$v19$m64+t1+p4$cmFuZG9t$BGUd3iZ46I0ewD88CVuSIoKieaPaq/PyosUnybUB6r4JuXY6tL5bcs762/6VIv0m0Uz+TIt5y8HE4H4rFskqQptjQEbpGe8ECDZIC5o1SHrhs7HuFKVGjWmNu6N/GwI5UHnkkm8YtQvWMbhuJClIAO6WAs8w3ps6zU785+zOSWk`},
		{salt: salt, raw: `0`, hash: `$a$v19$m64+t1+p4$cmFuZG9t$3dAMZ/CbyaYx01Bn61A0DgAgcY1yVHdNLLvpm1VFAMqXa9b/MNJ2hVTBGDR2xr/qogKdnQop21zu6FH5U7SjvHDAD4O/VFdB5okfTOELkoeEugBEyruciRHHZZddjX157LrsecdAhIkFumoC9RZ/lqte7st6ti7C+pKXZrQZa/k`},
		{salt: salt, raw: `qwe123!@#`, hash: `$a$v19$m64+t1+p4$cmFuZG9t$rLZrO7lDAr7g9+kAZg4i+15Q1cX+hOsu9KnMWjgWQOq5rZ+n8g4PRBVbRtTjlGaFJB/igOG1yrJ2uC2Ehjl+00zlzBjAkGwJQ0XDiKyiPdxWGsHU2Gu72p4W74G46goPX6yG2TCz0iTIJJgSTsxy4OU11e/6OoU9PWZ30RYbxGg`},
		{salt: salt, raw: `tiếng việt`, hash: `$a$v19$m64+t1+p4$cmFuZG9t$7q1wJTm+Nu+0nBi2Uv+uG0Y9zlLL91isd1d9cbmBknlXbqj27OVsGabbJSx/CbIRX79dhjfh88SGw1z1AWXHEefQ31jCjbyqW7JG2Dn+34KZ22Ezmf+GVOvdVkBcsn+PFbuk+TlS06G1V4o+Sa/+IGPYjPtTxAk2tdoWzw2HwMQ`},
		{salt: salt, raw: ``, hash: `$a$v19$m64+t1+p4$cmFuZG9t$8FeuhoiEgm0Z9bOD3HvRnyyLsMyiJpKRtUooNp7THa6Aei6tgYwIOWceYSnYURusSNvNVDJLcS2Zhvi13OosY17JdlogbnrauWJXwW8rIgEfPKKGdfUoJzzgU2c7eyFeMhHgGmVQOQYibQPmnE9V40wCb/dq47uSoTq6JQLre3U`},
	}
	for _, tc := range tb {
		t.Run(tc.raw, func(t *testing.T) {
			result, err := hash.GenerateWithSalt(tc.raw, []byte(tc.salt))
			assert.Nil(t, err)
			if result == nil {
				t.Fail()
				return
			}

			assert.Equal(t, tc.hash, result.HashEncoded)
		})
	}
	for _, tc := range tb {
		t.Run(tc.raw, func(t *testing.T) {
			err := hash.Compare(tc.raw, tc.hash)
			assert.Nil(t, err)
		})
	}
}
