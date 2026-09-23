import { useState, useEffect } from 'react'
import { api } from './api/client'
import TopBar from './components/TopBar'
import HeroBalance from './components/HeroBalance'
import ActionGrid from './components/ActionGrid'
import Panel from './components/Panel'
import ManagePanel from './components/ManagePanel'
import DepositWithdrawForm from './components/DepositWithdrawForm'
import TransferForm from './components/TransferForm'
import LoanForm from './components/LoanForm'
import TransactionFeed from './components/TransactionFeed'
import AuthForm from './components/AuthForm'

const SESSION_KEY = 'zuri_session'
const SESSION_TTL_MS = 24 * 60 * 60 * 1000

function loadSession() {
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    if (!raw) return null
    const { accountId, expiresAt } = JSON.parse(raw)
    if (!accountId || Date.now() > expiresAt) {
      localStorage.removeItem(SESSION_KEY)
      return null
    }
    return accountId
  } catch {
    localStorage.removeItem(SESSION_KEY)
    return null
  }
}

function saveSession(accountId) {
  localStorage.setItem(
    SESSION_KEY,
    JSON.stringify({ accountId, expiresAt: Date.now() + SESSION_TTL_MS })
  )
}

function clearSession() {
  localStorage.removeItem(SESSION_KEY)
}

export default function App() {
  const [accounts, setAccounts] = useState([])
  const [selectedId, setSelectedId] = useState(null)
  const [activePanel, setActivePanel] = useState(null)
  const [ledger, setLedger] = useState([])
  const [loan, setLoan] = useState(null)
  const [msg, setMsg] = useState({ text: '', type: '' })
  const [busy, setBusy] = useState(false)
  const [restoring, setRestoring] = useState(true)

  const selected = accounts.find((a) => a.id === selectedId) || null

  useEffect(() => {
    const accountId = loadSession()
    if (!accountId) {
      setRestoring(false)
      return
    }
    api.getAccount(accountId)
      .then((acc) => {
        setAccounts((p) => (p.some((a) => a.id === acc.id) ? p : [...p, acc]))
        setSelectedId(acc.id)
      })
      .catch(() => clearSession())
      .finally(() => setRestoring(false))
  }, [])

  useEffect(() => {
    setLoan(null)
    if (selected) loadLedger(selected.id)
    else setLedger([])
  }, [selectedId])

  const flash = (text, type = 'error') => setMsg({ text, type })
  const clearMsg = () => setMsg({ text: '', type: '' })

  const loadLedger = async (id) => {
    try {
      const data = await api.getLedger(id)
      setLedger(data || [])
    } catch {
      setLedger([])
    }
  }

  const refreshSelected = async () => {
    if (!selectedId) return
    const fresh = await api.getAccount(selectedId)
    setAccounts((prev) => prev.map((a) => (a.id === fresh.id ? fresh : a)))
    loadLedger(selectedId)
  }

  const togglePanel = (key) => {
    clearMsg()
    setActivePanel((prev) => (prev === key ? null : key))
  }

  const logOut = () => {
    setSelectedId(null)
    setActivePanel(null)
    clearMsg()
    clearSession()
  }

  if (restoring) {
    return (
      <div>
        <TopBar />
      </div>
    )
  }

  if (!selected) {
    return (
      <div>
        <TopBar />
        {msg.text && <div className={`msg ${msg.type}`}>{msg.text}</div>}
        <AuthForm
          busy={busy} setBusy={setBusy} flash={flash} clearMsg={clearMsg}
          onAuthenticated={(acc) => {
            setAccounts((p) => (p.some((a) => a.id === acc.id) ? p : [...p, acc]))
            setSelectedId(acc.id)
            saveSession(acc.id)
          }}
        />
      </div>
    )
  }

  return (
    <div>
      <TopBar />

      {msg.text && <div className={`msg ${msg.type}`}>{msg.text}</div>}

      <HeroBalance account={selected} onManage={() => togglePanel('manage')} onLogout={logOut} />

      <Panel open={activePanel === 'manage'}>
        <ManagePanel
          account={selected} busy={busy} setBusy={setBusy} flash={flash} clearMsg={clearMsg}
          onDeleted={() => {
            setAccounts((p) => p.filter((a) => a.id !== selected.id))
            logOut()
          }}
        />
      </Panel>

      <ActionGrid activePanel={activePanel} onToggle={togglePanel} />

      <Panel open={activePanel === 'deposit' || activePanel === 'withdraw'}>
        <DepositWithdrawForm
          mode={activePanel} account={selected} busy={busy} setBusy={setBusy}
          flash={flash} clearMsg={clearMsg} onDone={refreshSelected}
        />
      </Panel>

      <Panel open={activePanel === 'transfer'}>
        <TransferForm
          accounts={accounts} fromAccount={selected} busy={busy} setBusy={setBusy}
          flash={flash} clearMsg={clearMsg} onDone={refreshSelected}
        />
      </Panel>

      <Panel open={activePanel === 'loan'}>
        <LoanForm
          account={selected} loan={loan} setLoan={setLoan} busy={busy} setBusy={setBusy}
          flash={flash} clearMsg={clearMsg} onDone={refreshSelected}
        />
      </Panel>

      <div className="feed-title">Transaction history</div>
      <TransactionFeed entries={ledger} />
    </div>
  )
}