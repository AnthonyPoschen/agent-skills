---
name: technical-writing
description: >
  Use when writing or reviewing documentation, READMEs, tutorials, how to
  guides, RFCs, architecture notes, PR descriptions, release notes, or commit
  messages. Write clear, plain English for developers and maintainers.
---

# Technical Writing

Use this skill to write for a tired developer who should understand the first
read. It is adapted from pstack's technical writing skill under its MIT license,
but stands alone and has no pstack only dependencies.

## Write Plain English

- Cut words that do no work. Write “to”, not “in order to”. Delete “it is
  important to note that”.
- Use short, ordinary words unless a technical term adds needed precision. Write
  “use”, not “utilize”; “move”, not “relocate”; and “help”, not “facilitate”.
- Write the real symbol, file path, flag, command, type, or field name. Do not
  replace it with a vague description.
- Use one name for one thing throughout the document.
- Do not invent jargon, metaphors, or made up hyphenated phrases. Use words a
  developer would say out loud.
- Prefer a period to an em dash. Use a new sentence when the thought changes.
- Keep related words close. Put “only” and “not” next to the word they change.
- Break up dense noun strings. Write “the script that checks imports”, not “the
  import budget check script”.

The rules serve the reader. Rewrite a sentence that sounds mechanical instead
of following a rule in a way that makes the sentence worse.

## Choose One Document Mode

Pick one mode for each document. Split and link documents when their purposes
conflict.

- **Tutorial:** teach a newcomer by guiding them through a visible result.
- **How to guide:** help a capable reader complete one task. Keep only the steps
  and decisions needed for that task.
- **Reference:** describe facts, options, limits, inputs, outputs, and errors.
  Mirror the thing being described. Do not persuade or teach.
- **Explanation:** answer a bounded “why” question with context, constraints,
  alternatives, and tradeoffs.

## Write Clear Sentences

- Address the reader as “you” when giving instructions.
- Use active voice when the actor matters. State who does what.
- Write instructions as commands. Put a condition before the step it controls.
- Put the common case first. Put exceptions after it.
- Give each sentence one instruction or one main thought. Split a sentence that
  carries more.
- Keep articles and verbs when they make the sentence easier to parse.
- Make each “it”, “they”, “this”, and “which” point to one obvious noun. Repeat
  the noun when that is clearer.
- Use “both … and”, “either … or”, and “if … then” when they prevent ambiguity.
- Avoid passive constructions such as “the file should be updated”. Write “update
  the file”.
- Read awkward text aloud. Rewrite it if it still sounds awkward.

## Organize For Scanning

- Use headings that state the point. Write “Choose the document mode”, not
  “Document modes”.
- Use sentence case headings. Use a verb phrase for a task heading and a noun
  phrase for a concept heading.
- Use numbered lists for ordered steps. Use bullet lists for unordered facts.
- Introduce every list with a complete sentence and keep its items parallel.
- Link with words that name the destination. Do not write “click here”.
- Use code formatting for real paths, commands, flags, symbols, and output.
- Keep examples short and based on a real reader task.

## Apply Repository Rules

- Follow the repository's established documentation and release conventions.
- Keep claims about paths, symbols, counts, and behavior true for the current
  change. Run the command that proves a generated count or tree when one exists.
- Treat PR descriptions and commit messages as technical writing.
- Do not apply this skill to product UI text. Use the product's copy rules for
  that work.

## Review Checklist

Before finishing, check the following:

1. The document has one clear mode and purpose.
2. Each instruction states who acts and what they do.
3. Each sentence has one main thought.
4. Each concept has one consistent name.
5. The document uses ordinary words and real repository symbols.
6. Paths, commands, counts, and examples are current and verifiable.
7. The reader can find the task, fact, or explanation they need by scanning the
   headings.

## Sources

- [Pstack technical writing skill](https://github.com/cursor/plugins/blob/main/pstack/skills/technical-writing/SKILL.md)
- [Diátaxis](https://diataxis.fr/)
- [Google developer documentation style guide](https://developers.google.com/style)
