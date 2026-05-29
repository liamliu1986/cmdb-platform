import { Layout, Menu } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { DatabaseOutlined, AppstoreOutlined, ApartmentOutlined, BankOutlined, RadarChartOutlined, RobotOutlined } from '@ant-design/icons'

const { Header, Sider, Content } = Layout

export default function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const menuItems = [
    { key: '/citypes', icon: <AppstoreOutlined />, label: '模型管理' },
    { key: '/ci', icon: <DatabaseOutlined />, label: '资源管理' },
    { key: '/ipam', icon: <ApartmentOutlined />, label: 'IPAM' },
    { key: '/dcim', icon: <BankOutlined />, label: 'DCIM' },
    { key: '/discovery', icon: <RadarChartOutlined />, label: '自动发现' },
    { key: '/agents', icon: <RobotOutlined />, label: 'Agent管理' },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#fff', padding: '0 24px', borderBottom: '1px solid #f0f0f0' }}>
        <h2 style={{ margin: 0 }}>CMDB</h2>
      </Header>
      <Layout>
        <Sider theme="light" style={{ borderRight: '1px solid #f0f0f0' }}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Content style={{ padding: 24, background: '#f5f5f5' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
