const base = { fill: 'none', strokeWidth: 1.8, strokeLinecap: 'round', strokeLinejoin: 'round' }

export const DepositIcon = () => (
  <svg viewBox="0 0 24 24" {...base}><path d="M12 4v13m0 0l-5-5m5 5l5-5M4 20h16" /></svg>
)
export const WithdrawIcon = () => (
  <svg viewBox="0 0 24 24" {...base}><path d="M12 20V7m0 0l-5 5m5-5l5 5M4 4h16" /></svg>
)
export const TransferIcon = () => (
  <svg viewBox="0 0 24 24" {...base}><path d="M7 7h13m0 0l-4-4m4 4l-4 4M17 17H4m0 0l4 4m-4-4l4-4" /></svg>
)
export const LoanIcon = () => (
  <svg viewBox="0 0 24 24" {...base}><path d="M3 10l9-6 9 6M5 10v9m14-9v9M3 19h18M9 19v-6h6v6" /></svg>
)
export const AddIcon = () => (
  <svg viewBox="0 0 24 24" {...base}><path d="M12 5v14M5 12h14" /></svg>
)