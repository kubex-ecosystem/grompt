<!-- ---

title: Kubex Brand Visual Spec — v0.0.2 (Light‑First)
version: 0.0.2
owner: kubex
audience: dev|ops|stakeholder\languages: [en, pt-BR]
sources: [none]
assumptions: ["[ASSUMPTION] token names kept backward‑compatible; map provided below.", "[ASSUMPTION] asset pipeline exposes `generate:og` and `generate:thumbs` scripts."]
-----------  ----------------------------------------------------------------------------------------------------------------------------------------------------------------

## TL;DR

Light‑first aesthetic with dark as optional. Background **ice** `#f9fafb`, text **graphite** `#111827`, neon **cyan/lilac** accents preserved for icons/logos. Glows softened; contrast tuned for accessibility. This spec keeps Kubex DNA while modernizing enterprise appeal.

--- -->

# Kubex Brand Visual Spec — v0.0.2 (Light‑First)

> **Freedom, Engineered.** · **Independent by Design.**
> No artificial borders. Interoperability is our diplomacy.

## 1) Core Palette / Paleta Central

### EN

* **Base / Page**: Ice `#f9fafb`
* **Surface**: White `#ffffff`
* **Text / Graphite**: `#111827` (headings), `#334155` (body)
* **Primary (Cyan)**: `#06b6d4` (500), hover `#0891b2`, subtle `#ecfeff` (50)
* **Accent (Lilac/Fuchsia)**: `#a855f7` (500) / `#d946ef` (500)
* **Neon Glow (soft)**: cyan rgba(0,136,255,0.08) · lilac rgba(124,77,255,0.08)
* **Border**: `#e2e8f0` (slate‑200)
* **Success / Warn / Danger**: `#16a34a` / `#f59e0b` / `#ef4444`

### pt‑BR

* **Fundo / Página**: Gelo `#f9fafb`
* **Superfície**: Branco `#ffffff`
* **Texto / Grafite**: `#111827` (títulos), `#334155` (corpo)
* **Primária (Ciano)**: `#06b6d4` (500), hover `#0891b2`, sutil `#ecfeff` (50)
* **Acento (Lilás/Fúcsia)**: `#a855f7` (500) / `#d946ef` (500)
* **Brilho Neon (suave)**: ciano rgba(0,136,255,0.08) · lilás rgba(124,77,255,0.08)
* **Borda**: `#e2e8f0` (slate‑200)
* **Sucesso / Alerta / Erro**: `#16a34a` / `#f59e0b` / `#ef4444`

> **Note**: Keep **icon/logo hues** from v0.0.1 for backward compatibility. Use the new palette only for backgrounds, surfaces, text and UI accents.

---

## 2) Typography / Tipografia

* **Headings**: geometric/futuristic sans (Orbitron / Exo 2) — medium/semibold; tight tracking for H1.
* **Body**: system sans (Inter / UI) or IBM Plex Sans; comfortable line‑height (1.6).
* **Code/Monospace**: IBM Plex Mono / Source Code Pro.
* **Scale**: H1 40–64px, H2 28–40px, H3 22–28px, Body 16–18px.

### Do / Don’t

* ✅ High contrast on light backgrounds.
* 🚫 Avoid pure black text; prefer `#111827` to reduce glare.

---

## 3) Components / Componentes

### Badge (Light)

* Border `cyan-600/20`, BG `cyan-50`, Text `cyan-700`, radius **full**.
* Shadow **sm** only.

### Card / Section (Light)

* Surface `#ffffff`, Border `slate-200`, Shadow `sm`, Radius `2xl`.
* Prose (light): headings `slate-900`, body `slate-700`.

### Links

* Default text `slate-800`, underline decoration `cyan-600/40`, hover `cyan-700`.

### Glows / Hex Grid

* Reduce intensity to 8–10%; never exceed 15% on light.
* Hex stroke `rgba(0,76,153,0.05)`.

---

## 4) Accessibility / Acessibilidade

* **WCAG AA**: Body text ≥ 4.5:1, headings ≥ 3:1.
* Focus rings: `outline 2px` cyan `#06b6d4` offset `2px`.
* Motion: respect `prefers-reduced-motion`.
* Dark mode optional; maintain parity (token map below).

---

## 5) Token Map (v0.0.1 → v0.0.2)

| Old (Dark‑First)                     | New (Light‑First)                                         |
| ------------------------------------ | --------------------------------------------------------- |
| `bg.base = #0a0f14`                  | `bg.base = #f9fafb`                                       |
| `text.base = #e0f7fa`                | `text.head = #111827` · `text.body = #334155`             |
| `primary.cyan = #00f0ff`             | `primary.cyan = #06b6d4` (ui); **logos keep legacy neon** |
| `accent.lilac = #8e24aa`             | `accent.lilac = #a855f7`                                  |
| `border.muted = #0b1220`             | `border.muted = #e2e8f0`                                  |
| `glow.cyan = rgba(0,240,255,0.15)`   | `glow.cyan = rgba(0,136,255,0.08)`                        |
| `glow.lilac = rgba(124,77,255,0.25)` | `glow.lilac = rgba(124,77,255,0.08)`                      |

> Logos/Ícones: **não** alterar a matiz base; apenas ajustar **exposição** se necessário (‑5% a ‑10% de saturação em fundos muito claros).

---

## 6) Layout Rules / Regras de Layout

* Spacing scale 8px base (8/12/16/24/32/48/64).
* Grid **12 cols**, safe‑area padding 24–32px.
* Cards em colunas 2–3; uma coluna em mobile.
* Hero: título forte, subtítulo curto, uma cor de ênfase máx.
* Ícones em **hex** opcionais; linhas finas conectando.

---

## 7) Assets Kit

* **OG Cover**: 1200×630
* **YouTube Thumb**: 1280×720
* **Square**: 1080×1080
* Export PNG (sRGB) + SVG para logotipos.
* Badge “Powered by Kubex” quando couber.

---

## 8) Dark Mode Parity (Optional)

* Base `#0a0f14`, text `#e5f2f2`, borders `#0b1220`.
* Preserve the **same spacing and geometry**; only invert color tokens.

---

## 9) Migration Guide / Guia de Migração

1. Replace base tokens (`bg`, `surface`, `text`, `border`) conforme mapa.
2. Reduzir glows para ≤10%.
3. Atualizar componentes: `Badge`, `Card`, `Link` e `Section` (ver §3).
4. Validar contraste (AA) e foco.
5. Manter paleta de **logos/ícones** (apenas exposição se necessário).

***One‑command check / Um comando***

```bash
# [ASSUMPTION] pipeline de assets disponível
npm run check:brand && npm run generate:og && npm run generate:thumbs
```

---

## 10) Examples / Exemplos

### About Section (Light)

* See `AboutKubex.tsx` structure: hex grid suave, hero com “Freedom, Engineered.”, cards brancos.

### Badge

```html
<span class="rounded-full border border-cyan-600/20 bg-cyan-50 px-3 py-1 text-xs font-semibold tracking-wide text-cyan-700 shadow-sm">Open • Independent</span>
```

---

## Risks & Mitigations / Riscos & Mitigações

* **Logo washout** em fundos claros → ajustar exposição ‑5%/‑10%.
* **Contraste insuficiente** em lilás/fúcsia → usar cyan para estados críticos.
* **Inconsistência entre módulos** → script de lint de tokens + PR template reforçando §5.

---

## Next Steps / Próximos Passos

1. Publicar tokens v0.0.2 em cada módulo (TS/SCSS/JSON).
2. Rodar `check:brand` em todos os repositórios.
3. Atualizar OG images com capa light.
4. Adicionar dark‑mode opcional com paridade de tokens.
5. Revisar acessibilidade (AA) em páginas críticas.

---

### Changelog

* **0.0.2**: Light‑first, contraste revisado, mapeamento de tokens e guia de migração.
* **0.0.1**: Spec original dark‑first com neon tech.
