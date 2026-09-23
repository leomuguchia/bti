import { useState } from 'react'
import { api } from '../api/client'

export default function TransferForm({ accounts, fromAccount, onDone, busy, setBusy, flash, clearMsg }) {
  const [toId, setToId] = useState('')
  const [amount, setAmount] = useState('')
  const others = accounts.filter((a) => a.id !== fromAccount.id)

  const submit = async (e) => {
    e.preventDefault()
    const value = parseInt(amount, 10)
    if (!toId) return flash('choose a destination account')
    if (!value || value <= 0) return flash('enter a valid amount')
    setBusy(true)
    clearMsg()
    try {
      await api.transfer(fromAccount.id, toId, value)
      flash(`Sent ${value} KES`, 'success')
      setAmount('')
      onDone()
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  if (others.length === 0) {
    return <p className="muted-note" style={{ marginTop: 0 }}>Create a second account to enable transfers.</p>
  }

  return (
    <form onSubmit={submit}>
      <h3>Transfer funds</h3>
      <label>To account</label>
      <select value={toId} onChange={(e) => setToId(e.target.value)}>
        <option value="">Select account</option>
        {others.map((a) => <option key={a.id} value={a.id}>{a.owner.full_name}</option>)}
      </select>
      <label>Amount (KES)</label>
      <input type="number" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="0" />
      <button disabled={busy}>{busy ? 'Sending...' : 'Transfer'}</button>
    </form>
  )
}