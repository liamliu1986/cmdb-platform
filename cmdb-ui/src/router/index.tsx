import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Login from '@/modules/auth/Login'
import AppLayout from '@/layouts/AppLayout'
import CITypeList from '@/modules/core/CITypeList'
import CIList from '@/modules/core/CIList'
import SubnetList from '@/modules/ipam/SubnetList'
import IDCList from '@/modules/dcim/IDCList'
import RuleList from '@/modules/discovery/RuleList'
import AgentList from '@/modules/discovery/AgentList'
import CIRelationGraph from '@/modules/core/CIRelationGraph'
import TopologyView from '@/modules/core/TopologyView'
import IDCMap from '@/modules/dcim/IDCMap'
import Dashboard from '@/modules/dashboard/Dashboard'
import RackView from '@/modules/dcim/RackView'

export default function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<AppLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="citypes" element={<CITypeList />} />
          <Route path="ci" element={<CIList />} />
          <Route path="ipam" element={<SubnetList />} />
          <Route path="dcim" element={<IDCList />} />
          <Route path="discovery" element={<RuleList />} />
          <Route path="agents" element={<AgentList />} />
          <Route path="relations" element={<CIRelationGraph />} />
          <Route path="topology" element={<TopologyView />} />
          <Route path="idcmap" element={<IDCMap />} />
          <Route path="rack" element={<RackView />} />
        </Route>
        <Route path="*" element={<div>404</div>} />
      </Routes>
    </BrowserRouter>
  )
}
