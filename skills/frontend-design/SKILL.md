---
name: frontend-design
description: Design, build, restyle, or review polished frontend interfaces for web apps, dashboards, tools, landing pages, components, and responsive workflows. Use whenever a request involves frontend visual design, layout, styling, UI hierarchy, interface polish, design systems, color palettes, typography, color schemes, prefers-color-scheme, dark and light themes, responsive behavior, navigation and wayfinding, or making an existing web experience more usable and distinctive. Also use for motion, elevation, accessibility, and design critique. Apply it to both new screens and focused UI improvements, even when the user asks only for code. Avoid templated AI aesthetics. When designing or reviewing UI, visually inspect the running interface with the eyes-on skill, in both dark and light themes.
---

# Frontend Design

You already know how to design. These are reminders for the places where
agent-built interfaces most often fall short. The brief and the project's
`DESIGN.md` win over anything here. Make each visual choice serve the person
using the page and the page's one job.

## Before You Build

- Name the user, the outcome they need, and the page's single job. Design the
  smallest complete workflow first, including its empty, loading, error and
  success states, before you choose a shell or navigation.
- When improving existing UI, look at it running first (`eyes-on`). Keep the
  conventions that work and change the few things with the most impact.
- Ask only when an unknown would change the workflow, content priority, or
  visual direction. Otherwise state the assumption and proceed.
- Write a short plan before a new screen or redesign: the direction in a few
  concrete words, the structural system, the type voices and their jobs, the
  colour roles, and the one signature element you will spend boldness on.
- Choose the direction from the subject and the brief, not from habit. When
  the brief pins a direction, execute it fully and precisely. When an axis is
  free, ask whether your choice would appear unchanged for an unrelated
  product; if it would, find a reason in this subject or change it.

## Keep A Design Record

`DESIGN.md` in the target project is the durable record of design decisions.
Read it before a material decision. Create it at the first substantial design
consultation, and update it when feedback, testing or review changes a
decision. Record the decision and its reason, not the conversation. Use
[the `DESIGN.md` template](references/design-brief.md) for a substantial flow
or redesign. Project palettes, typefaces and signature elements belong there,
not in this skill.

When someone is choosing a palette, let them compare a few named candidates on
a real screen of the product (a local-only switcher), then keep only the
winner's tokens.

## Structure

- Give the content a structural system to sit in: a column, a grid, rules or
  a spacing rhythm, and make every section obey it, with the same spacing
  between sections and no stray empty bands. Visible structure can carry a
  page with very little decoration. Hold one shape language (corner radius,
  border weight) across every component, chips, tags and badges included.
- When items repeat (cards, rows, plans), every item shows the same fields in
  the same order, so people compare instead of re-reading.
- The primary action matches the page's job and leads to the thing the page
  is about. On a long page, repeat it at the end. Keep the first view lean:
  the page's job, the primary action and the key overview stay above the
  fold; anything else (contact, sharing, social, extra facts) goes lower.
- Give persistent information (navigation, current location, status) a
  stable position, and keep sticky or fixed chrome from covering content,
  headings or focused elements. Offset in-page targets for it.
- Keep text at a readable measure. On wide screens, use the spare width for
  navigation or deliberate space, not longer lines.
- A loading, empty or refreshed state keeps the settled state's box. Nothing
  the user can select moves when data arrives. Media reserves its aspect
  ratio.
- Feedback and indicators never reflow the text they describe.

## Type

- Use a small number of type voices, each with one job (for example display,
  reading, and a utility voice for labels, data and controls). Give each text
  role one repeatable treatment. When the brief or design record assigns a
  voice to a role (navigation, tags, buttons), apply it, casing included,
  everywhere that role appears: side navigation, tags, secondary controls and
  narrow layouts.
- Small labels that carry meaning still pass contrast. Make them quieter with
  size, weight or case, not faint grey. Keep the faintest tones for rules,
  ticks and decoration.
- When a label does not fit, shorten the words. Do not truncate meaning.
- Headings add information beyond a section's number or label rather than
  repeating it. Keep multi-line headings from stranding one short word.

## Colour

- Keep most of the surface neutral and give colour a job: action, state, or
  identity. One family per job.
- State colours mean state and stay the same on every item. Identity colours
  identify. Never let an identity colour read as a warning or error, and never
  colour state by identity.
- When sibling items each get an identity colour, choose the set together and
  check it side by side in both themes. Every pair should read as clearly
  different at a glance (aim for roughly ΔE 25), apart from the state colours,
  and saturated enough to read as a colour, not as ink. Neighbours in a
  spectrum (blue, indigo, violet) read as one family. When hue separation runs
  out, separate by lightness, shape, label or placement rather than forcing
  more hues. An item's identity (colour, type voice) carries everywhere the
  item appears: overview strips, navigation and its own card. Where it colours
  text, darken the light-theme variant until that text passes.
- Colour is never the only signal for status, selection or errors. When
  colour, shape or size encodes something that isn't self-evident, say what it
  means on the page with a short key or legend.
- Build both a light and a dark theme and follow `prefers-color-scheme`. If
  there is a theme toggle, remember the choice and switch without a flash.
- Design the dark theme rather than inverting it: lift accents so they hold
  contrast on dark, flip the primary emphasis, keep hairlines dim but
  visible, raise surfaces instead of relying on shadows, then recheck every
  label against the accessibility bar.

## Wayfinding

- On a long page or flow, help people answer: how far am I, where am I, and
  what comes next. One element can carry progress, contents and navigation
  together.
- Navigation aids carry real information and are operable. Section names in
  them are real links whose labels match the headings, the current one is
  marked (`aria-current`), and they sit early in the tab order. Any control
  that moves the reader (a progress scale, a minimap) also has a keyboard
  path. Do not draw an indicator that looks interactive and is not.
- When you build a section rail or contents list, make sure it includes:
  previous and next placed around the list, as links like the section items;
  the state of each section shown with words ("here", "next", "read") or
  proportion, not only a paler tone; every label at 4.5:1 or better; and after
  a jump, focus moved into the target section so keyboard and screen-reader
  users land where they asked to go.
- Keep numbering and names consistent between headings, navigation and
  progress.
- On narrow screens, collapse wayfinding into a compact, stable bar that keeps
  the current position, previous and next links, and a menu of every section.
  Keep progress clearly visible (a bar people can see, not a faint hairline)
  and the bar short enough not to crowd the content.

## Interaction And Motion

- Cover the states that can occur: hover, focus, active, selected, disabled,
  loading, empty, error, success. Hover, focus and selected look distinct.
  Information the user needs is in the resting state, not only on hover.
- Expose the decision the person is making, not every capability the system
  has. Put advanced options behind interaction, and keep committed choices
  visible and reversible.
- Make primary and destructive actions easy to tell apart before they are
  taken. If an action is unavailable, say why beside it.
- Motion shows cause, effect, continuity or reading order. Keep it short and
  interruptible, animate `transform` and `opacity`, and put hover effects
  under `@media (hover: hover)`. Skip motion people will see many times a
  day.
- Motion never makes the current state unreadable, not even mid-transition.
- Prefer native scrolling. Under `prefers-reduced-motion: reduce`, drop
  movement and keep the state changes.

## Copy

Words are design material. Name things by what people control and recognise,
not by how the system is built. Keep action names consistent through a flow.
Empty and error states say what happened and what to do next. A label labels;
an example demonstrates.

## Accessibility Bar

Meet this on every surface you build or change. `DESIGN.md` can raise it, not
waive it.

- Text someone must read is at least 4.5:1, including small labels, numerals
  and secondary text. Only text of 24px or more, or about 19px or more at
  weight 700+, counts as large and may drop to 3:1; 19px semibold is normal
  text. Judge size from the computed style at each width, since narrow layouts
  often shrink text below the large threshold. Coloured text on light
  backgrounds fails most often: compute its ratio in both themes at the narrow
  width. Meaningful non-text marks are at least 3:1.
- Every link and button keeps an accessible name at every width. When
  responsive CSS hides a control's text (a wordmark, a button label), give it
  an `aria-label`. Placeholder text is not a name.
- Focus is visible, follows the layout order, and is never hidden under
  sticky chrome. A keyboard user can reach every action and skip repeated
  chrome. Dialogs hold and return focus.
- Every interactive element is at least 24px, and about 44px on touch
  screens: buttons, toggles, navigation, chips, and stacked link lists such as
  footers.
- Toggles and segmented controls expose their state (`aria-pressed` or
  `aria-checked`), not only a visual change.
- Reduced motion is honoured.
- Before handoff, run an automated scan (axe or similar) at a narrow and a wide
  width in both themes, and fix every contrast and naming failure it reports.
  A clean scan is not proof; it skips text it cannot measure, so check its
  needs-review items and coloured text yourself.

## Narrow Screens

A narrow layout is a new priority order, not the desktop page squeezed until
it wraps. Keep the page's job, the primary action and the values a decision
needs in the first view. Move secondary detail into disclosures or later
sections. Collapse navigation on purpose, keeping the current location and the
way back visible. No horizontal scrolling of the page.

## Particular Screens

- **Reading pages** (articles, docs, marketing): the reading experience is the
  interface. Headings reveal the argument; quotes, examples and calls to
  action are meaningful interruptions, not interchangeable cards.
- **Forms:** group fields by the user's decisions, pair label, input, help and
  error as one unit, keep actions in a stable place, and separate destructive
  actions from the default path.
- **Data-dense screens:** give the strongest hierarchy to what must be
  understood at a glance; align values for comparison; reserve loud status
  treatment for conditions that need attention.

## Implement In Context

Follow the project's framework, components, CSS approach and tokens when they
are coherent. Use semantic HTML and tokens for repeated values. Do not add
fonts or assets whose licence or loading is unknown, and do not load families
the page never uses. Build the working interface with realistic content.

## Look At It

Markup is not visual proof. Use the `eyes-on` skill during the design loop,
not only at the end: build the smallest complete piece, look at it at desktop
and narrow widths in both themes, exercise the navigation and states, and fix
the largest problem before polishing. Look at scrolled states as well as the
top of the page, and look at repeated items and their colours side by side.

## Review Before Handoff

Fix the largest problems first, then make a detail pass. Check:

- Can a first-time visitor tell the page's purpose, primary action and
  current state at a glance?
- Does the content sit in one structural system and one shape language, and
  do repeated items share one order?
- Does every colour have one job, with state and identity kept apart, and do
  sibling identity colours stay distinct in both themes?
- Do navigation aids show real information, navigate, move focus to the
  target, and keep current, previous and next visible at every width?
- Did the automated scan pass at both widths in both themes, and does every
  control keep its name when text is hidden?
- Were both themes and both widths inspected, and does each meet the
  accessibility bar?
- Has unused type, tokens and decoration been removed?

When asked for a critique, report findings instead of fixing them. Rank each
as blocking, should-fix or polish; name who is affected, what you saw, and the
smallest fix. Blocking means someone cannot complete the task, tell the state,
or operate a control.

In the final response, state what changed, which workflow or state was
prioritised, and what visual verification was performed.
