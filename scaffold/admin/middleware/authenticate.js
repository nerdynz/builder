export default function ({ app, store, redirect, route }) {
  if (route.fullPath === '/login') {
    // cool
  } else if (!store.state.auth.isValid) {
    redirect('/login')
  }
}
