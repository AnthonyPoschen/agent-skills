---
name: eyes-on
description: >
  Look at a rendered change the way a person would, before calling it done.
  Use whenever you finish, review, or debug something someone can see or
  operate: a page, layout, style, component, dashboard, homepage, email, or
  other interface. Also use when the user says verify, looks good, check the
  browser, screenshot, or that they had to review it themselves. Exercise the
  feature, save screenshots with scripts/see.mjs, and read those image files.
  A DOM query, a CSS diff, or "the element exists" is not this look. Skip it
  for a change with no visible or operable surface.
---

# Eyes on

A change a person can see is done when you have looked at the rendered result
and the feature behaves. Source, a passing test, and a DOM query are support.
They do not show whether the piece fits the page.

## Loop

1. Name the action a person performs and the result they should see.
2. Run the project's own app. A repo skill or `AGENTS.md` owns the command,
   the URL, and where logs are. Use that. Do not invent a second stack.
3. Perform the action on the running surface. Click, type, and submit the
   way that person would. Loading the URL is the action only when the change
   is the resting page. Say which one you did.
4. Save pictures of that result. For a URL, run `scripts/see.mjs` in this
   skill. Open every PNG it prints with the image reader. A screenshot tool
   that returns an image counts only after you judge that image in the eyes
   note. For a terminal or native window, save a PNG from the computer-use
   tool and open that file the same way.
5. Write the eyes note from the pictures, then fix anything that fails the
   bar below. Shoot and look again after a visual fix.
6. In the handoff, name the action, the pictures, and any problem you left.

Stop before the pictures when the change has no surface a person sees or
operates. Say that, and use the project's ordinary checks.

If the app or the browser will not open, say so. Do not finish from source.

## Capture

```sh
node ~/.agents/skills/eyes-on/scripts/see.mjs \
  --url 'http://127.0.0.1:PORT/path' \
  --out /tmp/eyes-on
```

The script needs Playwright (`EYES_ON_PLAYWRIGHT` is the package root when it
is not on `NODE_PATH`) and Chromium on `PATH` (`EYES_ON_CHROMIUM` overrides).

Defaults are a desktop viewport (`1280x800`) and a phone viewport
(`390x844`), in light and dark. Add `--wait` or `--scroll` with a CSS
selector so the shot contains the change and its neighbors. Add `--full`
only when the change is the page's vertical composition. A long image hides
the detail you need to judge.

Read `eyes.txt` in the output directory. It lists the PNGs. Look at each one.

## What you are judging

Answer these from the pictures, in concrete visual terms. "It looks fine" and
"the node is present" both fail the note.

- **Behavior.** The action changed the picture in the way the feature claims.
  Compare the before and after when the change is a response to input.
- **Fit.** Type, radius, border, elevation, and color roles match the
  neighbors and the project's design record. A one-off treatment looks pasted
  on, even when the new element is the one you added.
- **Alignment.** Edges and baselines follow the surrounding rhythm. Text is
  not clipped, overlapping, or colliding with a control. Spacing is the
  page's rhythm, not a new gap invented for this block.
- **Hierarchy.** The picture shows the primary action and the current state
  without the source.
- **Width.** The phone picture still lets the person complete the action.
  The desktop picture is not the phone layout stretched.
- **Scheme.** When the surface follows the system theme, light and dark both
  keep text readable and surfaces distinct.

The project's design record wins on palette and character. This skill does
not invent a new visual direction.

## Eyes note

Write this before you claim the change is done:

```text
Eyes
- Action: what you did on the running surface
- Pictures: the paths you opened
- Saw: fit, alignment, hierarchy, and behavior, from the pixels
- Widths and schemes: what held and what broke
- Verdict: pass, or the fix you are making now
```

Pass names a specific thing you saw. This fails the note: "the card exists
in the DOM." This passes: "the new price uses the same panel radius and
hairline as the cards above it; the amount lines up with the title; at
390px the button wraps into the footnote."

Fix a miss, then replace the pictures and the note. Leave a miss in the
handoff only when you cannot fix it in this change, and say what you saw.
