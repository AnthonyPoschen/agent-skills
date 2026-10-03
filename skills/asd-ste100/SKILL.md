---
name: asd-ste100
description: >
  Write replies and documents in ASD-STE100 Simplified Technical English.
  Use this skill whenever you write a reply, a document, a skill, AGENTS.md,
  a comment, a commit message, or other prose for a person. Also use it when
  the user asks for a simple explanation, or says simplified English, STE,
  or ASD-STE100.
---

# ASD-STE100

Write each sentence so a reader can follow it on the first read. Follow the rules in this skill. The official word list is a licensed document. This skill keeps these rules and leaves that word list out.

## When this applies

Use these rules for prose that you write for a person. A chat reply, a document, a skill, a comment, a commit message, and a simple explanation are in this group. A simple explanation is descriptive text.

Keep the original words in these cases:

- A code identifier, a path, a command, a flag, an error string, or a quoted source.
- A name that the project already uses. Repeat that name.

Write a long project name in full one time. After that, use a shorter name only when the project already has that shorter name.

When `technical-writing` also applies, that skill selects the document type. This skill controls each sentence.

## Select the sentence type

A procedure tells the reader to do a task. A description gives information. Keep each paragraph to one type.

- A procedure sentence has a maximum of 20 words.
- A description sentence has a maximum of 25 words.
- A warning has a maximum of 20 words.
- A caution has a maximum of 20 words.
- A note gives information only. A note has a maximum of 25 words. Do not put a command in a note.
- Text that prevents harm is a warning or a caution.

Count these items as one word each: an identifier, a path, a command, a hyphenated name, and a number with its unit. Count each list item as its own sentence.

Split a sentence that is longer than its limit. Keep every fact.

## Write the sentence

- Put one action in each procedure sentence. Two actions can share a sentence only when the actions happen at the same time.
- Put one topic in each description sentence.
- Start a procedure sentence with a command verb.
- Put a condition before the action. Example: "If the file is absent, stop."
- Use the active voice. Name the actor, or use a command.
- In a description, use the passive voice only when the actor is unknown.
- Use only these verb forms: infinitive, command, simple present, simple past, and simple future.
- Use a past participle only as an adjective. The adjective shows a condition.
- An "-ing" word is a technical noun, or part of a technical name.
- Use "must" for an action the reader is required to do. Use "can" for an ability. Use "will" for the future.
- Write a recommendation as its own sentence. Do not use "should", "may", "might", "would", or "shall".
- Write "a", "an", or "the" when English needs an article.
- Write the subject and the verb in each sentence.
- Keep "that" after "make sure" and after "show".
- Repeat the noun when a pronoun can mean two nouns.
- Write each word in full. Example: "do not".
- Write a count or a measure as a numeral.
- Use American spelling.

## Use one word for one meaning

Use the short ordinary word when the meaning stays the same. Use the project name or the project verb when the ordinary word hides the meaning.

- Assign one word to each meaning in the text.
- Use that word every time you write that meaning.
- Choose the short ordinary word. Examples: "use", "start", "help", "before", "stop", and "remove".
- Write "for example" when you give an example.
- Write "and other items" when a list continues.
- A noun cluster has a maximum of three nouns.
- Three nouns is the ceiling. Prefer fewer nouns when the meaning stays clear.
- Use a preposition when the group needs more than three nouns.
- Use one verb for one action. Example: "remove the cover".
- Keep a phrasal verb when that verb is the project term. Examples: "roll back" and "check out".
- Put the action in the verb. Example: "install the seal".
- Replace slang with the action. Keep the term that the project uses.

## Organize the text

- Use a vertical list for several items or several actions.
- Introduce the list with a complete sentence and a colon.
- Put a warning or a caution immediately before the related step.
- A warning is a risk of injury, or an action that the reader cannot undo.
- A caution is a risk of damage to an object or a system.
- Start the warning or the caution with the command. Then name the risk.
- Keep the risk sentence in that same warning or caution.
- Give each prose paragraph one topic and a maximum of six sentences.
- A vertical list is not one paragraph. Each list item still counts as its own sentence.
- End each statement with a period.
- Join two alternatives with "and" or with "or".
- Write "and/or" only when both readings are true.

## Check before you send

1. Count the words with the rules in "Select the sentence type". A procedure sentence has 20 words or fewer. A description sentence has 25 words or fewer.
2. Each procedure sentence has one action, or two actions at the same time.
3. Each paragraph has one type, one topic, and a maximum of six sentences.
4. A requirement uses "must". An ability uses "can". A future fact uses "will". A recommendation is its own sentence. The text has no "should", "may", "might", "would", or "shall".
5. Each thing has the same name in the whole text.
6. Each code name, path, command, and quote keeps its original form.
7. Each warning and each caution starts with the command.
