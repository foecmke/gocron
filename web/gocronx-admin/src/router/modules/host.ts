import { AppRouteRecord } from '@/types/router'

export const hostRoutes: AppRouteRecord = {
  path: '/host',
  name: 'Host',
  component: '/index/index',
  meta: {
    title: 'menus.host.title',
    icon: 'ri:server-line',
    // Normal users may view hosts (backend urlAuth allows /api/host, /api/host/all).
    roles: ['R_SUPER', 'R_ADMIN', 'R_USER']
  },
  children: [
    {
      path: 'list',
      name: 'HostList',
      component: '/host/index',
      meta: {
        title: 'menus.host.list',
        icon: 'ri:hard-drive-2-line',
        keepAlive: true
      }
    },
    {
      path: 'create',
      name: 'HostCreate',
      component: '/host/edit',
      meta: {
        title: 'menus.host.create',
        isHide: true,
        keepAlive: false
      }
    },
    {
      path: 'edit/:id',
      name: 'HostEdit',
      component: '/host/edit',
      meta: {
        title: 'menus.host.edit',
        isHide: true,
        keepAlive: false
      }
    }
  ]
}
