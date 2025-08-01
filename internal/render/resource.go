package render

import (
	_ "embed"
	"encoding/base64"
	"html/template"
)

var (
	//go:embed templates/script/echarts5.6.0.min.js
	resourceEcharts string
	//go:embed templates/script/tailwindcss.js
	resourceTailwind string
	//go:embed templates/font/zsft184.woff2
	resourceZsft184 []byte
)

var (
	ResourceEcharts = template.JS(resourceEcharts)

	ResourceTailwind = template.JS(resourceTailwind)

	ResourceZsft184 = template.CSS(`@font-face {
        font-family: 'ZSFT-ENMIN-184';
        src: url(data:font/woff2;base64,` + base64.StdEncoding.EncodeToString(resourceZsft184) + `) format('woff2');
        unicode-range: U+0061-007A, U+0041-005A, U+0030-0039, U+FF0C, U+3002, U+002E, U+002C, U+0021, U+003F, U+003A, U+003B, U+002D, U+FF01, U+FF1F;
    }`)
)
