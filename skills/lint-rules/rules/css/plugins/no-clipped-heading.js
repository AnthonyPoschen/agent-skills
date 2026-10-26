const stylelint = require("stylelint")

const ruleName = "lint/no-clipped-heading"
const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected: (reason) =>
    `Heading text is clipped (${reason}). Use min-height, line-height >= 1.15, and visible overflow (why: a fixed box with overflow:hidden cuts glyphs).`,
})

function declMap(ruleNode) {
  const map = new Map()
  ruleNode.walkDecls((decl) => {
    map.set(decl.prop.toLowerCase(), decl)
  })
  return map
}

function isFixedLength(value) {
  return /^-?\d*\.?\d+(px|rem|em)$/.test((value || "").trim())
}

const rule = (primary, secondaryOptions) => {
  return (root, result) => {
    const validOptions = stylelint.utils.validateOptions(
      result,
      ruleName,
      { actual: primary },
      {
        actual: secondaryOptions,
        possible: {
          extraSelectors: (value) => typeof value === "string" && value.length > 0,
        },
        optional: true,
      },
    )
    if (!validOptions || primary === false) return

    const extras = (secondaryOptions && secondaryOptions.extraSelectors) || [
      ".heading",
      ".section-heading",
    ]
    const headingSelector = new RegExp(
      `(^|[\\s>+~,])h[1-6]\\b|${extras.map(escapeRegExp).join("|")}`,
    )

    root.walkRules((ruleNode) => {
      if (!headingSelector.test(ruleNode.selector || "")) return
      const decls = declMap(ruleNode)
      if (decls.has("clip") || decls.has("clip-path")) return
      const width = decls.get("width")
      const height = decls.get("height")
      if (width && width.value.trim() === "1px" && height && height.value.trim() === "1px") return

      const lineHeight = decls.get("line-height")
      if (lineHeight && /^\d*\.?\d+$/.test(lineHeight.value.trim())) {
        const leading = Number(lineHeight.value)
        if (leading < 1.15) {
          stylelint.utils.report({
            message: messages.rejected(`line-height: ${lineHeight.value}`),
            node: lineHeight,
            result,
            ruleName,
          })
        }
      }

      const overflow = decls.get("overflow")
      if (!overflow || overflow.value.trim() !== "hidden") return
      if (decls.has("text-overflow")) return
      for (const prop of ["height", "max-height", "width", "max-width"]) {
        const decl = decls.get(prop)
        if (!decl || !isFixedLength(decl.value) || decl.value.trim() === "1px") continue
        stylelint.utils.report({
          message: messages.rejected(`${prop}: ${decl.value} with overflow: hidden`),
          node: decl,
          result,
          ruleName,
        })
      }
    })
  }
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
