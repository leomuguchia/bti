import { DepositIcon, WithdrawIcon, TransferIcon, LoanIcon } from '../icons'

const ACTIONS = [
  { key: 'deposit', label: 'Deposit', Icon: DepositIcon },
  { key: 'withdraw', label: 'Withdraw', Icon: WithdrawIcon },
  { key: 'transfer', label: 'Transfer', Icon: TransferIcon },
  { key: 'loan', label: 'Loan', Icon: LoanIcon },
]

export default function ActionGrid({ activePanel, onToggle }) {
  return (
    <div className="grid">
      {ACTIONS.map(({ key, label, Icon }) => (
        <button
          key={key}
          className={`action-btn ${activePanel === key ? 'active' : ''}`}
          onClick={() => onToggle(key)}
        >
          <Icon />
          <span>{label}</span>
        </button>
      ))}
    </div>
  )
}