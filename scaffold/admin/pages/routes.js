export default function () {
  return {
    'index': {
      title: 'Home',
      icon: 'fa-home',
      isRouteOnly: true
    },
    'pages': {
      title: 'Pages',
      icon: 'fa-file',
      children: {
        'pages-ID-pageEdit': {
          title: function (instance) {
            return instance.$route.params.pageID === 0 ? 'Create Page' : 'Edit Page'
          }
        },
        'pages-ID-pageLayout': {
          title: 'Edit Page Layout'
        },
        'pages-ID-pageSlides': {
          title: 'Edit Page Slides'
        }
      }
    },
    'events': {
      title: 'Events',
      icon: 'fa-heart',
      children: {
        'events-ID-pageEdit': {
          title: function (instance) {
            return instance.$route.params.pageID === 0 ? 'Create Page' : 'Edit Page'
          }
        },
        'events-ID-pageLayout': {
          title: 'Edit Page Layout'
        },
        'events-ID-pageSlides': {
          title: 'Edit Page Slides'
        }
      }
    },
    'categoryImages': {
      title: 'Category Images',
      icon: 'fa-image',
      children: {
        'categoryImages-ID-categoryImageEdit': {
          title: function (instance) {
            return instance.$route.params.ID === 0 ? 'Create Category Image' : 'Edit Category Image'
          }
        }
      }
    },

    'people': {
      title: 'Person',
      icon: 'fa-circle-o',
      children: {
        'people-ID-personEdit': {
          title: function (instance) {
            return instance.$route.params.ID === 0 ? 'Create Person' : 'Edit Person'
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
