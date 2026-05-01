---
name: My Store
description: A calm, pastel inventory management interface built for portfolio craft.
colors:
  sage: "oklch(72% 0.09 155)"
  sage-deep: "oklch(65% 0.11 155)"
  pale-sage: "oklch(94% 0.015 155)"
  parchment: "oklch(99% 0.004 95)"
  canvas: "oklch(96% 0.008 95)"
  chalk: "oklch(87% 0.010 95)"
  stone: "oklch(82% 0.012 95)"
  ink: "oklch(22% 0.008 95)"
  ash: "oklch(57% 0.010 95)"
  clay: "oklch(60% 0.16 25)"
typography:
  title:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: "18px"
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: "normal"
  body:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.6
    letterSpacing: "normal"
  label:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: "13px"
    fontWeight: 600
    lineHeight: 1.4
    letterSpacing: "0.04em"
  caption:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: "13px"
    fontWeight: 400
    lineHeight: 1.4
    letterSpacing: "normal"
rounded:
  sm: "6px"
  md: "8px"
  lg: "10px"
spacing:
  xs: "4px"
  sm: "8px"
  md: "12px"
  lg: "20px"
  xl: "32px"
components:
  button-primary:
    backgroundColor: "{colors.sage}"
    textColor: "{colors.parchment}"
    rounded: "{rounded.md}"
    padding: "10px 18px"
  button-primary-hover:
    backgroundColor: "{colors.sage-deep}"
    textColor: "{colors.parchment}"
    rounded: "{rounded.md}"
    padding: "10px 18px"
  button-secondary:
    backgroundColor: "{colors.chalk}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "10px 18px"
  button-secondary-hover:
    backgroundColor: "{colors.stone}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "10px 18px"
  button-danger:
    backgroundColor: "{colors.clay}"
    textColor: "{colors.parchment}"
    rounded: "{rounded.md}"
    padding: "10px 18px"
  card:
    backgroundColor: "{colors.parchment}"
    rounded: "{rounded.md}"
    padding: "{spacing.md}"
  input:
    backgroundColor: "{colors.parchment}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "8px 12px"
  input-focus:
    backgroundColor: "{colors.parchment}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "8px 12px"
---

# Design System: My Store

## 1. Overview

**Creative North Star: "The Clear Ledger"**

A ledger is not decorative. It is precise, structured, and quietly confident — every entry in its place, every column legible at a glance. This interface carries that character: pastel surfaces that organize without competing, typography that delivers information without ceremony, and space used with intention. Nothing here is accidental.

The palette is warm and soft — sage greens, warm whites, quiet ash — but never precious. Color is working structure: it marks section boundaries, signals interaction, and guides the eye from section to action. The goal is that a reviewer pauses and sees considered decisions, not assembled defaults. Craft within constraints: no extra libraries, no downloaded fonts. All quality comes from OKLCH color, spacing rhythm, and type hierarchy.

This system rejects the generic purple-gradient SaaS aesthetic entirely. No dark-mode-by-default, no glowing cards, no gradient-filled buttons, no glassmorphism, no decorative blurs. Calm and efficient — the interface earns attention by staying out of the way.

**Key Characteristics:**
- Single pastel sage accent on a warm off-white ground
- Generous whitespace; sections breathe before they stack
- Uppercase tracked labels as structural dividers — no extra borders needed
- Flat cards at rest; elevation earned only through hover state
- Every neutral carries a fractional warm tint — the warmth is invisible as a fact and perceptible as a feeling

## 2. Colors: The Warm Ledger Palette

One accent, warm-tinted neutrals throughout, and a muted clay for errors. Pastels carry hierarchy as surface tints, not as decoration.

### Primary
- **Sage Moss** (`oklch(72% 0.09 155)`): The single accent. Primary buttons, active nav states, focus rings, and interactive affordances. Restricted to at most 15% of any given surface — scarcity is intentional.
- **Deep Moss** (`oklch(65% 0.11 155)`): Hover and pressed state of Sage Moss only. Never used as a standalone background.

### Tertiary
- **Pale Sage** (`oklch(94% 0.015 155)`): Ambient surface tint for category section headers, empty-state backgrounds, and any surface that needs gentle hierarchy without a border. The organized-drawer tone.

### Neutral
- **Parchment White** (`oklch(99% 0.004 95)`): All card and form backgrounds. Every elevated surface.
- **Warm Canvas** (`oklch(96% 0.008 95)`): Page background. The ground the interface sits on.
- **Warm Chalk** (`oklch(87% 0.010 95)`): Borders at rest; secondary button fill.
- **Warm Stone** (`oklch(82% 0.012 95)`): Borders on hover; dividers.
- **Warm Ink** (`oklch(22% 0.008 95)`): All primary body text and headings. Never pure black.
- **Quiet Ash** (`oklch(57% 0.010 95)`): Muted text — labels, captions, secondary metadata, placeholder text.
- **Muted Clay** (`oklch(60% 0.16 25)`): Error states only. Muted rather than alarming — consistent with the calm register.

### Named Rules
**The One Accent Rule.** Sage Moss appears on one interaction type per screen at a time — either buttons, or active states, or focus rings — never all three at full saturation simultaneously. Pale Sage handles ambient hierarchy. Sage Moss marks action.

**The Warm Tint Rule.** Every neutral carries a fractional warm tint (hue 95, chroma 0.004–0.012). Pure white (`oklch(100% 0 0)`) and pure black (`oklch(0% 0 0)`) are prohibited. The warmth is nearly invisible as a value and perceptible as a character.

## 3. Typography

**Body/UI Font:** System font stack — `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`

**Character:** No custom fonts; all hierarchy comes from weight, size, tracking, and spacing rhythm. The system stack reads cleanly on every OS without a network request, consistent with the no-extra-dependencies constraint. The hierarchy relies on scale contrast (title 18px to body 15px to label 13px) amplified by weight and uppercase tracking on labels.

### Hierarchy
- **Title** (600 weight, 18px, 1.3 line-height): Page headings and section headings. The top of the visible hierarchy on any screen.
- **Body** (400 weight, 15px, 1.6 line-height): All prose, metadata, descriptions, and values. Max line length 65ch on reading surfaces.
- **Label** (600 weight, 13px, 0.04em letter-spacing, uppercase): Form field labels, column headers, category section headings in ProductListPage. Creates structural division without borders.
- **Caption** (400 weight, 13px, 1.4 line-height): Supporting metadata — stock counts, secondary data, muted helper text. Uses Quiet Ash color.

### Named Rules
**The Label-as-Divider Rule.** Uppercase tracked labels (13px, 600, 0.04em) replace visual separators. Before adding a horizontal rule, border, or extra container, ask whether a Label-style heading can carry the structural weight instead.

## 4. Elevation

Flat by default. Depth is earned by interaction state, not applied to signal importance.

Cards and containers have no shadow at rest — they read as surfaces through background color difference (Parchment White on Warm Canvas). On hover and in active states, a warm-tinted shadow appears to signal responsiveness. This keeps the resting interface calm and gives hover states genuine meaning.

### Shadow Vocabulary
- **Flat** (no shadow): All cards, panels, and containers at rest. Tonal contrast from Parchment on Canvas is sufficient.
- **Hover Lift** (`0 4px 16px oklch(22% 0.008 95 / 0.10)`): Clickable cards on pointer hover. Soft, warm-tinted, directional downward shadow.
- **Elevated** (`0 8px 24px oklch(22% 0.008 95 / 0.14)`): Open dropdown menus, active popovers. Slightly stronger; same warm family.

### Named Rules
**The Earned Elevation Rule.** Shadows are a response to state, not a badge of importance. A static, non-interactive surface has no shadow. Period.

## 5. Components

### Buttons
Soft and tactile — gently rounded, clear weight contrast between primary and secondary. The primary button is the only fully saturated element on most screens.

- **Shape:** Gently curved (8px radius)
- **Primary:** Sage Moss fill, Parchment White text, 10px top/bottom 18px left/right padding. The single call-to-action.
- **Hover:** Deepens to Deep Moss; lifts 1px via `transform: translateY(-1px)`. No box-shadow added on hover — the color shift carries the state.
- **Focus:** 2px Sage Moss outline, 2px offset. Visible but not alarming.
- **Secondary:** Warm Chalk fill, Warm Ink text. Recedes behind primary. Used for Cancel and neutral actions.
- **Danger:** Muted Clay fill, Parchment White text. Muted, not alarming.
- **Disabled:** 50% opacity, `cursor: not-allowed`, no transform. Quiet, not hidden.

### Cards / Containers
The primary grouping unit. Flat at rest; responsive on hover.

- **Corner Style:** Gently curved (8px radius)
- **Background:** Parchment White
- **Shadow:** None at rest; Hover Lift (`0 4px 16px oklch(22% 0.008 95 / 0.10)`) on interactive cards
- **Border:** 1px Warm Chalk at rest; transitions to 1px Warm Stone on hover
- **Padding:** `md` (12px) interior

### Inputs / Fields
Legible and calm. Labels always above the field, uppercase tracked. No floating labels.

- **Style:** 1px Warm Chalk border, Parchment White fill, 6px radius, 8px/12px padding
- **Focus:** Border shifts to Sage Moss at 2px width. No glow or shadow — the color shift is sufficient.
- **Error:** Border shifts to Muted Clay; Muted Clay error text below at Caption size.
- **Disabled:** Warm Canvas fill, Quiet Ash text, `cursor: not-allowed`.
- **Textarea:** Same treatment; `resize: vertical` permitted.

### Navigation
Dark and grounding — the one surface that uses deep ink to anchor the page against the light ground.

- **Background:** Warm Ink (`oklch(22% 0.008 95)`)
- **Logo / Links:** Parchment White, 500 weight
- **Link Hover:** Underline + 85% opacity. No background highlight.
- **Layout:** `justify-content: space-between`, `padding: 15px 30px`. Simple and flat.

### Category Section Headers (Signature Component)
Category groupings on ProductListPage use uppercase Label typography on a Pale Sage background strip. This creates visual rhythm and section identity without added cards or heavy borders — the organization is legible before a single product renders.

- **Background:** Pale Sage (`oklch(94% 0.015 155)`)
- **Typography:** Label style — 13px, 600 weight, 0.04em tracking, uppercase, Warm Ink color
- **Border-bottom:** 2px solid Warm Chalk — a single grounding edge below the heading
- **Padding:** `sm` vertical, `lg` horizontal — consistent with PageShell padding

## 6. Do's and Don'ts

### Do:
- **Do** use Sage Moss as the only saturated accent. One color carries action — its restraint is the point.
- **Do** tint every neutral toward hue 95, chroma 0.004–0.012. Pure white and pure black are banned.
- **Do** use uppercase tracked labels (13px, 600, 0.04em) as structural dividers before reaching for borders or extra containers.
- **Do** keep all cards flat at rest. Shadows appear on hover only — they signal responsiveness, not rank.
- **Do** use Pale Sage tints for ambient hierarchy — category sections, empty states, raised panels — before using shadows or borders.
- **Do** keep page layouts open: one `PageShell` at `max-width: 1200px` with consistent `lg` (20px) padding. Sections breathe before they stack.
- **Do** use `transform` and `opacity` for all animation. Never animate layout properties.

### Don't:
- **Don't** use purple, violet, indigo, or any blue-shifted hue as an accent. The purple-gradient SaaS aesthetic ("AI slop") is explicitly banned by this project.
- **Don't** use gradient backgrounds, gradient-filled buttons, or `background-clip: text` gradient text. Prohibited.
- **Don't** use `border-left` or `border-right` greater than 1px as a colored accent stripe on cards, list items, or callouts. Rewrite with a background tint, a label, or nothing.
- **Don't** default to dark mode because "tools look dark." Warm Ink is used only for the Navbar — one anchoring surface, not a theme.
- **Don't** install additional CSS libraries, component frameworks, or icon packages. All visual quality comes from CSS custom properties, typography, and spacing.
- **Don't** add emoji, decorative icons, or illustrative elements to the UI. Structure comes from layout and type.
- **Don't** apply a shadow to a non-interactive surface. A static card that cannot be clicked has no shadow.
- **Don't** stack cards inside cards. Nested cards are always wrong.
- **Don't** use Bootstrap density — tight gutters, small form elements, crowded grids. The design principle is calm efficiency, which requires space.
