---
name: frontend-design
description: >
  Create, improve, or review production frontend interfaces: web pages, apps,
  dashboards, tools, components, games, posters, artifacts, landing pages, and
  HTML/CSS/JS or framework UI. Use whenever the user asks for frontend design,
  styling, layout, visual polish, beautification, responsive UI, or a web
  experience. Produces distinctive, context-appropriate interfaces with strong
  visual craft instead of generic AI aesthetics.
---

# Frontend Design

Use this skill to design and build frontend interfaces that feel intentionally
designed for their product, audience, and workflow. The goal is not decoration;
the goal is a coherent, memorable, usable interface.

## Design Direction

- Understand the purpose, audience, and expected workflow before choosing the
  visual direction.
- Choose a clear aesthetic point of view. It may be brutally minimal,
  maximalist, editorial, playful, industrial, luxurious, utilitarian, retro,
  organic, technical, or something more specific to the domain.
- Match intensity to the product category. A poster, portfolio, game, or
  landing page can be expressive; a dashboard, editor, admin surface, or
  repeated-use operational tool should usually be quieter, denser, and easier to
  scan.
- Avoid generic AI aesthetics: default card grids, purple-blue gradients on
  white, random glowing blobs, decorative glassmorphism, stock SaaS layouts, and
  unmotivated grain/noise.
- Make the interface memorable through concept, hierarchy, composition,
  typography, assets, interaction, or domain-specific details, not through
  decoration pasted on top.

## Product Fit

- Design for the real first screen, not a marketing page unless the user asked
  for one.
- For apps and tools, prioritize the actual workflow: navigation, controls,
  state, data density, filtering, comparison, repeated actions, and feedback.
- For marketing, brand, venue, portfolio, or object-focused pages, make the
  subject unmistakable in the first viewport.
- For games and interactive toys, make the playable or interactive surface the
  primary experience.
- Use visual assets when they help communicate the product, place, object,
  gameplay, person, or state. Prefer real or generated bitmap imagery over
  abstract SVG decoration when the user needs to inspect the subject.

## Practical Visual Refinement

Use a Refactoring UI style approach: improve interfaces through concrete visual
decisions.

- Design with realistic content and states, not empty placeholder boxes.
- Establish visual hierarchy with size, weight, color, spacing, contrast, and
  position. Semantic heading level does not automatically decide visual size.
- Make primary actions and important information visually obvious.
- Use contrast intentionally. Muted elements recede; important elements stand
  out.
- Use spacing as a design tool. Related items sit closer together; unrelated
  groups get more separation.
- Use a small set of strong, repeatable tokens for color, type, radius, shadow,
  spacing, borders, and motion.
- Add polish through precise details: alignment, icon sizing, border contrast,
  focus states, hover states, empty states, loading states, shadows, and
  responsive behavior.
- Do not treat decoration as design. Gradients, grain, glows, textures, custom
  cursors, and motion only help when they reinforce the product's tone and
  hierarchy.

## Typography And Color

- Avoid lazy defaults. Choose typography for readability, tone, performance, and
  brand fit.
- Distinctive display fonts can elevate expressive work; system fonts can be
  correct for dense tools and native-feeling apps.
- Pair type deliberately. Avoid using many font families or many display styles.
- Keep palettes intentional. Use enough contrast for hierarchy and readability.
- Avoid one-note palettes dominated by a single hue unless the concept demands
  it.
- Use color to convey meaning only when shape, text, or position also supports
  the state.

## Layout And Composition

- Create layouts that fit the workflow: dense and scannable for tools, more
  editorial or immersive for content and brand experiences.
- Use grids, alignment, and rhythm deliberately, then break the grid only when
  it improves attention or meaning.
- Keep fixed-format UI stable with explicit dimensions, aspect ratios, or grid
  tracks so labels, icons, hover states, and dynamic content do not shift the
  layout.
- Ensure text fits within its containers at desktop and mobile sizes. Long words
  and button labels must not overflow or overlap.
- Do not put cards inside cards. Use cards for repeated items, modals, or
  genuinely framed tools, not as default page-section wrappers.

## Interaction States

- Build complete controls and states a user would expect: hover, focus, active,
  disabled, loading, empty, error, success, selected, and expanded/collapsed.
- Use familiar controls for familiar jobs: icon buttons for tools, segmented
  controls for modes, toggles for binary settings, sliders/inputs for numeric
  values, menus for option sets, and tabs for views.
- Use icons where they improve scanning. Prefer the project's icon library when
  one exists.
- Motion should support comprehension, feedback, continuity, or delight. Avoid
  surprise motion in repeated-use tools.
- Respect reduced-motion preferences when adding substantial animation.

## Implementation Guidance

- Follow the existing app framework, component system, CSS strategy, and design
  conventions when they are coherent.
- If the existing UI is weak or inconsistent, preserve necessary behavior while
  improving the touched surface with the principles in this skill.
- Build the actual usable experience, not a feature list explaining the
  experience.
- Do not add visible instructional text about design decisions, keyboard
  shortcuts, or implementation unless the product itself needs it.
- Use CSS variables or design tokens for repeated values.
- Keep accessibility practical: semantic structure, keyboard reachability,
  visible focus states, readable contrast, labels, and meaningful alt text.

## Review Pass

Before finishing, review the interface like a designer and a user:

- Does the first screen make the product/task clear?
- Does the visual direction fit the domain and audience?
- Is hierarchy obvious without reading every word?
- Are primary actions and current state clear?
- Are spacing, alignment, type scale, border contrast, and icon sizes polished?
- Does text fit without overlap on mobile and desktop?
- Are loading, empty, disabled, focus, hover, and error states handled where
  relevant?
- Does motion support the experience instead of distracting from it?
- Are decorative effects justified by the concept?
- Does the final result avoid generic AI UI patterns?

## Verification

- When a dev server is needed, run it and provide the local URL.
- When a static HTML file is enough, provide the file path.
- For substantial frontend changes, verify visually at desktop and mobile
  viewport sizes. Use screenshots or browser checks when available.
- For canvas, WebGL, Three.js, generated media, or asset-heavy work, verify that
  the primary visual is nonblank, framed correctly, and assets load.
