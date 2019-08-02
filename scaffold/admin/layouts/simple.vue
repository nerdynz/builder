<template>
  <main class="simple main-body">
    <nuxt />
  </main>
</template>

<script>

export default {

}
</script>

<style lang="scss" >
  @import "~public/scss/admin.scss";
  @import "~public/scss/admin-overrides.scss";
  @import "~public/scss/utility.scss";
  @import "~public/scss/_variables";
  @import "~public/scss/utility";

  .simple.main-body {
    background: $sidebar-color;
    .logo {
      padding: 2rem 4rem;
    }
  }
</style>
