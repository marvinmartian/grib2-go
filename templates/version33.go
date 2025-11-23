package templates

import (
	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
)

// Version33 returns templates that are compatible with version 33 of the GRIB2 specification.
func Version33() record.Templates {
	return &templates{
		gridDefinitionEndingOctet: map[int]RetrieveEndOctet{
			0:     constantEnd(72),
			1:     constantEnd(84),
			2:     constantEnd(84),
			3:     constantEnd(96),
			4:     constantEnd(48),
			5:     constantEnd(60),
			10:    constantEnd(72),
			12:    constantEnd(84),
			13:    constantEnd(88),
			20:    constantEnd(65),
			23:    constantEnd(81),
			30:    constantEnd(81),
			31:    constantEnd(81),
			33:    constantEnd(97),
			40:    constantEnd(72),
			41:    constantEnd(84),
			42:    constantEnd(84),
			43:    constantEnd(96),
			50:    constantEnd(28),
			51:    constantEnd(40),
			52:    constantEnd(40),
			53:    constantEnd(52),
			61:    constantEnd(112),
			62:    constantEnd(106),
			63:    constantEnd(121),
			90:    constantEnd(80),
			100:   constantEnd(38),
			101:   constantEnd(35),
			110:   constantEnd(57),
			140:   constantEnd(64),
			150:   constantEnd(42),
			204:   constantEnd(72),
			1100:  constantEnd(82),
			32768: constantEnd(72),
			32769: constantEnd(80),
		},
		productDefinitionEndingOctet: map[int]RetrieveEndOctet{
			0: constantEnd(34),
			1: constantEnd(37),
			2: constantEnd(36),
			3: func(bytes []byte) int {
				forecasts := int(bytes[57])
				end := 58 + forecasts
				return end
			},
			4: func(bytes []byte) int {
				forecasts := int(bytes[53])
				end := 64 + forecasts
				return end
			},
			5: constantEnd(47),
			6: constantEnd(35),
			7: constantEnd(34),
			8: func(bytes []byte) int {
				timeRanges := int(bytes[41])
				end := 46 + (12 * timeRanges)
				return end
			},
			9: func(bytes []byte) int {
				timeRanges := int(bytes[54])
				end := 59 + (12 * timeRanges)
				return end
			},
		},
		gridDefinitionTemplates: map[int]record.GridDefinition{
			0:  grid.Template0{},
			1:  grid.Template1{},
			30: grid.Template30{},
			40: grid.Template40{},
		},
		productDefinitionTemplates: map[int]record.ProductDefinition{
			0: product.Template0{},
			1: product.Template1{},
			2: product.Template2{},
			3: product.Template3{},
			4: product.Template4{},
			5: product.Template5{},
			6: product.Template6{},
			7: product.Template7{},
			8: product.Template8{},
		},
		dataRepresentationTemplates: map[int]record.DataRepresentationDefinition{
			0:  data_representation.Template0{},
			2:  data_representation.Template2{},
			3:  data_representation.Template3{},
			40: data_representation.Template40{},
			41: data_representation.Template41{},
		},
	}
}

func constantEnd(end int) RetrieveEndOctet {
	return func(bytes []byte) int {
		return end
	}
}
