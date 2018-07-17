<template>
<div class="content u-fh">
  <div class="columns u-fh is-vcentered -u-m">
    <div class="column is-6 is-offset-3">
      <div class="logo-wrap has-text-centered u-fw">
        <img class="logo" src="~public/images/Logo-hoz.png" alt="">
      </div>
      <div class="box">
        <h1 class="is-title is-bold">Login</h1>
        <div v-show="error" style="color:red; word-wrap:break-word;">{{ error }}</div>
        <form v-on:submit.prevent="submitLogin">
          <div class="field">
            <label class="label">Email</label>
            <p class="control">
              <input v-model="details.Email" class="input" type="text" placeholder="email@example.org">
            </p>
          </div>
          <div class="field">
            <label class="label">Password</label>
            <p class="control">
              <input v-model="details.Password" class="input" type="password" placeholder="password">
            </p>
          </div>
          <div class="field has-addons u-mt">
            <p class="control">
              <button type="submit" class="button is-primary">Login</button>
            </p>
            <p class="control">
              <button class="button is-default">Cancel</button>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import {mapActions} from 'vuex'

export default {
  // COMPONENT
  // ______________________________________
  name: 'Login',
  layout: 'simple',
  components: {},
  props: {},
  computed: {},
  methods: {
    ...mapActions({
      login: 'auth/login'
    }),
    submitLogin () {
      this.login(this.details)
    }
  },
  watch: {},
  data () {
    return {
      details: {
        Email: null,
        Password: null,
        rememberMe: false
      },
      error: ''
    }
  },

  // LIFECYCLE METHODS
  // ______________________________________
  beforeCreate () {
  },
  created () {
  },
  beforeMount () {
  },
  mounted () {
  },
  beforeUpdate () {
  },
  updated () {
  },
  beforeDestroy () {
  },
  destroyed () {
  }
}
</script>

