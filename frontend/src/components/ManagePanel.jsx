import { api } from '../api/client'

export default function ManagePanel({ account, onDeleted, busy, setBusy, flash, clearMsg }) {
  const handleDelete = async () => {
    setBusy(true)
    clearMsg()
    try {
      await api.deleteAccount(account.id)
      onDeleted()
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <div>
      <h3>Manage account</h3>
      <p className="muted-note" style={{ marginTop: 0 }}>
        Deletion only succeeds when the balance is zero and the account has
        had no activity for the dormancy period.
      </p>
      <button className="danger" disabled={busy} onClick={handleDelete}>
        {busy ? 'Checking...' : 'Delete account'}
      </button>
    </div>
  )
}