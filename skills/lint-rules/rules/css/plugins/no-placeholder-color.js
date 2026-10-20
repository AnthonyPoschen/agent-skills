const stylelint = require("stylelint")

const ruleName = "lint/no-placeholder-color"
const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected: (themeFile, token) =>
    `Placeholder colour belongs to the single ::placeholder rule in ${themeFile}. Use var(${token}) there (why: per-field placeholder colours drift from the palette).`,
})

const placeholderSelector =
  /::placeholder|::-webkit-input-placeholder|::-moz-placeholder|:-ms-input-placeholder/

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
          placeholderToken: (value) => typeof value === "string" && value.startsWith("--"),
        },
        optional: true,
      },
    )
    if (!validOptions || primary === false) return

    const options = secondaryOptions || {}
    const themeFile = options.themeFile || "theme.css"
    const token = options.placeholderToken || "--placeholder"
    const file = ((root.source && root.source.input && root.source.input.file) || "").replaceAll("\\", "/")
    const inTheme = file.endsWith(themeFile.replaceAll("\\", "/"))

    root.walkRules((ruleNode) => {
      if (!placeholderSelector.test(ruleNode.selector || "")) return
      ruleNode.walkDecls((decl) => {
        const prop = decl.prop.toLowerCase()
        const colourProp = prop === "color" || prop === "-webkit-text-fill-color"
        if (!colourProp && prop !== "opacity") return
        const sharedColour =
          inTheme && colourProp && decl.value.replace(/\s+/g, "") === `var(${token})`
        const sharedOpacity = inTheme && prop === "opacity"
        if (sharedColour || sharedOpacity) return
        stylelint.utils.report({
          message: messages.rejected(themeFile, token),
          node: decl,
          result,
          ruleName,
        })
      })
    })
  }
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
