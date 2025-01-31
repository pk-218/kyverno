package common

import (
	"encoding/json"
	"testing"

	"github.com/kyverno/kyverno/pkg/engine/operator"
	"github.com/kyverno/kyverno/pkg/logging"
	"gotest.tools/assert"
)

func TestValidateValueWithPattern_Bool(t *testing.T) {
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), true, true))
	assert.Assert(t, !ValidateValueWithPattern(logging.GlobalLogger(), true, false))
	assert.Assert(t, !ValidateValueWithPattern(logging.GlobalLogger(), false, true))
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), false, false))
}

func TestValidateString_AsteriskTest(t *testing.T) {
	pattern := "*"
	value := "anything"
	empty := ""

	assert.Assert(t, validateString(logging.GlobalLogger(), value, pattern, operator.Equal))
	assert.Assert(t, validateString(logging.GlobalLogger(), empty, pattern, operator.Equal))
}

func TestValidateString_LeftAsteriskTest(t *testing.T) {
	pattern := "*right"
	value := "leftright"
	right := "right"

	assert.Assert(t, validateString(logging.GlobalLogger(), value, pattern, operator.Equal))
	assert.Assert(t, validateString(logging.GlobalLogger(), right, pattern, operator.Equal))

	value = "leftmiddle"
	middle := "middle"

	assert.Assert(t, !validateString(logging.GlobalLogger(), value, pattern, operator.Equal))
	assert.Assert(t, !validateString(logging.GlobalLogger(), middle, pattern, operator.Equal))
}

func TestValidateString_MiddleAsteriskTest(t *testing.T) {
	pattern := "ab*ba"
	value := "abbeba"
	assert.Assert(t, validateString(logging.GlobalLogger(), value, pattern, operator.Equal))

	value = "abbca"
	assert.Assert(t, !validateString(logging.GlobalLogger(), value, pattern, operator.Equal))
}

func TestValidateString_QuestionMark(t *testing.T) {
	pattern := "ab?ba"
	value := "abbba"
	assert.Assert(t, validateString(logging.GlobalLogger(), value, pattern, operator.Equal))

	value = "abbbba"
	assert.Assert(t, !validateString(logging.GlobalLogger(), value, pattern, operator.Equal))
}

func TestValidateValueWithPattern_BoolInJson(t *testing.T) {
	rawPattern := []byte(`
	{
		"key": true
	}
	`)

	rawValue := []byte(`
	{
		"key": true
	}
	`)

	var pattern, value map[string]interface{}
	err := json.Unmarshal(rawPattern, &pattern)
	assert.Assert(t, err)
	err = json.Unmarshal(rawValue, &value)
	assert.Assert(t, err)

	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value["key"], pattern["key"]))
}

func TestValidateValueWithPattern_NullPatternStringValue(t *testing.T) {
	rawPattern := []byte(`
	{
		"key": null
	}
	`)

	rawValue := []byte(`
	{
		"key": "value"
	}
	`)

	var pattern, value map[string]interface{}
	err := json.Unmarshal(rawPattern, &pattern)
	assert.Assert(t, err)
	err = json.Unmarshal(rawValue, &value)
	assert.Assert(t, err)

	assert.Assert(t, !ValidateValueWithPattern(logging.GlobalLogger(), value["key"], pattern["key"]))
}

func TestValidateValueWithPattern_NullPatternDefaultString(t *testing.T) {
	rawPattern := []byte(`
	{
		"key": null
	}
	`)

	rawValue := []byte(`
	{
		"key": ""
	}
	`)

	var pattern, value map[string]interface{}
	err := json.Unmarshal(rawPattern, &pattern)
	assert.Assert(t, err)
	err = json.Unmarshal(rawValue, &value)
	assert.Assert(t, err)

	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value["key"], pattern["key"]))
}

func TestValidateValueWithPattern_NullPatternDefaultFloat(t *testing.T) {
	rawPattern := []byte(`
	{
		"key": null
	}
	`)

	rawValue := []byte(`
	{
		"key": 0.0
	}
	`)

	var pattern, value map[string]interface{}
	err := json.Unmarshal(rawPattern, &pattern)
	assert.Assert(t, err)
	err = json.Unmarshal(rawValue, &value)
	assert.Assert(t, err)

	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value["key"], pattern["key"]))
}

func TestValidateValueWithPattern_NullPatternDefaultInt(t *testing.T) {
	rawPattern := []byte(`
	{
		"key": null
	}
	`)

	rawValue := []byte(`
	{
		"key": 0
	}
	`)

	var pattern, value map[string]interface{}
	err := json.Unmarshal(rawPattern, &pattern)
	assert.Assert(t, err)
	err = json.Unmarshal(rawValue, &value)
	assert.Assert(t, err)

	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value["key"], pattern["key"]))
}

func TestValidateValueWithPattern_NullPatternDefaultBool(t *testing.T) {
	rawPattern := []byte(`
	{
		"key": null
	}
	`)

	rawValue := []byte(`
	{
		"key": false
	}
	`)

	var pattern, value map[string]interface{}
	err := json.Unmarshal(rawPattern, &pattern)
	assert.Assert(t, err)
	err = json.Unmarshal(rawValue, &value)
	assert.Assert(t, err)

	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value["key"], pattern["key"]))
}

func TestValidateValueWithPattern_StringsLogicalOr(t *testing.T) {
	pattern := "192.168.88.1 | 10.100.11.*"
	value := "10.100.11.54"
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value, pattern))
}

func TestValidateValueWithPattern_StringsLogicalAnd(t *testing.T) {
	pattern := ">1 & <20"
	value := "10"
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value, pattern))
}

func TestValidateValueWithPattern_StringsAllLogicalOperators(t *testing.T) {
	pattern := ">1 & <20 | >31 & <33"
	value := "10"
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value, pattern))
	value = "32"
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), value, pattern))
	value = "21"
	assert.Assert(t, !ValidateValueWithPattern(logging.GlobalLogger(), value, pattern))
}

func TestValidateValueWithPattern_EqualTwoFloats(t *testing.T) {
	assert.Assert(t, ValidateValueWithPattern(logging.GlobalLogger(), 7.0, 7.000))
}

func TestValidateValueWithNilPattern_NullPatternStringValue(t *testing.T) {
	assert.Assert(t, !validateValueWithNilPattern(logging.GlobalLogger(), "value"))
}

func TestValidateValueWithNilPattern_NullPatternDefaultString(t *testing.T) {
	assert.Assert(t, validateValueWithNilPattern(logging.GlobalLogger(), ""))
}

func TestValidateValueWithNilPattern_NullPatternDefaultFloat(t *testing.T) {
	assert.Assert(t, validateValueWithNilPattern(logging.GlobalLogger(), 0.0))
}

func TestValidateValueWithNilPattern_NullPatternFloat(t *testing.T) {
	assert.Assert(t, !validateValueWithNilPattern(logging.GlobalLogger(), 0.1))
}

func TestValidateValueWithNilPattern_NullPatternDefaultInt(t *testing.T) {
	assert.Assert(t, validateValueWithNilPattern(logging.GlobalLogger(), 0))
}

func TestValidateValueWithNilPattern_NullPatternInt(t *testing.T) {
	assert.Assert(t, !validateValueWithNilPattern(logging.GlobalLogger(), 1))
}

func TestValidateValueWithNilPattern_NullPatternDefaultBool(t *testing.T) {
	assert.Assert(t, validateValueWithNilPattern(logging.GlobalLogger(), false))
}

func TestValidateValueWithNilPattern_NullPatternTrueBool(t *testing.T) {
	assert.Assert(t, !validateValueWithNilPattern(logging.GlobalLogger(), true))
}

func TestValidateValueWithFloatPattern_FloatValue(t *testing.T) {
	assert.Assert(t, validateValueWithFloatPattern(logging.GlobalLogger(), 7.9914, 7.9914))
}

func TestValidateValueWithFloatPattern_FloatValueNotPass(t *testing.T) {
	assert.Assert(t, !validateValueWithFloatPattern(logging.GlobalLogger(), 7.9914, 7.99141))
}

func TestValidateValueWithFloatPattern_FloatPatternWithoutFractionIntValue(t *testing.T) {
	assert.Assert(t, validateValueWithFloatPattern(logging.GlobalLogger(), 7, 7.000000))
}

func TestValidateValueWithFloatPattern_FloatPatternWithoutFraction(t *testing.T) {
	assert.Assert(t, validateValueWithFloatPattern(logging.GlobalLogger(), 7.000000, 7.000000))
}

func TestValidateValueWithIntPattern_FloatValueWithoutFraction(t *testing.T) {
	assert.Assert(t, validateValueWithFloatPattern(logging.GlobalLogger(), 7.000000, 7))
}

func TestValidateValueWithIntPattern_FloatValueWitFraction(t *testing.T) {
	assert.Assert(t, !validateValueWithFloatPattern(logging.GlobalLogger(), 7.000001, 7))
}

func TestValidateValueWithIntPattern_NotPass(t *testing.T) {
	assert.Assert(t, !validateValueWithFloatPattern(logging.GlobalLogger(), 8, 7))
}

func TestGetNumberAndStringPartsFromPattern_NumberAndString(t *testing.T) {
	number, str := getNumberAndStringPartsFromPattern("1024Gi")
	assert.Equal(t, number, "1024")
	assert.Equal(t, str, "Gi")
}

func TestGetNumberAndStringPartsFromPattern_OnlyNumber(t *testing.T) {
	number, str := getNumberAndStringPartsFromPattern("1024")
	assert.Equal(t, number, "1024")
	assert.Equal(t, str, "")
}

func TestGetNumberAndStringPartsFromPattern_OnlyString(t *testing.T) {
	number, str := getNumberAndStringPartsFromPattern("Gi")
	assert.Equal(t, number, "")
	assert.Equal(t, str, "Gi")
}

func TestGetNumberAndStringPartsFromPattern_StringFirst(t *testing.T) {
	number, str := getNumberAndStringPartsFromPattern("Gi1024")
	assert.Equal(t, number, "")
	assert.Equal(t, str, "Gi1024")
}

func TestGetNumberAndStringPartsFromPattern_Empty(t *testing.T) {
	number, str := getNumberAndStringPartsFromPattern("")
	assert.Equal(t, number, "")
	assert.Equal(t, str, "")
}

func TestValidateValueWithStringPattern_WithSpace(t *testing.T) {
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 4, ">= 3"))
}

func TestValidateValueWithStringPattern_Ranges(t *testing.T) {
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 0, "0-2"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 1, "0-2"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 2, "0-2"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 3, "0-2"))

	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 0, "10!-20"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 15, "10!-20"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 25, "10!-20"))

	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 0, "0.00001-2.00001"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 1, "0.00001-2.00001"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 2, "0.00001-2.00001"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 2.0001, "0.00001-2.00001"))

	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 0, "0.00001!-2.00001"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 1, "0.00001!-2.00001"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 2, "0.00001!-2.00001"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 2.0001, "0.00001!-2.00001"))

	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 2, "2-2"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), 2, "2!-2"))

	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 2.99999, "2.99998-3"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 2.99997, "2.99998!-3"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), 3.00001, "2.99998!-3"))

	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), "256Mi", "128Mi-512Mi"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), "1024Mi", "128Mi-512Mi"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), "64Mi", "128Mi-512Mi"))

	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), "256Mi", "128Mi!-512Mi"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), "1024Mi", "128Mi!-512Mi"))
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), "64Mi", "128Mi!-512Mi"))
}

func TestValidateNumberWithStr_LessFloatAndInt(t *testing.T) {
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), 7.00001, "7.000001", operator.More))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), 7.00001, "7", operator.NotEqual))

	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), 7.0000, "7", operator.Equal))
	assert.Assert(t, !validateNumberWithStr(logging.GlobalLogger(), 6.000000001, "6", operator.Less))
}

func TestValidateQuantity_InvalidQuantity(t *testing.T) {
	assert.Assert(t, !validateNumberWithStr(logging.GlobalLogger(), "1024Gi", "", operator.Equal))
	assert.Assert(t, !validateNumberWithStr(logging.GlobalLogger(), "gii", "1024Gi", operator.Equal))
}

func TestValidateQuantity_Equal(t *testing.T) {
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "1024Gi", "1024Gi", operator.Equal))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "1024Mi", "1Gi", operator.Equal))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "0.2", "200m", operator.Equal))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "500", "500", operator.Equal))
	assert.Assert(t, !validateNumberWithStr(logging.GlobalLogger(), "2048", "1024", operator.Equal))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), 1024, "1024", operator.Equal))
}

func TestValidateQuantity_Operation(t *testing.T) {
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "1Gi", "1000Mi", operator.More))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "1G", "1Gi", operator.Less))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "500m", "0.5", operator.MoreEqual))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "1", "500m", operator.MoreEqual))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "0.5", ".5", operator.LessEqual))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "0.2", ".5", operator.LessEqual))
	assert.Assert(t, validateNumberWithStr(logging.GlobalLogger(), "0.2", ".5", operator.NotEqual))
}

func TestGetOperatorFromStringPattern_OneChar(t *testing.T) {
	assert.Equal(t, operator.GetOperatorFromStringPattern("f"), operator.Equal)
}

func TestGetOperatorFromStringPattern_EmptyString(t *testing.T) {
	assert.Equal(t, operator.GetOperatorFromStringPattern(""), operator.Equal)
}

func TestValidateKernelVersion_NotEquals(t *testing.T) {
	assert.Assert(t, validateValueWithStringPattern(logging.GlobalLogger(), "5.16.5-arch1-1", "!5.10.84-1"))
	assert.Assert(t, !validateValueWithStringPattern(logging.GlobalLogger(), "5.10.84-1", "!5.10.84-1"))
	assert.Assert(t, validateValueWithStringPatterns(logging.GlobalLogger(), "5.16.5-arch1-1", "!5.10.84-1 & !5.15.2-1"))
	assert.Assert(t, !validateValueWithStringPatterns(logging.GlobalLogger(), "5.10.84-1", "!5.10.84-1 & !5.15.2-1"))
	assert.Assert(t, !validateValueWithStringPatterns(logging.GlobalLogger(), "5.15.2-1", "!5.10.84-1 & !5.15.2-1"))
}
