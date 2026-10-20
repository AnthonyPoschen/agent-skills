const stylelint = require("stylelint")

const ruleName = "lint/no-layout-animation"
const messages = stylelint.utils.ruleMessages(ruleName, {
  layout: (prop) =>
    `Do not animate ${prop}. Width, height, insets, margin, and padding shift layout (why: that animation causes cumulative layout shift). Animate transform, opacity, filter, or colour.`,
  all: "transition: all animates layout properties and shifts the page (why: cumulative layout shift). Name transform, opacity, filter, or a colour.",
  bare: "This transition shorthand names no property, so it defaults to all and shifts layout (why: cumulative layout shift). Name transform, opacity, filter, or a colour.",
})

const layoutProp = /^(width|height|top|left|right|bottom|margin|padding)(-|$)/
const timeToken = /^-?\d*\.?\d+m?s$/
const namedEasing = /^(ease|ease-in|ease-out|ease-in-out|linear|step-start|step-end)$/
const globalKeyword = /^(none|inherit|initial|unset|revert|revert-layer)$/

const rule = (primary) => {
  return (root, result) => {
    const validOptions = stylelint.utils.validateOptions(result, ruleName, { actual: primary })
    if (!validOptions || primary === false) return

    root.walkDecls((decl) => {
      if (insideKeyframes(decl) && layoutProp.test(decl.prop)) {
        report(decl, messages.layout(decl.prop), result)
      }
      if (decl.prop === "transition" || decl.prop === "transition-property") {
        checkTransition(decl, result)
      }
    })
  }
}

function insideKeyframes(node) {
  let parent = node.parent
  while (parent) {
    if (parent.type === "atrule" && parent.name === "keyframes") return true
    parent = parent.parent
  }
  return false
}

function checkTransition(decl, result) {
  for (const part of splitCommas(decl.value)) {
    const trimmed = part.trim()
    if (!trimmed) continue
    if (/\ball\b/.test(trimmed)) {
      report(decl, messages.all, result)
      continue
    }
    const tokens = tokenize(trimmed)
    const properties = tokens.filter((token) => !isTiming(token) && !globalKeyword.test(token))
    if (properties.length === 0 && !globalKeyword.test(trimmed)) {
      report(decl, messages.bare, result)
      continue
    }
    for (const prop of properties) {
      if (layoutProp.test(prop)) report(decl, messages.layout(prop), result)
    }
  }
}

function splitCommas(value) {
  const parts = []
  let start = 0
  let depth = 0
  for (let i = 0; i < value.length; i++) {
    const char = value[i]
    if (char === "(") depth++
    else if (char === ")") depth = Math.max(0, depth - 1)
    else if (char === "," && depth === 0) {
      parts.push(value.slice(start, i))
      start = i + 1
    }
  }
  parts.push(value.slice(start))
  return parts
}

function tokenize(value) {
  const tokens = []
  let depth = 0
  let current = ""
  for (const char of value) {
    if (char === "(") depth++
    else if (char === ")") depth = Math.max(0, depth - 1)
    if (depth === 0 && /\s/.test(char)) {
      if (current) tokens.push(current)
      current = ""
      continue
    }
    current += char
  }
  if (current) tokens.push(current)
  return tokens
}

function isTiming(token) {
  if (timeToken.test(token)) return true
  if (namedEasing.test(token)) return true
  return /^(cubic-bezier|steps)\(/.test(token)
}

function report(decl, message, result) {
  stylelint.utils.report({ message, node: decl, result, ruleName })
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
