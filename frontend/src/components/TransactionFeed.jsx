import { DepositIcon, WithdrawIcon } from '../icons'

export default function TransactionFeed({ entries }) {
  if (entries.length === 0) {
    return <div className="empty">No transactions yet</div>
  }

  const isCredit = (op) => op === 'deposit' || op === 'transfer_in' || op === 'loan_disbursement'

  return (
    <div>
      {entries.slice().reverse().map((e) => (
        <div className="txn" key={e.id}>
          <div className="txn-icon">
            {isCredit(e.operation) ? <DepositIcon /> : <WithdrawIcon />}
          </div>
          <div className="txn-body">
            <div className="txn-op">{e.operation.replace('_', ' ')}</div>
            <div className="txn-meta">{new Date(e.created_at).toLocaleString()}</div>
          </div>
          <div className={`txn-amount ${isCredit(e.operation) ? 'in' : 'out'}`}>
            {isCredit(e.operation) ? '+' : '-'}{e.amount.toLocaleString()}
          </div>
        </div>
      ))}
    </div>
  )
}