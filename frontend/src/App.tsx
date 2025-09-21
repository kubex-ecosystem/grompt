/**
 * Grompt - AI Prompt Engineering Tool
 * Main App component with integrated backend API and Analyzer design system
 */

import * as React from 'react'
import PromptCrafter from './PromptCrafter'

// Analyzer design system background effects
const BackgroundEffects: React.FC = () => (
  <>
    {/* Grid pattern background - following Analyzer design */}
    <div
      className="fixed inset-0 z-[-10]"
      style={{
        backgroundImage: `
          linear-gradient(rgba(128,128,128,0.1) 1px, transparent 1px),
          linear-gradient(to right, rgba(128,128,128,0.1) 1px, transparent 1px)
        `,
        backgroundSize: '24px 24px'
      }}
    />

    {/* Vignette effect for depth - following Analyzer design */}
    <div
      className="fixed inset-0 z-[-9]"
      style={{
        backgroundImage: 'radial-gradient(circle at center, transparent 40%, #030712 90%)'
      }}
    />
  </>
)

export default function App() {
  return (
    <div className="min-h-screen bg-[#030712] text-white font-sans selection:bg-purple-500/30">
      <BackgroundEffects />
      <PromptCrafter />
    </div>
  )
}
