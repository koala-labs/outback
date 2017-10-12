<template>
  <div>
    <button :disabled="disabled" @click="onClick" class="button">start deploy <span class="arrow">➤</span></button>
  </div>
</template>

<script>

export default {
  props: {
    onClick: {
      type: Function,
      required: true,
    },
    disabled: {
      type: Boolean,
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
