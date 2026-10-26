const stylelint = require("stylelint")

const ruleName = "lint/colors-in-theme"
const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected: (prop, themeFile) =>
    `Colour custom property ${prop} is defined outside ${themeFile}. Define it in the theme file (why: a second palette definition drifts from the source of colour).`,
})

const colorLiteral =
  /#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{4}|[0-9a-fA-F]{3})\b|\b(?:rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch)\(|\b(?:white|black|red|green|blue|gray|grey|yellow|orange|purple|pink|brown|cyan|magenta|silver|maroon|navy|teal|olive|lime|aqua|fuchsia)\b|color-mix\(/i

const rule = (primary, secondaryOptions) => {
  return (root, result) => {
    const validOptions = stylelint.utils.validateOptions(
      result,
      ruleName,
      { actual: primary },
      {
        actual: secondaryOptions,
        possible: {
          themeFile: (value) => typeof value === "string" && value.length > 0,
          colorPrefixes: (value) => typeof value === "string",
        },
        optional: true,
      },
    )
    if (!validOptions || primary === false) return

    const options = secondaryOptions || {}
    const themeFile = options.themeFile || "theme.css"
    const prefixes = options.colorPrefixes || []
    const file = (root.source && root.source.input && root.source.input.file) || ""
    if (file.replaceAll("\\", "/").endsWith(themeFile.replaceAll("\\", "/"))) return

    const colorVar = prefixes.length
      ? new RegExp(
          `var\\(\\s*--(?:${prefixes.map(escapeRegExp).join("|")})`,
        )
      : null

    root.walkDecls((decl) => {
      if (!decl.prop.startsWith("--")) return
      const value = decl.value || ""
      if (!colorLiteral.test(value) && !(colorVar && colorVar.test(value))) return
      stylelint.utils.report({
        message: messages.rejected(decl.prop, themeFile),
        node: decl,
        result,
        ruleName,
        word: decl.prop,
      })
    })
  }
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
