const fs = require("fs")
const path = require("path")
const stylelint = require("stylelint")

const ruleName = "lint/unknown-token"
const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected: (name, themeFile) =>
    `var(${name}) is not defined in ${themeFile}. Add the token there (why: an unknown token is a silent fallback, including when var() has a fallback).`,
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
          prefix: (value) => typeof value === "string" && value.length > 0,
          themeFile: (value) => typeof value === "string" && value.length > 0,
        },
        optional: false,
      },
    )
    if (!validOptions || primary === false) return

    const prefix = secondaryOptions.prefix
    const themeFile = secondaryOptions.themeFile
    const defined = definedTokens(themeFile, prefix)
    const file = (root.source && root.source.input && root.source.input.file) || ""
    if (file.replaceAll("\\", "/").endsWith(themeFile.replaceAll("\\", "/"))) return

    root.walkDecls((decl) => {
      for (const name of customProperties(decl.value, prefix)) {
        if (defined.has(name)) continue
        stylelint.utils.report({
          message: messages.rejected(name, themeFile),
          node: decl,
          result,
          ruleName,
          word: name,
        })
      }
    })
  }
}

function definedTokens(themeFile, prefix) {
  const abs = path.isAbsolute(themeFile) ? themeFile : path.join(process.cwd(), themeFile)
  const text = fs.readFileSync(abs, "utf8")
  const names = new Set()
  const matcher = new RegExp(`--${escapeRegExp(prefix)}[A-Za-z0-9_-]*`, "g")
  for (const match of text.matchAll(matcher)) names.add(match[0])
  return names
}

function customProperties(value, prefix) {
  const names = []
  const matcher = new RegExp(`var\\(\\s*(--${escapeRegExp(prefix)}[A-Za-z0-9_-]*)`, "g")
  for (const match of String(value).matchAll(matcher)) names.push(match[1])
  return names
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
