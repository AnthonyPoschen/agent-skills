// Stylelint config for palette, hygiene, type, and motion rules.
// Point it at a repo:
//   createConfig({
//     themeFiles: ["web/styles/theme.css"],
//     themePath: "web/styles/theme.css",
//     prefix: "palette-",
//     deprecated: ["^--palette-accent(-|$)"],
//     colorPrefixes: ["palette-"],
//     placeholderToken: "--palette-placeholder",
//   })
const path = require("path")

const why =
  "Use a var(--token) from the theme file so recolours and theming do not miss values."

function createConfig(options = {}) {
  const themeFiles = options.themeFiles || ["**/theme.css"]
  const themeFile = options.themeFile || "theme.css"
  const themePath =
    options.themePath || path.join(__dirname, "fixtures/theme.css")
  const prefix = options.prefix || "palette-"
  const deprecated = options.deprecated || ["^--palette-accent(-|$)"]
  const colorPrefixes = options.colorPrefixes || [prefix]
  const placeholderToken = options.placeholderToken || `--${prefix}placeholder`
  const extraSelectors = options.extraSelectors || [".heading", ".section-heading"]

  const plugin = (name) => path.join(__dirname, "plugins", name)

  return {
    plugins: [
      plugin("no-hex-color.js"),
      plugin("no-old-scheme.js"),
      plugin("colors-in-theme.js"),
      plugin("no-placeholder-color.js"),
      plugin("no-clipped-heading.js"),
      plugin("no-layout-animation.js"),
      plugin("unknown-token.js"),
      plugin("no-page-component-css.js"),
      "stylelint-declaration-strict-value",
    ],
    rules: {
      "color-no-hex": [true, { message: `Hard-coded hex colour. ${why}` }],
      "color-named": ["never", { message: `Named colour. ${why}` }],
      "function-disallowed-list": [
        ["rgb", "rgba", "hsl", "hsla", "hwb", "lab", "lch", "oklab", "oklch", "color"],
        { message: `Colour function. ${why}` },
      ],
      "lint/no-hex-color": [true, { themeFile }],
      "lint/no-old-scheme": [true, { themeFile, prefix, deprecated }],
      "lint/colors-in-theme": [true, { themeFile, colorPrefixes }],
      "lint/no-placeholder-color": [true, { themeFile, placeholderToken }],
      "lint/no-clipped-heading": [true, { extraSelectors }],
      "lint/no-layout-animation": true,
      "lint/unknown-token": [true, { prefix, themeFile: themePath }],
      "lint/no-page-component-css": [
        true,
        {
          componentClasses: options.componentClasses || [
            "card",
            "btn",
            "chip",
            "row",
            "stat",
            "input",
          ],
          componentFile: options.componentFile || "components.css",
          exempt: options.exempt || ["theme.css", "tokens.css", "components.css", "layout.css"],
        },
      ],
      "declaration-no-important": [
        true,
        {
          message:
            "!important breaks the cascade (why: it overrides later, more specific rules). Remove it, or disable the rule with a comment that says why.",
        },
      ],
      "declaration-block-no-duplicate-properties": [
        true,
        {
          message:
            "Duplicate property in one block (why: the earlier declaration is dead and a later edit can miss it).",
        },
      ],
      "no-duplicate-selectors": [
        true,
        {
          message:
            "Duplicate selector (why: one rule split in two drifts when only half is updated).",
        },
      ],
      "declaration-block-no-shorthand-property-overrides": [
        true,
        {
          message:
            "Longhand before a shorthand is overwritten (why: the longhand never applies).",
        },
      ],
      "property-no-unknown": [
        true,
        {
          message:
            "Unknown property (why: the browser drops it, so the declaration does nothing).",
        },
      ],
      "declaration-property-value-allowed-list": [
        {
          "font-size": ["/^var\\(--/", "inherit"],
          "z-index": ["/^var\\(--/", "auto", "0", "-1"],
        },
        {
          message:
            "Use a var(--token) from the theme file (why: raw font-size and z-index scatter type and stacking).",
        },
      ],
      "scale-unlimited/declaration-strict-value": [
        ["/color$/", "fill", "stroke"],
        {
          ignoreValues: ["currentColor", "transparent", "inherit", "initial", "unset", "none"],
          message:
            "${property} must use a var(--token) colour (why: single source for theming).",
        },
      ],
    },
    overrides: [
      {
        files: ["**/*.html"],
        customSyntax: "postcss-html",
      },
      {
        files: themeFiles,
        rules: {
          "color-no-hex": null,
          "color-named": null,
          "function-disallowed-list": null,
          "declaration-property-value-allowed-list": null,
          "scale-unlimited/declaration-strict-value": null,
        },
      },
    ],
  }
}

module.exports = createConfig
module.exports.createConfig = createConfig
