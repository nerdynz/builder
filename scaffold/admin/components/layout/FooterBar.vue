<template>
  <footer class="footer" :class="{'no-sidebar': sidebarIsOpen}" :style="{'padding-left': sidebarIsOpen ? '260px' : '10px'}">
    <div class="level">
      <div class="level-left">
        <div v-for="item in leftButtons" class="level-item">
          <button class="button" :class="'is-' + item.kind" :type="item.kind" @click="(ev) => { item.click(ev) }" :disabled="item.isDisabled || !item.click">{{item.text}}</button>
        </div>
      </div>
      <div v-for="item in centerButtons" class="level-item">
        <button class="button" :class="'is-' + item.kind" :type="item.kind" @click="(ev) => { item.click(ev) }" :disabled="item.isDisabled || !item.click">{{item.text}}</button>
      </div>
      <div class="level-right">
        <div v-for="item in rightButtons" class="level-item">
          <button class="button" :class="'is-' + item.kind" :type="item.kind" @click="(ev) => { item.click(ev) }" :disabled="item.isDisabled || !item.click">{{item.text}}</button>
        </div>
      </div>
    </div>
  </footer>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'
export default {
  computed: {
    ...mapGetters({
      leftButtons: 'app/leftButtons',
      centerButtons: 'app/centerButtons',
      rightButtons: 'app/rightButtons',
      navbar: 'app/navbar',
      sidebarIsOpen: 'app/sidebarIsOpen'
    })
  },
  methods: {
    ...mapActions([
      'toggleJobView'
    ])
  },
  data () {
    return {
    }
  },
  mounted () {
  }
}
</script>

<style lang="scss">
@import "~public/scss/_variables";
@import '~bulma/utilities/mixins';

$footer-bar-height: 58px;

.footer {
  position: fixed;
  width: 100%;
  height: $footer-bar-height;
  padding: 10px;
  background: $grey-lightest;
  z-index: 500;
  bottom: 0;
  
  .level-left, .level-right {
    display: flex;
  }

  &.no-sidebar{
    padding-left: 10px;
  }

  @include mobile() {
    margin-left: 10px;
  }
  
  .social a {
    border-bottom: none !important;
  }

  .fa.fa-heart {
    color: red;
  }
}
.main-body {
  &:after {
    content: "";
    display: block;
    height: 78px;
    width: 100%;
  }
}

</style>
