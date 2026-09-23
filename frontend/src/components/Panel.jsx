export default function Panel({ open, children }) {
  return (
    <div className={`panel ${open ? 'open' : ''}`}>
      <div className="panel-inner">{children}</div>
    </div>
  )
}