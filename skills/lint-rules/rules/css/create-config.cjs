// Build a Stylelint config that bans hard-coded colours outside the theme
// file. Pass themeFiles to name the exempt palette file(s), e.g.
//   createConfig({ themeFiles: ["web/styles/theme.css"] })
// Default: any path ending in theme.css.
const why =
  "Use a var(--token) from the theme file so recolours and theming do not miss values."

function createConfig(options = {}) {
  const themeFiles = options.themeFiles || ["**/theme.css"]
  return {
    rules: {
      "color-no-hex": [
        true,
        { message: `Hard-coded hex colour. ${why}` },
      ],
      "color-named": [
        "never",
        { message: `Named colour. ${why}` },
      ],
      "function-disallowed-list": [
        [
          "rgb",
          "rgba",
          "hsl",
          "hsla",
          "hwb",
          "lab",
          "lch",
          "oklab",
          "oklch",
          "color",
        ],
        { message: `Colour function. ${why}` },
      ],
    },
    overrides: [
      {
        files: themeFiles,
        rules: {
          "color-no-hex": null,
          "color-named": null,
          "function-disallowed-list": null,
        },
      },
    ],
  }
}

module.exports = createConfig
module.exports.createConfig = createConfig
