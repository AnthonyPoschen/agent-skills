const stylelint = require("stylelint")

const ruleName = "lint/no-page-component-css"
const messages = stylelint.utils.ruleMessages(ruleName, {
  color: (value) =>
    `Unexpected colour ${value} in page CSS. Use a token from the theme file (why: page-local colours are missed by a recolour).`,
  component: (name, componentFile) =>
    `Unexpected component class .${name}. Define it in ${componentFile} (why: the same component styled in a page and a stylesheet drifts).`,
})

const hexColor = /#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{4}|[0-9a-fA-F]{3})\b/g
const colorFunction = /\b(?:rgba?|hsla?|hwb|lab|lch|oklab|oklch|color)\(/gi
const namedColor =
  /\b(?:black|white|red|blue|green|gray|grey|silver|navy|orange|yellow|purple|pink|brown|cyan|magenta|lime|maroon|olive|teal|aqua|fuchsia)\b/gi
const colorProperty =
  /^(?:color|background(?:-color|-image)?|border(?:-(?:top|right|bottom|left|inline|block))?(?:-color)?|outline(?:-color)?|fill|stroke|box-shadow|text-shadow|caret-color|column-rule-color|text-decoration-color|accent-color)$/

const rule = (primary, secondaryOptions) => {
  return (root, result) => {
    const validOptions = stylelint.utils.validateOptions(
      result,
      ruleName,
      { actual: primary },
      {
        actual: secondaryOptions,
        possible: {
          componentClasses: (value) => typeof value === "string" && value.length > 0,
          componentFile: (value) => typeof value === "string" && value.length > 0,
          exempt: (value) => typeof value === "string" && value.length > 0,
        },
        optional: true,
      },
    )
    if (!validOptions || primary === false) return

    const options = secondaryOptions || {}
    const classes = options.componentClasses || ["card", "btn", "chip", "row", "stat", "input"]
    const componentFile = options.componentFile || "components.css"
    const exempt = options.exempt || ["theme.css", "tokens.css", "components.css", "layout.css"]
    const componentClass = new RegExp(`\\.(${classes.map(escapeRegExp).join("|")})(?![-\\w])`, "g")

    const file = (root.source && root.source.input && root.source.input.file) || ""
    if (isExempt(file, exempt)) return

    root.walkRules((ruleNode) => {
      if (ruleNode.parent && ruleNode.parent.type === "atrule" && ruleNode.parent.name === "keyframes") return
      const selector = ruleNode.selector || ""
      for (const match of selector.matchAll(componentClass)) {
        stylelint.utils.report({
          message: messages.component(match[1], componentFile),
          node: ruleNode,
          index: match.index,
          endIndex: match.index + match[0].length,
          result,
          ruleName,
          word: match[0],
        })
      }
    })

    root.walkDecls((decl) => {
      const value = decl.value || ""
      const reportValue = (match) => {
        stylelint.utils.report({
          message: messages.color(match[0]),
          node: decl,
          index: match.index,
          endIndex: match.index + match[0].length,
          result,
          ruleName,
          word: match[0],
        })
      }
      for (const match of value.matchAll(hexColor)) reportValue(match)
      for (const match of value.matchAll(colorFunction)) reportValue(match)
      if (!colorProperty.test(decl.prop)) return
      for (const match of value.matchAll(namedColor)) reportValue(match)
    })
  }
}

function isExempt(file, exempt) {
  const normalized = file.replaceAll("\\", "/")
  const base = normalized.slice(normalized.lastIndexOf("/") + 1)
  return exempt.some((name) => base === name || normalized.endsWith("/" + name))
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
