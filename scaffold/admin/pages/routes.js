export default function () {
  return {
    'index': {
      title: 'Home',
      icon: 'fa-home',
      isRouteOnly: true
    },
    'pages': {
      title: 'Pages',
      icon: 'fa-file-text',
      children: {
        'pages-ID-pageEdit': {
          title: function (instance) {
            return instance.$route.params.pageID === 0 ? 'Create Page' : 'Edit Page'
          }
        },
        'pages-ID-pageLayout': {
          title: 'Edit Page Layout'
        }
      }
    },
    'works': {
      title: 'Portfolio',
      icon: 'fa-heart-o',
      children: {
        'works-workID': {
          title: function (instance) {
            return instance.$route.params.workID === '0' ? 'Create Portfolio Job' : 'Edit Portfolio Job'
          }
        }
      }
    },
    'people': {
      title: 'Administrators',
      icon: 'fa-users',
      children: {
        'people-personID': {
          title: function (instance) {
            return instance.$route.params.personID === 0 ? 'Create Administrator' : 'Edit Administrator'
          }
        }
      }
    },
    // 'settings': {
    //   title: 'Site Settings',
    //   icon: 'fa-cog'
    // },
    'login': {
      isRouteOnly: true
    },
    'logout': {
      name: 'logout',
      title: 'Logout',
      icon: 'fa-sign-out',
      path: '/logout' // custom path
    }
  }
}
