declare global {
  interface Window {
    turnstile?: {
      render: (selector: string, options: Record<string, unknown>) => string
      remove: (widgetId: string) => void
      reset: (widgetId: string | number) => void
    }
  }
}

export {}
