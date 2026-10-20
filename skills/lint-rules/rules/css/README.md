# CSS palette and hygiene rules

Colours, type size, and stacking live in one theme file. Copy
[create-config.cjs](create-config.cjs) and set `themeFiles`, `themeFile`,
`themePath`, `prefix`, `deprecated`, and `colorPrefixes`.

```js
module.exports = require("./lint/css/create-config.cjs")({
  themeFiles: ["web/styles/theme.css"],
  themeFile: "web/styles/theme.css",
  themePath: "web/styles/theme.css",
  prefix: "palette-",
  deprecated: ["^--palette-accent(-|$)"],
  placeholderToken: "--palette-placeholder",
})
```

```sh
npm install -D stylelint@16.23.1 postcss-html@1.8.0 stylelint-declaration-strict-value@1.10.11
npx stylelint "src/**/*.css" "web/**/*.css" "web/**/*.html"
```

Plugins live in [plugins/](plugins/). `inline-color.js` and `svg-palette.js`
scan HTML, JS, and SVG; call them from the same CI step as Stylelint. Every
message says why the rule exists.

| rule | bad |
| --- | --- |
| `color-no-hex`, `color-named`, `function-disallowed-list`, `lint/no-hex-color`, `scale-unlimited/declaration-strict-value` | raw colour outside the theme file |
| `lint/no-old-scheme` | a token on the configurable deprecated list |
| `lint/colors-in-theme` | a colour custom property defined outside the theme file |
| `lint/no-placeholder-color` | `::placeholder` colour other than the shared token |
| `lint/no-clipped-heading` | fixed height or width with `overflow: hidden`, or `line-height` under 1.15, on a heading |
| `lint/no-layout-animation` | animating width, height, insets, margin, or padding |
| `lint/unknown-token` | `var(--prefix-*)` missing from the theme file |
| `lint/no-page-component-css` | component class or raw colour in page CSS |
| `declaration-no-important`, duplicate properties, duplicate selectors, shorthand overrides, `property-no-unknown` | cascade and dead-declaration bugs |
| `declaration-property-value-allowed-list` | raw `font-size` or `z-index` |
| inline-color, svg-palette | colour in `style=`, `element.style`, canvas paint, or an SVG fill outside the theme hex set |

```sh
npm --prefix skills/lint-rules/rules/css test
```
