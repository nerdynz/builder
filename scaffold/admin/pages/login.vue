<template>
  <div class="columns u-fh is-vcentered -u-m">
    <div class="column is-6 is-offset-3">
      <div class="logo-wrap has-text-centered u-fw">
        <logo />
      </div>
      <div class="box">
        <h1 class="title">
          Login
        </h1>
        <div v-show="error" style="color:red; word-wrap:break-word;">
          {{ error }}
        </div>
        <form @submit.prevent="submitLogin">
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
              <button type="submit" class="button is-primary">
                Login
              </button>
            </p>
            <p class="control">
              <button class="button is-default">
                Cancel
              </button>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import {
  mapActions
} from 'vuex'
import Logo from '~/components/Logo'

export default {
  // COMPONENT
  // ______________________________________
  name: 'Login',
  layout: 'simple',
  components: {
    Logo
  },
  props: {},
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
  computed: {},
  watch: {},
  methods: {
    ...mapActions({
      login: 'auth/login'
    }),
    submitLogin () {
      this.login(this.details)
    }
  }
}
</script>
