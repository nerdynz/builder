<template>
  <div class="picture-avatar">
    <div class="avatar-img">
      <img v-if="user.Picture" :src="user.Picture" :alt="user.Name">
      <avatar v-else :username="user.Name" />
    </div>
    <div class="name">
      <span v-if="fullname">{{ user.Name }}</span>
      <span v-else>{{ firstName }}</span>
    </div>
  </div>
</template>

<script>
import Avatar from 'vue-avatar'

export default {
  // COMPONENT
  // ______________________________________
  name: 'PictureAvatar',
  components: {
    Avatar
  },
  props: {
    user: {
      type: Object,
      default: null
    },
    fullname: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
    }
  },
  computed: {
    firstName: function () {
      if (this.user.Name) {
        return this.user.Name.split(' ')[0]
      }
      return ''
    }
  }
}
</script>

<style lang="scss">
  @import "~public/scss/_variables";
  .picture-avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    .avatar-img {
      padding-top: 0.3rem;
      img {
        width: 50px;
        height: 50px;
        max-height: 50px;
        border-radius: 50px;
      }
    }
    .name {
      padding: 0.5rem;
      color: $white;
      font-weight: bold;
    }
  }
</style>
