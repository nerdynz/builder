export default function ({ app, store }) {
  app.$axios.interceptors.request.use(function (config) {
    // console.log(config)
    if (store.state.auth.isValid) {
      let token = store.state.auth.details.token
      config.headers.Authorization = `Basic ${token}`
    }
    // Do something before request is sent
    // if (config.showProgress === false) { // if undefined then we still do it
    //   nprogress.start()
    // }
    return config
  }, function (error) {
    // console.log(error)
    // Do something with request error
    return Promise.reject(error)
  })

  app.$axios.interceptors.response.use(function (response) {
    // console.log(response)
    // Do something with response data
    // nprogress.done()

    return response
  }, function (error) {
    console.log(error)
    // nprogress.done()
    if (!error.response) {
      app.$toast.open({ message: error.message, type: 'is-danger', duration: 5000 })
      return
    }

    let response = error.response
    if (response && (response.status === 403 || (response.data && response.data.indexOf && response.data.indexOf('ciphertext too short') >= 0))) {
      // todo notify
      // app.$notify({title: `${response.status} ${response.statusText}`, message: response.body, type: 'error'})
      app.router.replace({ path: 'login' })
    }
    if (response.status !== 200) {
      let errorData = response.body
      if (response.data) {
        errorData = response.data
      }

      // Vue.rollbar.debug(errorData)
      try {
        errorData = JSON.parse(errorData)
      } catch (e) {
      }

      // console.log(errorData)
      let message = errorData
      if (errorData.Friendly) {
        message = errorData.Friendly
      }

      app.$snackbar.open({ message: message, type: 'is-danger', duration: 5000, position: 'is-top-right' })
    }
    // Do something with request error
    return Promise.reject(error)
  })
}
