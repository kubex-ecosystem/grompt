export interface Example {
  purpose: string;
  ideas: string[];
}

export const examples: Example[] = [
  {
    purpose: "Code Generation",
    ideas: [
      "Create a React hook for fetching data from an API.",
      "It should handle loading, error, and data states.",
      "Use the native `fetch` API.",
      "The hook should be written in TypeScript and be well-documented.",
    ],
  },
  {
    purpose: "Creative Writing",
    ideas: [
      "Write a short story opening.",
      "The setting is a neon-lit cyberpunk city in 2077.",
      "The main character is a grizzled detective who is part-cyborg.",
      "It's perpetually raining and the streets are reflective.",
    ],
  },
  {
    purpose: "Data Analysis",
    ideas: [
      "Analyze a dataset of customer sales from the last quarter.",
      "The dataset includes columns: 'Date', 'CustomerID', 'ProductCategory', 'Revenue', 'UnitsSold'.",
      "Identify the top 3 product categories by total revenue.",
      "Calculate the average revenue per customer.",
      "Look for any weekly sales trends or seasonality.",
    ],
  },
  {
    purpose: "Marketing Copy",
    ideas: [
      "Draft an email campaign for a new productivity app.",
      "The target audience is busy professionals and university students.",
      "Highlight features like AI-powered task scheduling, calendar sync, and focus mode.",
      "The tone should be encouraging, professional, and slightly urgent.",
    ],
  },
  {
    purpose: "Technical Documentation",
    ideas: [
        "Write the 'Getting Started' section for a new JavaScript library.",
        "The library is called 'ChronoWarp' and it simplifies date manipulation.",
        "Include a simple installation guide using npm.",
        "Provide a clear, concise code example for its primary use case."
    ]
  }
];

export const purposeKeys: Record<string, string> = {
  "Code Generation": "purposeCodeGeneration",
  "Creative Writing": "purposeCreativeWriting",
  "Data Analysis": "purposeDataAnalysis",
  "Technical Documentation": "purposeTechnicalDocumentation",
  "Marketing Copy": "purposeMarketingCopy",
  "General Summarization": "purposeGeneralSummarization",
};