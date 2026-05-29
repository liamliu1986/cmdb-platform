import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Login from '@/modules/auth/Login'

export default function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<div>404</div>} />
      </Routes>
    </BrowserRouter>
  )
}
