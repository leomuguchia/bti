export default function HeroBalance({ account, onManage, onLogout }) {
  return (
    <div className="hero">
      <div className="owner">{account.owner.full_name}</div>
      <div className="amount">
        {account.balance.toLocaleString()}
        <span className="code">KES</span>
      </div>
      <div className="hero-foot">
        <span className="pill">{account.status}</span>
        <div style={{ display: 'flex', gap: 12 }}>
          <button className="link-muted" onClick={onManage}>Manage account</button>
          <button className="link-muted" onClick={onLogout}>Switch account</button>
        </div>
      </div>
    </div>
  )
}