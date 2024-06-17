package utils_test

import (
	"pixie79/utils"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = Describe("Masking", func() {
	// Define common variables for use in test cases
	var (
		s             string
		maskChar      string
		option        string
		result        string
		unmaskedCount int
	)

	BeforeEach(func() {
		// Setup initial values
		s = "HelloWorld"
		maskChar = "*"
	})

	Context("when option is 'first'", func() {
		It("should return the string with first part unmasked", func() {
			option = "first"
			unmaskedCount = 3
			expected := "Hel*******" // trunk-ignore(codespell/misspelled)
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})

		It("should return full unmasked if unmaskedCount is greater than string length", func() {
			option = "first"
			unmaskedCount = 15
			expected := "HelloWorld"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})

		It("should return all masked if unmaskedCount is zero", func() {
			option = "first"
			unmaskedCount = 0
			expected := "**********"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})
	})

	Context("when option is 'last'", func() {
		It("should return the string with last part unmasked", func() {
			option = "last"
			unmaskedCount = 4
			expected := "******orld"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})

		It("should return all unmasked if unmaskedCount is greater than string length", func() {
			option = "last"
			unmaskedCount = 12
			expected := "HelloWorld"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})

		It("should return all masked if unmaskedCount is zero", func() {
			option = "last"
			unmaskedCount = 0
			expected := "**********"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})
	})

	Context("when option is 'fixed'", func() {
		It("should return a fixed number of masked characters to the length of unmaskedCount (0)", func() {
			option = "fixed"
			unmaskedCount = 0
			expected := ""
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})

		It("should return a fixed number of masked characters to the length of unmaskedCount (2)", func() {
			option = "fixed"
			unmaskedCount = 2
			expected := "**"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})

		It("should return a fixed number of masked characters to the length of unmaskedCount (8)", func() {
			option = "fixed"
			unmaskedCount = 8
			expected := "********"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})
	})

	Context("when option is invalid", func() {
		It("should return an error message", func() {
			option = "unknown"
			unmaskedCount = 5
			expected := "XXERRORXX"
			result = utils.MaskString(s, maskChar, option, unmaskedCount)
			gomega.Expect(result).To(gomega.Equal(expected))
		})
	})
})

var _ = Describe("utils.MaskStringInMap", func() {
	// Define common variables for use in test cases
	var (
		inputMap      map[string]string
		maskChar      string
		option        string
		unmaskedCount int
		expectedMap   map[string]string
		resultMap     map[string]string
	)

	BeforeEach(func() {
		// Setup initial values
		maskChar = "*"
		inputMap = map[string]string{
			"name":    "Alice",
			"city":    "Seattle",
			"country": "USA",
		}
	})

	Context("when option is 'first'", func() {
		It("should mask all values with the first part unmasked", func() {
			option = "first"
			unmaskedCount = 2
			expectedMap = map[string]string{
				"name":    "Al***",
				"city":    "Se*****",
				"country": "US*",
			}
			resultMap = utils.MaskStringInMap(inputMap, maskChar, option, unmaskedCount)
			gomega.Expect(resultMap).To(gomega.Equal(expectedMap))
		})
	})

	Context("when option is 'last'", func() {
		It("should mask all values with the last part unmasked", func() {
			option = "last"
			unmaskedCount = 3
			expectedMap = map[string]string{
				"name":    "**ice",
				"city":    "****tle",
				"country": "USA",
			}
			resultMap = utils.MaskStringInMap(inputMap, maskChar, option, unmaskedCount)
			gomega.Expect(resultMap).To(gomega.Equal(expectedMap))
		})
	})

	Context("when option is 'fixed'", func() {
		It("should mask all values completely to the length of unmaskedCount (4)", func() {
			option = "fixed"
			unmaskedCount = 4
			expectedMap = map[string]string{
				"name":    "****",
				"city":    "****",
				"country": "****",
			}
			resultMap = utils.MaskStringInMap(inputMap, maskChar, option, unmaskedCount)
			gomega.Expect(resultMap).To(gomega.Equal(expectedMap))
		})

		It("should mask all values completely to the length of unmaskedCount (2)", func() {
			option = "fixed"
			unmaskedCount = 2
			expectedMap = map[string]string{
				"name":    "**",
				"city":    "**",
				"country": "**",
			}
			resultMap = utils.MaskStringInMap(inputMap, maskChar, option, unmaskedCount)
			gomega.Expect(resultMap).To(gomega.Equal(expectedMap))
		})
	})

	Context("when option is invalid", func() {
		It("should return an error message for all values", func() {
			option = "unknown"
			unmaskedCount = 3
			expectedMap = map[string]string{
				"name":    "XXERRORXX",
				"city":    "XXERRORXX",
				"country": "XXERRORXX",
			}
			resultMap = utils.MaskStringInMap(inputMap, maskChar, option, unmaskedCount)
			gomega.Expect(resultMap).To(gomega.Equal(expectedMap))
		})
	})
})
