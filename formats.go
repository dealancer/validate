package validate

import (
	"bytes"
	"crypto/sha256"
	"net"
	"net/url"
	"os"
	"strings"

	urn "github.com/leodido/go-urn"
)

// FormatType is used for format validator type definitions.
type FormatType string

// Following string formats are available.
// E.g. `validate:"format=email"`
const (
	FormatAlpha                FormatType = "alpha"
	FormatAlnum                FormatType = "alnum"
	FormatAlphaUnicode         FormatType = "alpha_unicode"
	FormatAlnumUnicode         FormatType = "alnum_unicode"
	FormatNumeric              FormatType = "numeric"
	FormatNumber               FormatType = "number"
	FormatHexadecimal          FormatType = "hexadecimal"
	FormatHEXColor             FormatType = "hexcolor"
	FormatRGB                  FormatType = "rgb"
	FormatRGBA                 FormatType = "rgba"
	FormatHSL                  FormatType = "hsl"
	FormatHSLA                 FormatType = "hsla"
	FormatEmail                FormatType = "email"
	FormatURL                  FormatType = "url"
	FormatURI                  FormatType = "uri"
	FormatUrnRFC2141           FormatType = "urn_rfc2141" // RFC 2141
	FormatFile                 FormatType = "file"
	FormatBase64               FormatType = "base64"
	FormatBase64URL            FormatType = "base64url"
	FormatISBN                 FormatType = "isbn"
	FormatISBN10               FormatType = "isbn10"
	FormatISBN13               FormatType = "isbn13"
	FormatEthereumAddress      FormatType = "eth_addr"
	FormatBitcoinAddress       FormatType = "btc_addr"
	FormatBitcoinBech32Address FormatType = "btc_addr_bech32"
	FormatUUID                 FormatType = "uuid"
	FormatUUID3                FormatType = "uuid3"
	FormatUUID4                FormatType = "uuid4"
	FormatUUID5                FormatType = "uuid5"
	FormatASCII                FormatType = "ascii"
	FormatPrintableASCII       FormatType = "ascii_print"
	FormatDataURI              FormatType = "datauri"
	FormatLatitude             FormatType = "latitude"
	FormatLongitude            FormatType = "longitude"
	FormatSSN                  FormatType = "ssn"
	FormatIPv4                 FormatType = "ipv4"
	FormatIPv6                 FormatType = "ipv6"
	FormatIP                   FormatType = "ip"
	FormatCIDRv4               FormatType = "cidrv4"
	FormatCIDRv6               FormatType = "cidrv6"
	FormatCIDR                 FormatType = "cidr"
	FormatMAC                  FormatType = "mac"
	FormatHostnameRFC952       FormatType = "hostname"         // RFC 952
	FormatHostnameRFC1123      FormatType = "hostname_rfc1123" // RFC 1123
	FormatFQDN                 FormatType = "fqdn"
	FormatURLEncoded           FormatType = "url_encoded"
	FormatDir                  FormatType = "dir"
	FormatPostcode             FormatType = "postcode"
)

// formatFunc is an interface for format validator func
type formatFunc func(value string) bool

func getFormatTypeMap() map[FormatType]formatFunc {
	return map[FormatType]formatFunc{
		FormatAlpha:                formatAlpha,
		FormatAlnum:                formatAlnum,
		FormatAlphaUnicode:         formatAlphaUnicode,
		FormatAlnumUnicode:         formatAlnumUnicode,
		FormatNumeric:              formatNumeric,
		FormatNumber:               formatNumber,
		FormatHexadecimal:          formatHexadecimal,
		FormatHEXColor:             formatHEXColor,
		FormatRGB:                  formatRGB,
		FormatRGBA:                 formatRGBA,
		FormatHSL:                  formatHSL,
		FormatHSLA:                 formatHSLA,
		FormatEmail:                formatEmail,
		FormatURL:                  formatURL,
		FormatURI:                  formatURI,
		FormatUrnRFC2141:           formatUrnRFC2141,
		FormatFile:                 formatFile,
		FormatBase64:               formatBase64,
		FormatBase64URL:            formatBase64URL,
		FormatISBN:                 formatISBN,
		FormatISBN10:               formatISBN10,
		FormatISBN13:               formatISBN13,
		FormatEthereumAddress:      formatEthereumAddress,
		FormatBitcoinAddress:       formatBitcoinAddress,
		FormatBitcoinBech32Address: formatBitcoinBech32Address,
		FormatUUID:                 formatUUID,
		FormatUUID3:                formatUUID3,
		FormatUUID4:                formatUUID4,
		FormatUUID5:                formatUUID5,
		FormatASCII:                formatASCII,
		FormatPrintableASCII:       formatPrintableASCII,
		FormatDataURI:              formatDataURI,
		FormatLatitude:             formatLatitude,
		FormatLongitude:            formatLongitude,
		FormatSSN:                  formatSSN,
		FormatIPv4:                 formatIPv4,
		FormatIPv6:                 formatIPv6,
		FormatIP:                   formatIP,
		FormatCIDRv4:               formatCIDRv4,
		FormatCIDRv6:               formatCIDRv6,
		FormatCIDR:                 formatCIDR,
		FormatMAC:                  formatMAC,
		FormatHostnameRFC952:       formatHostnameRFC952,
		FormatHostnameRFC1123:      formatHostnameRFC1123,
		FormatFQDN:                 formatFQDN,
		FormatURLEncoded:           formatURLEncoded,
		FormatDir:                  formatDir,
		FormatPostcode:             formatPostcode,
	}
}

func formatURLEncoded(value string) bool {
	return uRLEncodedRegex.MatchString(value)
}

// formatMAC is the validation function for validating if the field's value is a valid MAC address.
func formatMAC(value string) bool {
	_, err := net.ParseMAC(value)

	return err == nil
}

// formatCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func formatCIDRv4(value string) bool {
	ip, _, err := net.ParseCIDR(value)

	return err == nil && ip.To4() != nil
}

// formatCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func formatCIDRv6(value string) bool {

	ip, _, err := net.ParseCIDR(value)

	return err == nil && ip.To4() == nil
}

// formatCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func formatCIDR(value string) bool {
	_, _, err := net.ParseCIDR(value)

	return err == nil
}

// formatIPv4 is the validation function for validating if a value is a valid v4 IP address.
func formatIPv4(value string) bool {
	ip := net.ParseIP(value)

	return ip != nil && ip.To4() != nil
}

// formatIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func formatIPv6(value string) bool {
	ip := net.ParseIP(value)

	return ip != nil && ip.To4() == nil
}

// formatIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func formatIP(value string) bool {
	ip := net.ParseIP(value)

	return ip != nil
}

// formatSSN is the validation function for validating if the field's value is a valid SSN.
func formatSSN(value string) bool {
	if len(value) != 11 {
		return false
	}

	return sSNRegex.MatchString(value)
}

// formatLongitude is the validation function for validating if the field's value is a valid longitude coordinate.
func formatLongitude(value string) bool {
	return longitudeRegex.MatchString(value)
}

// formatLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
func formatLatitude(value string) bool {
	return latitudeRegex.MatchString(value)
}

// formatDataURI is the validation function for validating if the field's value is a valid data URI.
func formatDataURI(value string) bool {
	uri := strings.SplitN(value, ",", 2)

	if len(uri) != 2 {
		return false
	}

	if !dataURIRegex.MatchString(uri[0]) {
		return false
	}

	return base64Regex.MatchString(uri[1])
}

// formatPrintableASCII is the validation function for validating if the field's value is a valid printable ASCII character.
func formatPrintableASCII(value string) bool {
	return printableASCIIRegex.MatchString(value)
}

// formatASCII is the validation function for validating if the field's value is a valid ASCII character.
func formatASCII(value string) bool {
	return aSCIIRegex.MatchString(value)
}

// formatUUID5 is the validation function for validating if the field's value is a valid v5 UUID.
func formatUUID5(value string) bool {
	return uUID5Regex.MatchString(value)
}

// formatUUID4 is the validation function for validating if the field's value is a valid v4 UUID.
func formatUUID4(value string) bool {
	return uUID4Regex.MatchString(value)
}

// formatUUID3 is the validation function for validating if the field's value is a valid v3 UUID.
func formatUUID3(value string) bool {
	return uUID3Regex.MatchString(value)
}

// formatUUID is the validation function for validating if the field's value is a valid UUID of any version.
func formatUUID(value string) bool {
	return uUIDRegex.MatchString(value)
}

// formatISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
func formatISBN(value string) bool {
	return formatISBN10(value) || formatISBN13(value)
}

// formatISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
func formatISBN13(value string) bool {
	s := strings.Replace(strings.Replace(value, "-", "", 4), " ", "", 4)

	if !iSBN13Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	factor := []int32{1, 3}

	for i = 0; i < 12; i++ {
		checksum += factor[i%2] * int32(s[i]-'0')
	}

	return (int32(s[12]-'0'))-((10-(checksum%10))%10) == 0
}

// formatISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
func formatISBN10(value string) bool {
	s := strings.Replace(strings.Replace(value, "-", "", 3), " ", "", 3)

	if !iSBN10Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	for i = 0; i < 9; i++ {
		checksum += (i + 1) * int32(s[i]-'0')
	}

	if s[9] == 'X' {
		checksum += 10 * 10
	} else {
		checksum += 10 * int32(s[9]-'0')
	}

	return checksum%11 == 0
}

// formatEthereumAddress is the validation function for validating if the field's value is a valid ethereum address based currently only on the format
func formatEthereumAddress(value string) bool {
	address := value

	if !ethAddressRegex.MatchString(address) {
		return false
	}

	if ethAddressRegexUpper.MatchString(address) || ethAddressRegexLower.MatchString(address) {
		return true
	}

	// checksum validation is blocked by https://github.com/golang/crypto/pull/28

	return true
}

// formatBitcoinAddress is the validation function for validating if the field's value is a valid btc address
func formatBitcoinAddress(value string) bool {
	address := value

	if !btcAddressRegex.MatchString(address) {
		return false
	}

	alphabet := []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	decode := [25]byte{}

	for _, n := range []byte(address) {
		d := bytes.IndexByte(alphabet, n)

		for i := 24; i >= 0; i-- {
			d += 58 * int(decode[i])
			decode[i] = byte(d % 256)
			d /= 256
		}
	}

	h := sha256.New()
	_, _ = h.Write(decode[:21])
	d := h.Sum([]byte{})
	h = sha256.New()
	_, _ = h.Write(d)

	validchecksum := [4]byte{}
	computedchecksum := [4]byte{}

	copy(computedchecksum[:], h.Sum(d[:0]))
	copy(validchecksum[:], decode[21:])

	return validchecksum == computedchecksum
}

// formatBitcoinBech32Address is the validation function for validating if the field's value is a valid bech32 btc address
func formatBitcoinBech32Address(value string) bool {
	address := value

	if !btcLowerAddressRegexBech32.MatchString(address) && !btcUpperAddressRegexBech32.MatchString(address) {
		return false
	}

	am := len(address) % 8

	if am == 0 || am == 3 || am == 5 {
		return false
	}

	address = strings.ToLower(address)

	alphabet := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

	hr := []int{3, 3, 0, 2, 3} // the human readable part will always be bc
	addr := address[3:]
	dp := make([]int, 0, len(addr))

	for _, c := range addr {
		dp = append(dp, strings.IndexRune(alphabet, c))
	}

	ver := dp[0]

	if ver < 0 || ver > 16 {
		return false
	}

	if ver == 0 {
		if len(address) != 42 && len(address) != 62 {
			return false
		}
	}

	values := append(hr, dp...)

	GEN := []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

	p := 1

	for _, v := range values {
		b := p >> 25
		p = (p&0x1ffffff)<<5 ^ v

		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				p ^= GEN[i]
			}
		}
	}

	if p != 1 {
		return false
	}

	b := uint(0)
	acc := 0
	mv := (1 << 5) - 1
	var sw []int

	for _, v := range dp[1 : len(dp)-6] {
		acc = (acc << 5) | v
		b += 5
		for b >= 8 {
			b -= 8
			sw = append(sw, (acc>>b)&mv)
		}
	}

	if len(sw) < 2 || len(sw) > 40 {
		return false
	}

	return true
}

// formatBase64 is the validation function for validating if the current field's value is a valid base 64.
func formatBase64(value string) bool {
	return base64Regex.MatchString(value)
}

// formatBase64URL is the validation function for validating if the current field's value is a valid base64 URL safe string.
func formatBase64URL(value string) bool {
	return base64URLRegex.MatchString(value)
}

// formatURI is the validation function for validating if the current field's value is a valid URI.
func formatURI(value string) bool {
	// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
	// emulate browser and strip the '#' suffix prior to validation. see issue-#237
	if i := strings.Index(value, "#"); i > -1 {
		value = value[:i]
	}

	if len(value) == 0 {
		return false
	}

	_, err := url.ParseRequestURI(value)

	return err == nil
}

// formatURL is the validation function for validating if the current field's value is a valid URL.
func formatURL(value string) bool {
	var i int

	// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
	// emulate browser and strip the '#' suffix prior to validation. see issue-#237
	if i = strings.Index(value, "#"); i > -1 {
		value = value[:i]
	}

	if len(value) == 0 {
		return false
	}

	url, err := url.ParseRequestURI(value)

	if err != nil || url.Scheme == "" {
		return false
	}

	return err == nil
}

// formatUrnRFC2141 is the validation function for validating if the current field's value is a valid URN as per RFC 2141.
func formatUrnRFC2141(value string) bool {
	_, match := urn.Parse([]byte(value))

	return match

}

// formatFile is the validation function for validating if the current field's value is a valid file path.
func formatFile(value string) bool {
	fileInfo, err := os.Stat(value)
	if err != nil {
		return false
	}

	return !fileInfo.IsDir()
}

// formatEmail is the validation function for validating if the current field's value is a valid email address.
func formatEmail(value string) bool {
	return emailRegex.MatchString(value)
}

// formatHSLA is the validation function for validating if the current field's value is a valid HSLA color.
func formatHSLA(value string) bool {
	return hslaRegex.MatchString(value)
}

// formatHSL is the validation function for validating if the current field's value is a valid HSL color.
func formatHSL(value string) bool {
	return hslRegex.MatchString(value)
}

// formatRGBA is the validation function for validating if the current field's value is a valid RGBA color.
func formatRGBA(value string) bool {
	return rgbaRegex.MatchString(value)
}

// formatRGB is the validation function for validating if the current field's value is a valid RGB color.
func formatRGB(value string) bool {
	return rgbRegex.MatchString(value)
}

// formatHEXColor is the validation function for validating if the current field's value is a valid HEX color.
func formatHEXColor(value string) bool {
	return hexcolorRegex.MatchString(value)
}

// formatHexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
func formatHexadecimal(value string) bool {
	return hexadecimalRegex.MatchString(value)
}

// formatNumber is the validation function for validating if the current field's value is a valid number.
func formatNumber(value string) bool {
	return numberRegex.MatchString(value)
}

// formatNumeric is the validation function for validating if the current field's value is a valid numeric value.
func formatNumeric(value string) bool {
	return numericRegex.MatchString(value)
}

// formatAlnum is the validation function for validating if the current field's value is a valid alphanumeric value.
func formatAlnum(value string) bool {
	return alnumRegex.MatchString(value)
}

// formatAlpha is the validation function for validating if the current field's value is a valid alpha value.
func formatAlpha(value string) bool {
	return alphaRegex.MatchString(value)
}

// formatAlnumUnicode is the validation function for validating if the current field's value is a valid alphanumeric unicode value.
func formatAlnumUnicode(value string) bool {
	return alnumUnicodeRegex.MatchString(value)
}

// formatAlphaUnicode is the validation function for validating if the current field's value is a valid alpha unicode value.
func formatAlphaUnicode(value string) bool {
	return alphaUnicodeRegex.MatchString(value)
}

func formatHostnameRFC952(value string) bool {
	return hostnameRegexRFC952.MatchString(value)
}

func formatHostnameRFC1123(value string) bool {
	return hostnameRegexRFC1123.MatchString(value)
}

func formatFQDN(value string) bool {
	val := value

	if val == "" {
		return false
	}

	if val[len(val)-1] == '.' {
		val = val[0 : len(val)-1]
	}

	return strings.ContainsAny(val, ".") &&
		hostnameRegexRFC952.MatchString(val)
}

// formatDir is the validation function for validating if the current field's value is a valid directory.
func formatDir(value string) bool {
	fileInfo, err := os.Stat(value)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// formatPostcode is the validation function for validating if the current field's value is a valid postcode.
func formatPostcode(value string) bool {
	return postcodeRegex.MatchString(value)
}
