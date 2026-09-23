import { useState } from 'react'
import { api } from '../api/client'

export default function DepositWithdrawForm({ mode, account, onDone, busy, setBusy, flash, clearMsg }) {
  const [amount, setAmount] = useState('')
  const isDeposit = mode === 'deposit'

  const submit = async (e) => {
    e.preventDefault()
    const value = parseInt(amount, 10)
    if (!value || value <= 0) return flash('enter a valid amount')
    setBusy(true)
    clearMsg()
    try {
      if (isDeposit) await api.deposit(account.id, value)
      else await api.withdraw(account.id, value)
      flash(`${isDeposit ? 'Deposited' : 'Withdrew'} ${value} KES`, 'success')
      setAmount('')
      onDone()
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <form onSubmit={submit}>
      <h3>{isDeposit ? 'Deposit funds' : 'Withdraw funds'}</h3>
      <label>Amount (KES)</label>
      <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="0" />
      <button disabled={busy}>{busy ? 'Processing...' : isDeposit ? 'Deposit' : 'Withdraw'}</button>
    </form>
  )
}