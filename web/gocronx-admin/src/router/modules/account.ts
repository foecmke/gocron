import { AppRouteRecord } from '@/types/router'

/**
 * 账户自助路由（个人中心、两步验证）。
 *
 * 这些是所有登录用户都能访问的自助页面（后端 urlAuth 对普通用户放行
 * /api/user/editMyPassword 和 /api/user/2fa/*），因此父级不限制 roles，
 * 任何角色都会注册这些路由。
 *
 * 父级 isHide，子项也 isHide，所以不出现在侧边栏；子项用绝对路径保留原有
 * URL（/system/user-center、/system/two-factor），书签不失效。
 */
export const accountRoutes: AppRouteRecord = {
  path: '/account',
  name: 'Account',
  component: '/index/index',
  meta: {
    title: 'menus.system.userCenter',
    isHide: true
  },
  children: [
    {
      path: '/system/user-center',
      name: 'UserCenter',
      component: '/system/user-center',
      meta: {
        title: 'menus.system.userCenter',
        isHide: true,
        keepAlive: true,
        isHideTab: true,
        allowAbsolutePath: true
      }
    },
    {
      path: '/system/two-factor',
      name: 'TwoFactor',
      component: '/system/user-center/two-factor',
      meta: {
        title: 'menus.system.twoFactor',
        isHide: true,
        keepAlive: false,
        allowAbsolutePath: true
      }
    }
  ]
}
