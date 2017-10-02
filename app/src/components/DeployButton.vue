<template>
  <div>
    <button :disabled="!isCompleteUnit" @click="makeDeployment" class="button">start deploy <span class="arrow">➤</span></button>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  methods: {
    ...mapActions([
      'runDeploy',
    ]),
    makeDeployment() {
      this.runDeploy(this.deployUnit);
      this.$notify({
        group: 'ufo',
        title: 'UFO',
        text: 'Your deployment has been scheduled',
      });
    },
  },
  computed: {
    ...mapState({
      deployUnit: state => state.deployUnit,
    }),
    isCompleteUnit() {
      if (this.deployUnit.cluster && this.deployUnit.service && this.deployUnit.version) {
        return true;
      }
      return false;
    },
  },
};
</script>

<style lang="scss" scoped>
.arrow {
  font-size: 14px;
}

.button {
  box-sizing: border-box;
  cursor: pointer;
  display: inline-block;
  transition: all .3s ease-in-out;
  background-color: #4688F1;
  color: white;
  font-size: 20px;
  font-weight: bold;
  letter-spacing: 2px;
  padding: 13px 23px;
  text-align: center;
  text-decoration: none;
  border: none;
  
  &:disabled {
    background-color: white;
    color: #4688F1;
  }
}
</style>
