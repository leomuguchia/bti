import axios from 'axios'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${import.meta.env.VITE_API_KEY}`,
  },
})

client.interceptors.response.use(
  (res) => res.data,
  (err) => {
    const message = err.response?.data?.error || err.message || 'request failed'
    return Promise.reject(new Error(message))
  }
)

export const api = {
  createAccount: (owner) => client.post('/accounts', { owner }),
  getAccount: (id) => client.get(`/accounts/${id}`),
  login: (nationalId, phoneNumber) => client.post('/login', { national_id: nationalId, phone_number: phoneNumber }),
  deleteAccount: (id) => client.delete(`/accounts/${id}`),
  getLedger: (id) => client.get(`/accounts/${id}/ledger`),
  deposit: (id, amount) => client.post(`/accounts/${id}/deposit`, { amount }),
  withdraw: (id, amount) => client.post(`/accounts/${id}/withdraw`, { amount }),
  transfer: (fromId, toId, amount) =>
    client.post('/transfer', { from_id: fromId, to_id: toId, amount }),
  createLoan: (accountId) => client.post('/loans', { account_id: accountId }),
  disburseLoan: (loanId, amount) =>
    client.post(`/loans/${loanId}/disburse`, { amount }),
}