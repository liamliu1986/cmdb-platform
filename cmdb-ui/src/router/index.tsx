import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Login from '@/modules/auth/Login'
import AppLayout from '@/layouts/AppLayout'
import CITypeList from '@/modules/core/CITypeList'

export default function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<AppLayout />}>
          <Route index element={<CITypeList />} />
          <Route path="citypes" element={<CITypeList />} />
        </Route>
        <Route path="*" element={<div>404</div>} />
      </Routes>
    </BrowserRouter>
  )
}
