<template>
  <div class="list-placeholder rel-overlay">
    <div v-if="showing">
      <slot></slot>
    </div>
    <div v-else-if="loading" class="u-rel">
     <b-loading :active="true" :canCancel="true"></b-loading>
    </div>
    <div class="no-data" v-else>
      <svg viewBox="0 0 80 80" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
          <defs></defs>
          <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
              <g id="folder-smile" fill-rule="nonzero" fill="#000000">
                  <path d="M69.999,6.667 L40,6.667 L33.333,0 L10,0 C8.167,0 6.667,1.497 6.667,3.333 L6.667,13.333 L73.333,13.333 L73.333,10 C73.332,8.164 71.832,6.667 69.999,6.667 Z" id="Shape"></path>
                  <path d="M74.338,20 L5.661,20 C1.995,20 -0.511,22.956 0.091,26.575 L8.994,80 L71.005,80 L79.908,26.575 C80.514,22.956 78.007,20 74.338,20 Z M49.999,37.5 C52.301,37.5 54.166,39.362 54.166,41.667 C54.166,43.965 52.301,45.833 49.999,45.833 C47.697,45.833 45.832,43.965 45.832,41.667 C45.832,39.362 47.697,37.5 49.999,37.5 Z M30,37.5 C32.302,37.5 34.167,39.362 34.167,41.667 C34.167,43.965 32.302,45.833 30,45.833 C27.698,45.833 25.833,43.965 25.833,41.667 C25.833,39.362 27.698,37.5 30,37.5 Z M56.666,58.9658752 L53.131,62.5008752 C49.511,58.8808752 44.755,57.0708752 39.999,57.0708752 C35.244,57.0708752 30.488,58.8808752 26.868,62.5008752 L23.333,58.9658752 C32.542,49.7598752 47.467,49.7668752 56.666,58.9658752 Z" id="Shape"></path>
              </g>
          </g>
      </svg>
      <h3 class="title is-3">Oops. There is nothing here!</h3>
      <button @click="add" class="button is-success is-large">Add New</button>
    </div>
  </div>
</template>

<script>
// ... imports

export default {
  // COMPONENT
  // ______________________________________
  name: 'ListPlaceholder',
  components: {},
  props: {
    showing: Boolean,
    loading: Boolean
  },
  computed: {},
  methods: {
    add () {
      this.$emit('create')
    }
  },
  watch: {},
  data () {
    return {
      hello: true
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

<style lang="scss">
  @import "~public/scss/_variables";

  .list-placeholder {
    height: 100%;
   .no-data {
     user-select: none;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      svg {
        width: 25%;
      }
      #folder-smile {
      fill: $grey-light;
      }
      .title {
        margin-top: 2rem;
        color: $grey-light;
        font-weight: bold;
        text-align: center;
        text-transform: uppercase;
      }
   } 
  }
</style>
