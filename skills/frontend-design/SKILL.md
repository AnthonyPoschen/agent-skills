---
name: frontend-design
description: Design, build, restyle, or review polished frontend interfaces for web apps, dashboards, tools, landing pages, components, and responsive workflows. Use whenever a request involves frontend visual design, layout, styling, UI hierarchy, interface polish, design systems, color palettes, typography, color schemes, prefers-color-scheme, dark and light themes, responsive behavior, or making an existing web experience more usable and distinctive. Apply it to both new screens and focused UI improvements, even when the user asks only for code. Avoid templated AI aesthetics. When designing or reviewing UI, visually inspect the running interface in both dark and light themes with browser-use or computer-use tools.
---

# Frontend Design

Create interfaces that make a real task easy to understand and complete. Make
visual choices serve the product and its users; do not treat a collection of
components as a finished design.

## Design Consultation

When consulted before implementation, establish the few decisions that coding
cannot safely infer. Do not block a small, reversible change for a full brief;
use repository evidence and state reasonable assumptions. Ask only when an
unknown would materially change the workflow, content priority, or visual
direction.

Treat `DESIGN.md` in the target project as the durable record of design
knowledge. Look for it before making a material design decision. At the first
substantive design consultation, create it if it does not exist. Update it when
user feedback, repository evidence, testing, or a design review changes an
assumption, constraint, priority, or decision. Record the decision and its
reason, not a transcript of the conversation.

Capture these decisions when the work is substantial, likely to span multiple
sessions, or needs handoff to another implementer:

- **Outcome and success:** the user outcome, the primary action or decision,
  and how the team will know the design worked.
- **Audience and context:** who uses it, their familiarity and urgency, the
  device or environment, and any access needs that shape the interaction.
- **Primary workflow:** entry point, happy path, key alternatives, and the
  feedback or recovery needed after success, failure, or interruption.
- **Content and data:** the real information, its relative importance, expected
  density, long or empty values, and any content that must remain visible.
- **States and rules:** permissions, validation, loading, empty, error,
  destructive, selected, and time-sensitive states that change what a person
  can see or do.
- **Constraints:** existing design system or product conventions, technical
  limits, supported viewports and input methods, accessibility requirements,
  brand constraints, and licensed assets or fonts.
- **Visual direction:** a few concrete adjectives, comparable product
  references, and anti-references that clarify the desired character without
  copying another product.
- **Scope and tradeoffs:** what is deliberately out of scope, what can be a
  sensible default, and which unresolved questions need a prototype or user
  decision.

Use [the `DESIGN.md` template](references/design-brief.md) to add a brief for
each substantial flow or redesign. Keep multiple related designs in the root
`DESIGN.md` while it remains easy to scan. Once it grows beyond about 1,000
lines or develops independently navigable areas, make `DESIGN.md` a concise
router and move individual briefs to `design/<topic>.md`. Keep shared design
principles and the linked index in the root file. Do not turn the record into a
retrospective diary.

## Start With The Product

Before choosing an app shell or a visual style, establish:

- the user and the outcome they need;
- the smallest complete workflow or feature being designed;
- the information needed to act, decide, or verify success;
- the important content states: normal, empty, loading, validation, error, and
  success; and
- the device, viewport, input method, and accessibility constraints that matter.

Start with a useful feature or state, not a navigation layout. Let repeated
workflows reveal whether the product needs a sidebar, top navigation, dense
table, wizard, canvas, feed, or another structure.

When improving existing UI, inspect the current behavior and visual language
first, including the rendered page when a browser-use or computer-use tool is
available. Preserve useful conventions and identify the few changes with the
most impact on comprehension, efficiency, and confidence.

## Establish A Clear Direction

Choose a deliberate visual character appropriate to the product. A repeated-use
tool should usually be calm, dense, and quick to scan; an editorial or
brand-focused experience can be more expressive. State the intended direction
in a few concrete terms before implementation, such as "quiet financial
workbench", "warm local guide", or "precise technical console".

Use that direction to constrain decisions:

- Pick a small, purposeful type system, spacing scale, color system, radius
  family, and elevation model.
- Reuse tokens and component patterns. Variation should signal a meaningful
  difference, not indecision.
- Use real-looking content, labels, data, and edge cases. Placeholder-only
  layouts hide the density and hierarchy problems users will actually face.
- Make decoration earn its place by reinforcing the product, a state, or the
  reading order.

Avoid generic visual filler: unmotivated gradients, floating glass panels,
oversized rounded cards, decorative glow, or ornamental motion that does not
help the task.

## Make It Distinctive, Not Default

Treat visual identity as a product decision. A usable layout that could belong
to any SaaS is unfinished. Ground choices in the subject's world: materials,
instruments, vernacular, and the page's single job.

Before coding a new screen or a visual redesign, write a compact plan:

- **Color:** 4–6 named hex values with roles (page, surface, ink, accent,
  semantic). Do not invent extra accents while implementing.
- **Type:** a characterful display face used with restraint, a complementary
  body face, and a utility face for captions or data if needed. Do not load
  families that never appear in the UI.
- **Layout:** one-sentence concept plus a small ASCII wireframe when comparing
  options.
- **Signature:** the one memorable element this surface should be remembered
  by. Spend boldness there; keep the rest quiet.

Current AI-generated work clusters around three looks: (1) warm cream
(`#F4F1EA`) with a high-contrast serif and terracotta; (2) near-black with a
single acid-green or vermilion accent; (3) broadsheet hairline rules, zero
radius, dense newspaper columns. All three are legitimate for some briefs,
but they are defaults rather than choices. Where the brief pins a direction,
follow it. Where an axis is free, do not spend that freedom on one of these
clusters.

If a palette is "dark gray plus one mint/teal/green accent", or "cream plus
terracotta", name the subject-specific reason it belongs here or replace it.
Project-specific palettes, typefaces, and signature elements belong in that
project's `DESIGN.md`, not in this skill.

Review the plan against the brief before building. If any part would appear
for an unrelated product, revise that part and say why.

## Write Copy As Design Material

Words exist to make the interface easier to understand. Name things by what
people control and recognize, never by how the system is built. Use active
voice. Keep action names consistent through the flow ("Save changes" →
"Saved"). Failure and empty states explain what happened and what to do next;
they do not apologize or sell. Keep register conversational, sentence case,
and matched to the audience. A label labels; an example demonstrates; nothing
quietly does two jobs.

## Compose The Interface

Design the reading order and interaction order together.

- Establish the primary action, the current state, and the most important
  information before styling secondary detail.
- Group related controls and information through proximity, alignment, shared
  containers, and repeated rhythm. Separate unrelated groups decisively.
- Use size, weight, contrast, placement, and whitespace to create hierarchy.
  Do not depend on font size alone.
- Keep semantic HTML correct while styling according to visual importance. The
  document outline and the visual hierarchy have different jobs.
- Give persistent or repeated information a stable position. Avoid layouts that
  jump when data, labels, or feedback change.
- Use cards sparingly: for repeated objects, a truly bounded tool, or a clear
  comparison unit. Do not put every page section in a card.
- Choose page width, grid, and density to fit the content. Empty space is useful
  only when it clarifies structure or focus.

For forms, put instructions next to the decision they support, make required
actions apparent, give errors a specific recovery path, and avoid forcing users
to remember information from another part of the page.

## Type, Color, And Depth

Use typography and color as functional systems, not isolated decoration.

- Choose typefaces and a compact scale for readability, tone, and the content
  density. Use a small number of text roles, then differentiate them by a
  purposeful combination of size, weight, line height, case, color, and space.
- Keep long-form text at a comfortable measure and line height. Break up dense
  copy with meaningful headings, lists, quotations, examples, or calls to
  action instead of shrinking the type until it fits.
- Treat labels, supporting copy, data, and headings as distinct roles. Make
  lower-priority information quieter without making it illegible.
- Build a palette by roles: neutral surfaces and text; an accent for primary
  interaction; and restrained semantic colors for success, warning, and error.
  Define tonal steps for each role before styling individual components.
- Implement both a dark color theme and a light color theme. Wire them with the
  CSS media query `prefers-color-scheme` (`@media (prefers-color-scheme: dark)`
  and `@media (prefers-color-scheme: light)`) so the interface follows the
  user's browser or OS preference. Do not ship a dark-only or light-only design
  unless the project's `DESIGN.md` already forbids one of those schemes. A
  product theme button is not required when `prefers-color-scheme` is the
  switch.
- Use the color system consistently: darker or stronger tones for emphasis and
  interaction, lighter tones for backgrounds and selected states, and neutral
  tones for the majority of the interface. A usable interface needs more than
  a page background, a brand color, and one gray.
- Keep contrast sufficient in every state, including muted text on tinted
  surfaces and text over imagery.
- Reinforce success, warning, error, selection, and status with text, icons,
  shape, or position as well as color.
- Use borders, shadows, background changes, and overlap to distinguish layers
  and interactive surfaces. Keep lighting and elevation consistent across the
  page.

## Design Text-Heavy Screens

For articles, documentation, guides, onboarding, and marketing pages, make the
reading experience the interface.

- Set a readable text column rather than allowing paragraphs to span the whole
  viewport. On wide screens, use the surrounding space for navigation, related
  content, or deliberate breathing room.
- Establish the headline, metadata, introduction, body, and calls to action as
  distinct text roles. Give each role a repeatable treatment rather than
  hand-tuning every block.
- Use headings to reveal the argument or sequence. Keep a heading visually and
  spatially attached to the content it introduces.
- Use emphasis sparingly. Bold, links, colored text, and highlighted panels all
  compete for attention; reserve each for a clear purpose.
- Treat quotes, lists, examples, author information, and sign-up prompts as
  meaningful interruptions in the reading rhythm, not interchangeable cards.
- Check the narrow-screen layout early. Reduce or relocate secondary content
  while protecting line length, hierarchy, and tap targets.

## Design Complex Forms

Organize long settings, account, checkout, and configuration forms around the
user's decisions, not around the underlying data model.

- Divide the form into named sections and make the section boundaries obvious.
  Put closely related fields together; do not make users infer the grouping from
  a long uninterrupted column.
- Keep the form at a readable working width. Use a second column only when the
  relationship remains clear and the screen has enough room.
- Pair a label, input, help text, and validation message as one visual unit.
  Use spacing and surface contrast to distinguish inputs from their container.
- Give selected plans, options, and payment methods a clear selected state that
  combines text, control state, and visual treatment. Do not depend on color
  alone.
- Place destructive actions away from the default save or continue path and
  label their consequence precisely. Make the normal next step visually
  dominant without making every button loud.
- Keep workflow actions in a stable location across related screens. When a
  narrow layout stacks controls, preserve their priority and their relationship
  to the section they affect.

## Design Data-Dense Screens

Dashboards, lists, reports, and operational tools should optimize recognition
and comparison, not imitate a marketing page.

- Decide what a user must understand at a glance, then dedicate the strongest
  hierarchy to those values, trends, alerts, or next actions.
- Separate summary, recent activity, detail, and navigation through layout and
  tonal structure. Avoid giving every metric, card, and table row equal visual
  weight.
- Align values so comparisons are easy: use consistent units, tabular number
  alignment where available, and right alignment for magnitudes when it helps
  scanning. Keep text labels readable and stable.
- Make the interactive target clear. A card may be entirely clickable, but its
  affordance and primary destination should still be evident.
- Use status color with text, icons, or badges, and reserve high-contrast
  treatment for conditions that require attention. Do not turn every value into
  a colored label.
- Use dividers and borders only where they improve row, column, or section
  separation. Repeated heavy boxes make dense information harder to scan.
- On narrow screens, retain the columns and controls needed for the immediate
  decision; move secondary detail into a disclosure, detail view, or deliberate
  horizontal-scroll region rather than silently losing it.

## Build Complete Interactions

Use familiar controls for familiar work, then make their state visible.

- Provide keyboard access, visible focus, labels, and appropriate semantics.
- Cover hover, focus, active, selected, disabled, loading, empty, error, and
  success states when they can occur.
- Make primary and destructive actions easy to distinguish before an action is
  taken.
- Use motion to explain cause, effect, or continuity. Respect reduced-motion
  preferences and avoid movement that distracts from repeated work.
- Design responsive behavior as a deliberate rearrangement of priority and
  controls, not just a desktop layout squeezed narrower.

## Implement In Context

Follow the project's existing framework, component library, CSS approach, and
design tokens when they are coherent. When they are not, improve the touched
surface without a speculative rewrite of unrelated areas.

- Build the working interface, not a presentation describing the intended UI.
- Use semantic structure and reusable variables or tokens for repeat values.
- Make constraints explicit with grid tracks, widths, aspect ratios, wrapping,
  truncation, and overflow handling where content demands them.
- Prefer the project’s established icon and asset strategy. Do not introduce
  assets or fonts whose licenses or loading behavior are unknown.
- Keep content, state, and interaction logic realistic enough to verify the
  visual result.

## Inspect With Browser-Use Or Computer-Use

Markup and CSS are not visual proof. When designing, restyling, or reviewing a
UI, open the running interface with the available **browser-use** or
**computer-use** tools and look at pixels, not only source.

- Run or attach to the local app when the project has a known start command.
- Capture the first viewport and the full page at a desktop width and a narrow
  phone width. Scroll through the whole surface you changed.
- Inspect **both** the dark theme and the light theme on the rendered page.
  Switch schemes with the browser or computer-use emulation of
  `prefers-color-scheme` (or the user's OS setting). Screenshot each scheme.
  Do not treat CSS or DOM inspection as proof that both themes look correct.
- Interact the way a user would: primary actions, hover and focus where it
  matters, empty or error states if they are in scope.
- Compare comparable live products in the same tools when judging whether the
  page sells, reads, or composes like its category.
- If the tool cannot connect (remote debugging, missing server), stop and get
  that unblocked rather than finishing from code alone.
- Fix what the screenshots show, then inspect again.

Reading the DOM or describing the CSS is supporting evidence. It does not
replace a screenshot of the rendered page.

## Review Before Handoff

After inspecting in the browser as above, correct the largest problems first,
then make a detail pass.

Check:

- Can a first-time user identify the screen's purpose, primary action, and
  current state quickly?
- Does the layout prioritize the workflow rather than the app chrome?
- Are grouping, alignment, spacing, and text hierarchy unambiguous?
- Do labels, long values, realistic data, and feedback fit without overlap or
  unstable layout shifts?
- Are all relevant interactive, loading, empty, error, and focus states clear?
- Does the visual direction fit the audience and product instead of resembling
  a generic template or one of the three AI-default clusters above?
- Is the palette a role-based system (surfaces, ink, one accent, semantics)
  rather than a leftover placeholder or a second competing accent?
- Do both the dark and light themes, driven by `prefers-color-scheme`, keep
  readable contrast, and were both visually inspected?
- Is there one signature element, and has unused type, unused tokens, and
  leftover decoration been removed?
- Is the result usable with keyboard navigation, reduced motion, and without
  relying on color alone?

In the final response, state what changed, which workflow or state was
prioritized, and what visual verification was performed.
