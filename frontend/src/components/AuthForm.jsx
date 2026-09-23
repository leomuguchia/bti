import { useState } from 'react'
import { api } from '../api/client'

export default function AuthForm({ onAuthenticated, busy, setBusy, flash, clearMsg }) {
  const [tab, setTab] = useState('login')

  return (
    <div>
      <div className="row" style={{ marginBottom: 16 }}>
        <button type="button" className={tab === 'login' ? '' : 'secondary'} onClick={() => { setTab('login'); clearMsg() }}>
          Log in
        </button>
        <button type="button" className={tab === 'create' ? '' : 'secondary'} onClick={() => { setTab('create'); clearMsg() }}>
          Open account
        </button>
      </div>

      {tab === 'login'
        ? <LoginTab onAuthenticated={onAuthenticated} busy={busy} setBusy={setBusy} flash={flash} clearMsg={clearMsg} />
        : <CreateTab onAuthenticated={onAuthenticated} busy={busy} setBusy={setBusy} flash={flash} clearMsg={clearMsg} />}
    </div>
  )
}

function LoginTab({ onAuthenticated, busy, setBusy, flash, clearMsg }) {
  const [nationalId, setNationalId] = useState('')
  const [phoneNumber, setPhoneNumber] = useState('')

  const submit = async (e) => {
    e.preventDefault()
    if (!nationalId || !phoneNumber) return flash('enter your national ID and phone number')
    setBusy(true)
    clearMsg()
    try {
      const acc = await api.login(nationalId, phoneNumber)
      onAuthenticated(acc)
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <form onSubmit={submit}>
      <h3>Log in to Zuri Bank demo</h3>
      <label>National ID number</label>
      <input value={nationalId} onChange={(e) => setNationalId(e.target.value)} placeholder="e.g. 30112233" />
      <label>Phone number</label>
      <input value={phoneNumber} onChange={(e) => setPhoneNumber(e.target.value)} placeholder="e.g. 0712345678" />
      <button disabled={busy}>{busy ? 'Checking...' : 'Log in'}</button>
      <p className="muted-note">
        Identified by ID and phone rather than a password for this demo.
      </p>
    </form>
  )
}

function CreateTab({ onAuthenticated, busy, setBusy, flash, clearMsg }) {
  const [form, setForm] = useState({
    fullName: '', nationalId: '', phoneNumber: '', physicalAddress: '', kraPin: '',
  })
  const [idPhoto, setIdPhoto] = useState(null)
  const [passportPhoto, setPassportPhoto] = useState(null)

  const update = (key) => (e) => setForm((f) => ({ ...f, [key]: e.target.value }))

  const submit = async (e) => {
    e.preventDefault()
    const { fullName, nationalId, phoneNumber, physicalAddress } = form
    if (!fullName || !nationalId || !phoneNumber || !physicalAddress) {
      return flash('fill in all required fields')
    }
    if (!idPhoto || !passportPhoto) {
      return flash('ID photo and passport photo are both required')
    }
    setBusy(true)
    clearMsg()
    try {
      const acc = await api.createAccount({
        full_name: fullName,
        national_id: nationalId,
        phone_number: phoneNumber,
        physical_address: physicalAddress,
        kra_pin: form.kraPin || undefined,
        id_photo_ref: idPhoto.name,
        passport_photo_ref: passportPhoto.name,
      })
      onAuthenticated(acc)
    } catch (err) {
      flash(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <form onSubmit={submit}>
      <h3>Open a Zuri Bank account</h3>
      <label>Full name</label>
      <input value={form.fullName} onChange={update('fullName')} placeholder="e.g. Alex Muguchia" />
      <label>National ID number</label>
      <input value={form.nationalId} onChange={update('nationalId')} placeholder="e.g. 30112233" />
      <label>Phone number</label>
      <input value={form.phoneNumber} onChange={update('phoneNumber')} placeholder="e.g. 0712345678" />
      <label>Physical address</label>
      <input value={form.physicalAddress} onChange={update('physicalAddress')} placeholder="e.g. Ongata Rongai" />
      <label>KRA PIN (optional)</label>
      <input value={form.kraPin} onChange={update('kraPin')} placeholder="e.g. A012345678B" />
      <label>ID photo</label>
      <input type="file" accept="image/*" onChange={(e) => setIdPhoto(e.target.files[0] || null)} />
      <label>Passport-size photo</label>
      <input type="file" accept="image/*" onChange={(e) => setPassportPhoto(e.target.files[0] || null)} />
      <button disabled={busy}>{busy ? 'Opening account...' : 'Open account'}</button>
    </form>
  )
}