# Assets Checklist for Kubex.world

This document lists all the assets needed for the complete Grompt v2.0 deployment.

## Required Icons & Images

### ✅ Already Referenced (Need to Create/Update)

```
/assets/
├── favicon.ico                     # Browser favicon
├── apple-touch-icon.png           # iOS home screen icon (180x180)
├── safari-pinned-tab.svg          # Safari pinned tab icon
├── kubex_og.png                   # Open Graph/Social media image
└── icons/
    ├── favicon.svg                # SVG favicon
    ├── android-chrome-192x192.png # Android Chrome icon
    └── android-chrome-512x512.png # Android Chrome icon
```

## Icon Specifications

### Favicon Package

- **favicon.ico**: 16x16, 32x32, 48x48 (multi-size ICO)
- **favicon.svg**: Vector favicon (dark/light mode compatible)

### Mobile Icons

- **apple-touch-icon.png**: 180x180 PNG for iOS
- **android-chrome-192x192.png**: 192x192 PNG for Android
- **android-chrome-512x512.png**: 512x512 PNG for Android

### Safari

- **safari-pinned-tab.svg**: Monochrome SVG, color: #00f0ff

### Social Media

- **kubex_og.png**: 1200x630 PNG for Open Graph/Twitter Cards

## Design Guidelines

### Colors

- **Primary**: #00f0ff (Kubex cyan)
- **Secondary**: #00e676 (Kubex green)
- **Dark Background**: #0a0f14
- **Light Background**: #f9fafb

### Style

- **Minimalist**: Clean, geometric design
- **Hexagonal Theme**: Following Kubex hexagon logo pattern
- **High Contrast**: Works on both light and dark backgrounds
- **Scalable**: Vector-based when possible

## Logo Concept

Based on the Kubex logo (hexagon with circuit-like patterns), the Grompt icon should:

1. **Maintain Kubex Identity**: Use the hexagonal shape
2. **Add Prompt Context**: Include elements suggesting text/AI (maybe stylized "G" or chat bubble)
3. **Circuit Aesthetics**: Keep the tech/circuit theme
4. **Glow Effects**: Neon glow for dark mode

## Current Status

### ❌ Missing Files

These files are referenced but don't exist yet:

- [ ] `/assets/favicon.ico`
- [ ] `/assets/apple-touch-icon.png`
- [ ] `/assets/safari-pinned-tab.svg`
- [ ] `/assets/kubex_og.png`
- [ ] `/assets/icons/favicon.svg`
- [ ] `/assets/icons/android-chrome-192x192.png`
- [ ] `/assets/icons/android-chrome-512x512.png`

### 📝 Creation Tools Suggested

- **Vector Design**: Figma, Adobe Illustrator, or Inkscape
- **PNG Export**: Export from vector at proper resolutions
- **ICO Creation**: Use online ICO generator or ImageMagick
- **Optimization**: TinyPNG for file size optimization

### 🎨 Design Priority

1. **Create base SVG logo first** (favicon.svg)
2. **Export PNG versions** at required sizes
3. **Generate ICO file** from PNG versions
4. **Create social media banner** (kubex_og.png)

### 🚀 Quick Win

For immediate deployment, you could:

1. Use a simple geometric "G" in a hexagon
2. Apply Kubex color scheme (#00f0ff)
3. Export at all required sizes
4. Test on different devices/browsers

## SEO Impact

Having proper icons improves:

- **User Trust**: Professional appearance
- **Brand Recognition**: Consistent visual identity
- **Mobile Experience**: Proper home screen icons
- **Social Sharing**: Rich previews with custom images

---

**Note**: These assets are crucial for the professional launch at kubex.world. The current deployment will work but may show default browser icons until these are created.
