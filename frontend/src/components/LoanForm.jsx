import { useState } from 'react'
import { api } from '../api/client'

export default function LoanForm({ account, loan, setLoan, onDone, busy, setBusy, flash, clearMsg }) {
  const [amount, setAmount] = useState('')

  const create = async () => {
    setBusy(true)
    clearMsg()
    try {
      const l = await api.createLoan(account.id)
      setLoan(l)
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  const disburse = async (e) => {
    e.preventDefault()
    const value = parseInt(amount, 10)
    if (!value || value <= 0) return flash('enter a valid amount')
    setBusy(true)
    clearMsg()
    try {
      await api.disburseLoan(loan.id, value)
      flash(`Disbursed ${value} KES`, 'success')
      setAmount('')
      onDone()
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  if (!loan) {
    return (
      <div>
        <h3>Loan account</h3>
        <button disabled={busy} onClick={create}>{busy ? 'Creating...' : 'Create loan account'}</button>
      </div>
    )
  }

  return (
    <form onSubmit={disburse}>
      <h3>Disburse loan</h3>
      <p className="muted-note" style={{ marginTop: 0 }}>Loan status: {loan.status}</p>
      <label>Amount (KES)</label>
      <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="10000" />
      <button disabled={busy}>{busy ? 'Disbursing...' : 'Disburse'}</button>
    </form>
  )
}