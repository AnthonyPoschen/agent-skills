const stylelint = require("stylelint")

const ruleName = "lint/no-old-scheme"
const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected: (name) =>
    `Deprecated palette token ${name}. Use a current token from the theme file so one palette stays the source of colour (why: leftover schemes survive a recolour).`,
})

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
          prefix: (value) => typeof value === "string",
          deprecated: (value) => typeof value === "string" && value.length > 0,
        },
        optional: true,
      },
    )
    if (!validOptions || primary === false) return

    const options = secondaryOptions || {}
    const themeFile = options.themeFile || "theme.css"
    const prefix = options.prefix == null ? "" : options.prefix
    const deprecated = (options.deprecated || []).map((source) => new RegExp(source))
    if (deprecated.length === 0) return

    const file = (root.source && root.source.input && root.source.input.file) || ""
    if (file.replaceAll("\\", "/").endsWith(themeFile.replaceAll("\\", "/"))) return

    const matcher = new RegExp(`--${escapeRegExp(prefix)}[A-Za-z0-9_-]+`, "g")
    root.walkDecls((decl) => {
      const seen = new Set()
      for (const chunk of [decl.prop, decl.value]) {
        for (const match of String(chunk).matchAll(matcher)) {
          const name = match[0]
          if (seen.has(name) || !deprecated.some((pattern) => pattern.test(name))) continue
          seen.add(name)
          stylelint.utils.report({
            message: messages.rejected(name),
            node: decl,
            result,
            ruleName,
            word: name,
          })
        }
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
