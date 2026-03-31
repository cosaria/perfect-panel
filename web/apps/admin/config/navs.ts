export const navs = [
  {
    title: 'Dashboard',
    url: '/dashboard',
    icon: 'flat-color-icons:globe',
  },

  {
    title: 'Maintenance',
    icon: 'flat-color-icons:data-protection',
    items: [
      {
        title: 'Server Management',
        url: '/dashboard/servers',
        icon: 'flat-color-icons:data-protection',
      },
      { title: 'Node Management', url: '/dashboard/nodes', icon: 'flat-color-icons:mind-map' },
      { title: 'Subscribe Config', url: '/dashboard/subscribe', icon: 'flat-color-icons:ruler' },
      { title: 'Product Management', url: '/dashboard/product', icon: 'flat-color-icons:shop' },
    ],
  },

  {
    title: 'Commerce',
    icon: 'flat-color-icons:sales-performance',
    items: [
      { title: 'Order Management', url: '/dashboard/order', icon: 'flat-color-icons:todo-list' },
      { title: 'Coupon Management', url: '/dashboard/coupon', icon: 'flat-color-icons:bookmark' },
      {
        title: 'Marketing Management',
        url: '/dashboard/marketing',
        icon: 'flat-color-icons:bullish',
      },
      {
        title: 'Announcement Management',
        url: '/dashboard/announcement',
        icon: 'flat-color-icons:advertising',
      },
    ],
  },

  {
    title: 'Users & Support',
    icon: 'flat-color-icons:collaboration',
    items: [
      {
        title: 'User Management',
        url: '/dashboard/user',
        icon: 'flat-color-icons:conference-call',
      },
      {
        title: 'Ticket Management',
        url: '/dashboard/ticket',
        icon: 'flat-color-icons:collaboration',
      },
      {
        title: 'Document Management',
        url: '/dashboard/document',
        icon: 'flat-color-icons:document',
      },
    ],
  },

  {
    title: 'System',
    icon: 'flat-color-icons:services',
    items: [
      { title: 'System Config', url: '/dashboard/system', icon: 'flat-color-icons:services' },
      {
        title: 'Auth Control',
        url: '/dashboard/auth-control',
        icon: 'flat-color-icons:lock-portrait',
      },
      {
        title: 'Payment Config',
        url: '/dashboard/payment',
        icon: 'flat-color-icons:currency-exchange',
      },
      { title: 'ADS Config', url: '/dashboard/ads', icon: 'flat-color-icons:electrical-sensor' },
    ],
  },

  {
    title: 'Logs & Analytics',
    icon: 'flat-color-icons:statistics',
    items: [
      { title: 'Login', url: '/dashboard/log/login', icon: 'flat-color-icons:unlock' },
      { title: 'Register', url: '/dashboard/log/register', icon: 'flat-color-icons:contacts' },
      { title: 'Email', url: '/dashboard/log/email', icon: 'flat-color-icons:feedback' },
      { title: 'Mobile', url: '/dashboard/log/mobile', icon: 'flat-color-icons:sms' },
      { title: 'Subscribe', url: '/dashboard/log/subscribe', icon: 'flat-color-icons:workflow' },
      {
        title: 'Reset Subscribe',
        url: '/dashboard/log/reset-subscribe',
        icon: 'flat-color-icons:refresh',
      },
      {
        title: 'Subscribe Traffic',
        url: '/dashboard/log/subscribe-traffic',
        icon: 'flat-color-icons:statistics',
      },
      {
        title: 'Server Traffic',
        url: '/dashboard/log/server-traffic',
        icon: 'flat-color-icons:statistics',
      },
      {
        title: 'Traffic Details',
        url: '/dashboard/log/traffic-details',
        icon: 'flat-color-icons:combo-chart',
      },
      {
        title: 'Balance',
        url: '/dashboard/log/balance',
        icon: 'flat-color-icons:sales-performance',
      },
      { title: 'Commission', url: '/dashboard/log/commission', icon: 'flat-color-icons:debt' },
      { title: 'Gift', url: '/dashboard/log/gift', icon: 'flat-color-icons:donate' },
    ],
  },
];

export function findNavByUrl(url: string) {
  function matchDynamicRoute(pattern: string, path: string): boolean {
    const regexPattern = pattern.replace(/:[^/]+/g, '[^/]+').replace(/\//g, '\\/');
    const regex = new RegExp(`^${regexPattern}$`);
    return regex.test(path);
  }
  function findNav(items: any[], url: string, path: any[] = []): any[] {
    for (const item of items) {
      if (item.url === url || (item.url && matchDynamicRoute(item.url, url))) {
        return [...path, item];
      }
      if (item.items) {
        const result = findNav(item.items, url, [...path, item]);
        if (result.length) return result;
      }
    }
    return [];
  }
  return findNav(navs, url);
}
