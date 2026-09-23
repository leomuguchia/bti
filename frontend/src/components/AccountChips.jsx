export default function AccountChips({ accounts, selectedId, onSelect, onNew }) {
  return (
    <div className="chip-row">
      {accounts.map((a) => (
        <button
          key={a.id}
          className={`chip ${a.id === selectedId ? 'active' : ''}`}
          onClick={() => onSelect(a.id)}
        >
          {a.owner.full_name}
        </button>
      ))}
      <button className="chip new" onClick={onNew}>+ New account</button>
    </div>
  )
}